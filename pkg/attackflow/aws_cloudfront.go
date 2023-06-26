package attackflow

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type CloudFrontMetadata struct {
	DistributionID    string   `json:"distribution_id"`
	Status            string   `json:"status"` // Deployed or InProgress
	Enabled           bool     `json:"enabled"`
	DomainName        string   `json:"domain_name"`
	DefaultRootObject string   `json:"default_root_object"`
	Ailiases          []string `json:"ailiases"`
	Origins           []string `json:"origins"`
	GeoRestriction    []string `json:"geo_restriction"`
	Logging           string   `json:"logging"`
	WebACLId          string   `json:"web_acl_id"`
}

func (a *AWSAttackFlowAnalyzer) analyzeCloudFront(ctx context.Context, arn string) (*datasource.AnalyzeAttackFlowResponse, error) {
	// analyze cloudfront resource
	cf, meta, err := a.analyzeCloudFronResource(ctx, arn)
	if err != nil {
		return nil, err
	}

	// next
	if err := a.analyzeCloudFrontNext(ctx, cf, meta); err != nil {
		return nil, err
	}

	return &datasource.AnalyzeAttackFlowResponse{
		Nodes: a.nodes,
		Edges: a.edges,
	}, nil
}

func (a *AWSAttackFlowAnalyzer) analyzeCloudFronResource(ctx context.Context, arn string) (*datasource.Resource, *CloudFrontMetadata, error) {
	r := getAWSInfoFromARN(arn)

	// https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetDistribution.html#API_GetDistribution_ResponseElements
	resp, err := a.cf.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: aws.String(r.ShortName),
	})
	if err != nil {
		return nil, nil, err
	}

	var enabled bool
	var aliases, origins, geoRestriction []string
	var logging, defaultRootObject, waf string
	if resp.Distribution.DistributionConfig != nil {
		enabled = *resp.Distribution.DistributionConfig.Enabled
		if resp.Distribution.DistributionConfig.Aliases != nil {
			aliases = resp.Distribution.DistributionConfig.Aliases.Items
		}
		if resp.Distribution.DistributionConfig.Origins != nil {
			for _, origin := range resp.Distribution.DistributionConfig.Origins.Items {
				origins = append(origins, *origin.DomainName)
			}
		}
		defaultRootObject = *resp.Distribution.DistributionConfig.DefaultRootObject
		geoRestriction = resp.Distribution.DistributionConfig.Restrictions.GeoRestriction.Items
		logging = *resp.Distribution.DistributionConfig.Logging.Bucket + "/" + *resp.Distribution.DistributionConfig.Logging.Prefix
		waf = *resp.Distribution.DistributionConfig.WebACLId
	}
	meta := &CloudFrontMetadata{
		DistributionID:    *resp.Distribution.Id,
		Status:            *resp.Distribution.Status,
		Enabled:           enabled,
		DomainName:        *resp.Distribution.DomainName,
		Ailiases:          aliases,
		Origins:           origins,
		DefaultRootObject: defaultRootObject,
		GeoRestriction:    geoRestriction,
		Logging:           logging,
		WebACLId:          waf,
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return nil, nil, err
	}
	r.Layer = LAYER_CDN
	r.MetaData = string(metaJSON)

	// add node
	if meta.Enabled {
		a.addInternetNode(getNodeName(SERVICE_CLOUDFRONT, r.ResourceName), meta.DomainName)
	}
	a.nodes = append(a.nodes, r)
	return r, meta, nil
}

func (a *AWSAttackFlowAnalyzer) analyzeCloudFrontNext(ctx context.Context, cf *datasource.Resource, meta *CloudFrontMetadata) error {
	if len(meta.Origins) == 0 {
		return nil
	}
	for _, origin := range meta.Origins {
		awsService := findAWSServiceFromDomain(origin)
		if !isSupportedAWSService(awsService) {
			a.logger.Warnf(ctx, "Not supported origin: %s", origin)
			continue
		}

		switch awsService {
		case SERVICE_S3:
			s3ARN := getS3ARNFromDomain(origin)
			sr := getAWSInfoFromARN(s3ARN)
			a.addEdge(getNodeName(SERVICE_CLOUDFRONT, cf.ResourceName), getNodeName(SERVICE_S3, sr.ShortName), origin)
			if _, err := a.analyzeS3(ctx, s3ARN); err != nil {
				return err
			}
		}
	}
	return nil
}
