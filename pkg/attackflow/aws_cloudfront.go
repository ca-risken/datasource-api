package attackflow

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type cloudFrontAnalyzer struct {
	resource  *datasource.Resource
	metadata  *CloudFrontMetadata
	awsConfig *aws.Config
	client    *cloudfront.Client
	logger    logging.Logger
}

func newCloudFrontAnalyzer(arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	return &cloudFrontAnalyzer{
		resource:  getAWSInfoFromARN(arn),
		metadata:  &CloudFrontMetadata{},
		awsConfig: cfg,
		client:    cloudfront.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

type CloudFrontMetadata struct {
	DistributionID    string   `json:"distribution_id"`
	Description       string   `json:"description"`
	Status            string   `json:"status"` // Deployed or InProgress
	Enabled           bool     `json:"enabled"`
	DomainName        string   `json:"domain_name"`
	DefaultRootObject string   `json:"default_root_object"`
	Aliases           []string `json:"aliases"`
	Origins           []string `json:"origins"`
	GeoRestriction    []string `json:"geo_restriction"`
	Logging           string   `json:"logging"`
	WebACLId          string   `json:"web_acl_id"`
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

	var enabled bool
	var aliases, origins, geoRestriction []string
	var description, logging, defaultRootObject, waf string
	if d.Distribution.DistributionConfig != nil {
		description = *d.Distribution.DistributionConfig.Comment
		enabled = *d.Distribution.DistributionConfig.Enabled
		if d.Distribution.DistributionConfig.Aliases != nil {
			aliases = d.Distribution.DistributionConfig.Aliases.Items
		}
		if d.Distribution.DistributionConfig.Origins != nil {
			for _, origin := range d.Distribution.DistributionConfig.Origins.Items {
				origins = append(origins, *origin.DomainName)
			}
		}
		defaultRootObject = *d.Distribution.DistributionConfig.DefaultRootObject
		geoRestriction = d.Distribution.DistributionConfig.Restrictions.GeoRestriction.Items
		if d.Distribution.DistributionConfig.Logging != nil && *d.Distribution.DistributionConfig.Logging.Bucket != "" {
			logging = *d.Distribution.DistributionConfig.Logging.Bucket + "/" + *d.Distribution.DistributionConfig.Logging.Prefix
		}
		waf = *d.Distribution.DistributionConfig.WebACLId
	}
	meta := &CloudFrontMetadata{
		DistributionID:    *d.Distribution.Id,
		Description:       description,
		Status:            *d.Distribution.Status,
		Enabled:           enabled,
		DomainName:        *d.Distribution.DomainName,
		Aliases:           aliases,
		Origins:           origins,
		DefaultRootObject: defaultRootObject,
		GeoRestriction:    geoRestriction,
		Logging:           logging,
		WebACLId:          waf,
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	c.resource.MetaData = string(metaJSON)
	c.metadata = meta

	// add node
	if meta.Enabled {
		internet := getInternetNode()
		if !existsInternetNode(resp.Nodes) {
			resp.Nodes = append(resp.Nodes, internet)
		}
		resp.Edges = append(resp.Edges, getEdge(internet.ResourceName, c.resource.ResourceName, meta.DomainName))
	}
	resp.Nodes = append(resp.Nodes, c.resource)
	return resp, nil
}

func (c *cloudFrontAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	if c.metadata == nil || len(c.metadata.Origins) == 0 {
		return resp, analyzers, nil
	}
	for _, origin := range c.metadata.Origins {
		awsService := findAWSServiceFromDomain(origin)
		if !isSupportedAWSService(awsService) {
			c.logger.Warnf(ctx, "Not supported origin: %s", origin)
			continue
		}

		switch awsService {
		case SERVICE_S3:
			s3ARN := getS3ARNFromDomain(origin)
			resp.Edges = append(resp.Edges, getEdge(c.resource.ResourceName, s3ARN, "origin"))
			s3Analyzer, err := newS3Analyzer(s3ARN, c.awsConfig, c.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, s3Analyzer)
		}
	}
	return resp, analyzers, nil
}
