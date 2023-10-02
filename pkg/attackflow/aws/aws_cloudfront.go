package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type cloudFrontAnalyzer struct {
	resource  *datasource.Resource
	metadata  *CloudFrontMetadata
	awsConfig *aws.Config
	client    *cloudfront.Client
	logger    logging.Logger
}

func newCloudFrontAnalyzer(arn string, cfg *aws.Config, logger logging.Logger) (attackflow.CloudServiceAnalyzer, error) {
	return &cloudFrontAnalyzer{
		resource:  getAWSInfoFromARN(arn),
		metadata:  &CloudFrontMetadata{},
		awsConfig: cfg,
		client:    cloudfront.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

type CloudFrontMetadata struct {
	DistributionID    string    `json:"distribution_id"`
	Description       string    `json:"description"`
	Status            string    `json:"status"` // Deployed or InProgress
	Enabled           bool      `json:"enabled"`
	DomainName        string    `json:"domain_name"`
	DefaultRootObject string    `json:"default_root_object"`
	Aliases           []string  `json:"aliases"`
	Origins           []*origin `json:"origins"`
	GeoRestriction    []string  `json:"geo_restriction"`
	Logging           string    `json:"logging"`
	WebACLId          string    `json:"web_acl_id"`
}

type origin struct {
	DomainName string `json:"domain_name"`
	HTTPOnly   bool   `json:"http_only"`
}

func (c *cloudFrontAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetDistribution.html#API_GetDistribution_ResponseElements
	d, err := c.client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: aws.String(c.resource.ShortName),
	})
	if err != nil {
		return nil, err
	}

	c.metadata.DistributionID = aws.ToString(d.Distribution.Id)
	c.metadata.Status = aws.ToString(d.Distribution.Status)
	c.metadata.DomainName = aws.ToString(d.Distribution.DomainName)
	if d.Distribution.DistributionConfig != nil {
		c.metadata.Description = aws.ToString(d.Distribution.DistributionConfig.Comment)
		c.metadata.Enabled = aws.ToBool(d.Distribution.DistributionConfig.Enabled)
		if d.Distribution.DistributionConfig.Aliases != nil {
			c.metadata.Aliases = d.Distribution.DistributionConfig.Aliases.Items
		}
		if d.Distribution.DistributionConfig.Origins != nil {
			for _, o := range d.Distribution.DistributionConfig.Origins.Items {
				httpOnly := false
				if o.CustomOriginConfig != nil {
					httpOnly = o.CustomOriginConfig.OriginProtocolPolicy == types.OriginProtocolPolicyHttpOnly
				}
				c.metadata.Origins = append(c.metadata.Origins, &origin{
					DomainName: aws.ToString(o.DomainName),
					HTTPOnly:   httpOnly,
				})
			}
		}
		c.metadata.DefaultRootObject = *d.Distribution.DistributionConfig.DefaultRootObject
		c.metadata.GeoRestriction = d.Distribution.DistributionConfig.Restrictions.GeoRestriction.Items
		if d.Distribution.DistributionConfig.Logging != nil && *d.Distribution.DistributionConfig.Logging.Bucket != "" {
			c.metadata.Logging = *d.Distribution.DistributionConfig.Logging.Bucket + "/" + *d.Distribution.DistributionConfig.Logging.Prefix
		}
		c.metadata.WebACLId = *d.Distribution.DistributionConfig.WebACLId
	}
	c.resource.MetaData, err = attackflow.ParseMetadata(c.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(c.metadata.Enabled, c.metadata.DomainName, c.resource, resp)
	return resp, nil
}

func (c *cloudFrontAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []attackflow.CloudServiceAnalyzer, error,
) {
	analyzers := []attackflow.CloudServiceAnalyzer{}
	if c.metadata == nil || len(c.metadata.Origins) == 0 {
		return resp, analyzers, nil
	}
	for _, o := range c.metadata.Origins {
		awsService := findAWSServiceFromDomain(o.DomainName)
		if !isSupportedAWSService(awsService) {
			c.setCustomDomain(o, resp)
			continue
		}

		switch awsService {
		case SERVICE_S3:
			s3ARN := getS3ARNFromDomain(o.DomainName)
			s3Analyzer, err := newS3Analyzer(s3ARN, c.awsConfig, c.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, s3Analyzer)
			resp.Edges = append(resp.Edges, attackflow.GetEdge(c.resource.ResourceName, s3ARN, "origin"))
		case SERVICE_EC2:
			ec2ARN, err := getEC2ARNFromPublicDomain(ctx, o.DomainName, c.resource.CloudId, c.awsConfig)
			if err != nil {
				c.logger.Warnf(ctx, "failed to get ec2 ARN from public domain: %s", err)
				continue
			}
			ec2Analyzer, err := newEC2Analyzer(ctx, ec2ARN, c.awsConfig, c.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, ec2Analyzer)
			resp.Edges = append(resp.Edges, attackflow.GetEdge(c.resource.ResourceName, ec2ARN, "origin"))
		case SERVICE_ELB:
			elbARN, err := getElbARNFromPublicDomain(ctx, o.DomainName, c.resource.CloudId, c.awsConfig)
			if err != nil {
				c.logger.Warnf(ctx, "failed to get elb ARN from public domain: %s", err)
				continue
			}
			elbAnalyzer, err := newELBAnalyzer(ctx, elbARN, c.awsConfig, c.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, elbAnalyzer)
			resp.Edges = append(resp.Edges, attackflow.GetEdge(c.resource.ResourceName, elbARN, "origin"))
		case SERVICE_API_GATEWAY:
			apiGatewayARN, err := getAPIGatewayARNFromPublicDomain(ctx, o.DomainName, c.resource.CloudId, c.awsConfig)
			if err != nil {
				c.logger.Warnf(ctx, "failed to get api gateway ARN from public domain: %s", err)
				continue
			}
			apiGatewayAnalyzer, err := newAPIGatewayAnalyzer(ctx, apiGatewayARN, c.awsConfig, c.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, apiGatewayAnalyzer)
			resp.Edges = append(resp.Edges, attackflow.GetEdge(c.resource.ResourceName, apiGatewayARN, "origin"))
		default:
			c.logger.Warnf(ctx, "unsupported aws service: %s", awsService)
			c.setCustomDomain(o, resp)
		}
	}
	return resp, analyzers, nil
}

func (c *cloudFrontAnalyzer) setCustomDomain(o *origin, resp *datasource.AnalyzeAttackFlowResponse) {
	resourceName := "https://" + o.DomainName
	if o.HTTPOnly {
		resourceName = "http://" + o.DomainName
	}
	resp.Nodes = append(resp.Nodes, attackflow.GetExternalServiceNode(resourceName))
	resp.Edges = append(resp.Edges, attackflow.GetEdge(c.resource.ResourceName, resourceName, "origin"))
}
