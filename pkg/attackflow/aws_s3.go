package attackflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type s3Analyzer struct {
	resource  *datasource.Resource
	metadata  *S3Metadata
	awsConfig *aws.Config
	client    *s3.Client
	logger    logging.Logger
}

func newS3Analyzer(arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	return &s3Analyzer{
		resource:  getAWSInfoFromARN(arn),
		metadata:  &S3Metadata{},
		awsConfig: cfg,
		client:    s3.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

type S3Metadata struct {
	Encryption string
	IsPublic   bool
	Versioning bool
}

func (s *s3Analyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	bucketName := aws.String(s.resource.ShortName)

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLocation.html
	location, err := s.client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, err
	}
	regionCode := fmt.Sprint(location.LocationConstraint)
	s.logger.Debugf(ctx, "s3: location=%v", regionCode)
	if regionCode != "" {
		s.resource.Region = regionCode
		if err := s.updateS3ClientWithRegion(ctx, regionCode); err != nil {
			return nil, err
		}
	}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketEncryption.html
	encryption, err := s.client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, err
	}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicyStatus.html
	policyStatus, err := s.client.GetBucketPolicyStatus(ctx, &s3.GetBucketPolicyStatusInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, err
	}

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketVersioning.html
	versioning, err := s.client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, err
	}

	sseEncrypt := ""
	for _, rule := range encryption.ServerSideEncryptionConfiguration.Rules {
		sseEncrypt = fmt.Sprint(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		break
	}

	meta := &S3Metadata{
		Encryption: sseEncrypt,
		IsPublic:   policyStatus.PolicyStatus != nil && policyStatus.PolicyStatus.IsPublic,
		Versioning: versioning.Status == types.BucketVersioningStatusEnabled,
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	s.resource.Layer = LAYER_DATASTORE
	s.resource.MetaData = string(metaJSON)
	s.metadata = meta

	// add node
	if meta.IsPublic {
		internet := getInternetNode()
		if !existsInternetNode(resp.Nodes) {
			resp.Nodes = append(resp.Nodes, internet)
		}
		resp.Edges = append(resp.Edges, getEdge(internet.ResourceName, s.resource.ResourceName, ""))
	}
	resp.Nodes = append(resp.Nodes, s.resource)
	return resp, nil
}

func (s *s3Analyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	return resp, []CloudServiceAnalyzer{}, nil
}

func (s *s3Analyzer) updateS3ClientWithRegion(ctx context.Context, region string) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryMaxAttempts(RETRY_MAX_ATTEMPT),
	)
	if err != nil {
		return err
	}
	cfg.Credentials = s.awsConfig.Credentials
	s.client = s3.NewFromConfig(cfg) // update s3 client
	return nil
}

func getS3ARNFromDomain(domain string) string {
	if !domainPatternS3.MatchString(domain) {
		return ""
	}
	// bucket-name.com.s3.{region}.amazonaws.com -> bucket-name.com
	bucketName := domainPatternS3.ReplaceAll([]byte(domain), []byte{})
	return fmt.Sprintf("arn:aws:s3:::%s", bucketName)
}
