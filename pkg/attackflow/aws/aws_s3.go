package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type s3Analyzer struct {
	resource  *datasource.Resource
	metadata  *S3Metadata
	awsConfig *aws.Config
	client    *s3.Client
	logger    logging.Logger
}

func newS3Analyzer(arn string, cfg *aws.Config, logger logging.Logger) (attackflow.CloudServiceAnalyzer, error) {
	return &s3Analyzer{
		resource:  getAWSInfoFromARN(arn),
		metadata:  &S3Metadata{},
		awsConfig: cfg,
		client:    s3.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

type S3Metadata struct {
	Encryption string `json:"encryption"`
	IsPublic   bool   `json:"is_public"`
	Versioning bool   `json:"versioning"`

	// S3 Notification
	LambdaConfiguration      []string `json:"lambda_configuration"`
	SQSConfiguration         []string `json:"sqs_configuration"`
	SNSConfiguration         []string `json:"sns_configuration"`
	EventBridgeConfiguration string   `json:"event_bridge_configuration"`
}

func (s *s3Analyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// cache
	cachedResource, cachedMeta, err := getS3AttackFlowCache(s.resource.CloudId, s.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil && cachedMeta != nil {
		s.resource = cachedResource
		s.metadata = cachedMeta
		resp = attackflow.SetNode(cachedMeta.IsPublic, "", cachedResource, resp)
		return resp, nil
	}

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
		s.awsConfig, err = retrieveAWSCredentialWithRegion(ctx, *s.awsConfig, regionCode)
		if err != nil {
			return nil, err
		}
		s.client = s3.NewFromConfig(*s.awsConfig)
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

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketNotificationConfiguration.html
	notification, err := s.client.GetBucketNotificationConfiguration(ctx, &s3.GetBucketNotificationConfigurationInput{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, err
	}

	// metadata
	for _, rule := range encryption.ServerSideEncryptionConfiguration.Rules {
		s.metadata.Encryption = fmt.Sprint(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		break
	}
	s.metadata.IsPublic = policyStatus.PolicyStatus != nil && policyStatus.PolicyStatus.IsPublic
	s.metadata.Versioning = versioning.Status == types.BucketVersioningStatusEnabled
	for _, config := range notification.LambdaFunctionConfigurations {
		s.metadata.LambdaConfiguration = append(s.metadata.LambdaConfiguration, aws.ToString(config.LambdaFunctionArn))
	}
	for _, config := range notification.QueueConfigurations {
		s.metadata.SQSConfiguration = append(s.metadata.SQSConfiguration, aws.ToString(config.QueueArn))
	}
	for _, config := range notification.TopicConfigurations {
		s.metadata.SNSConfiguration = append(s.metadata.SNSConfiguration, aws.ToString(config.TopicArn))
	}
	if notification.EventBridgeConfiguration != nil {
		s.metadata.EventBridgeConfiguration = fmt.Sprint(notification.EventBridgeConfiguration)
	}

	s.resource.MetaData, err = attackflow.ParseMetadata(s.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(s.metadata.IsPublic, "", s.resource, resp)

	// set cache
	if err := attackflow.SetAttackFlowCache(s.resource.CloudId, s.resource.ResourceName, s.resource); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *s3Analyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []attackflow.CloudServiceAnalyzer, error,
) {
	analyzer := []attackflow.CloudServiceAnalyzer{}
	// TODO: support for EventBridge, SNS, SQS
	for _, arn := range s.metadata.SNSConfiguration {
		r := getAWSInfoFromARN(arn)
		resp.Edges = append(resp.Edges, attackflow.GetEdge(s.resource.ResourceName, r.ResourceName, "event"))
		resp.Nodes = append(resp.Nodes, r)
	}
	for _, arn := range s.metadata.SQSConfiguration {
		r := getAWSInfoFromARN(arn)
		resp.Edges = append(resp.Edges, attackflow.GetEdge(s.resource.ResourceName, r.ResourceName, "event"))
		resp.Nodes = append(resp.Nodes, r)
	}
	for _, arn := range s.metadata.LambdaConfiguration {
		r := getAWSInfoFromARN(arn)
		lambdaAnalyzer, err := newLambdaAnalyzer(ctx, r.ResourceName, s.awsConfig, s.logger)
		if err != nil {
			return nil, nil, err
		}
		resp.Edges = append(resp.Edges, attackflow.GetEdge(s.resource.ResourceName, r.ResourceName, "event"))
		analyzer = append(analyzer, lambdaAnalyzer)
	}

	return resp, analyzer, nil
}

func getS3ARNFromDomain(domain string) string {
	if !domainPatternS3.MatchString(domain) {
		return ""
	}
	// bucket-name.com.s3.{region}.amazonaws.com -> bucket-name.com
	bucketName := domainPatternS3.ReplaceAll([]byte(domain), []byte{})
	return fmt.Sprintf("arn:aws:s3:::%s", bucketName)
}

func getS3AttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *S3Metadata, error) {
	resource, err := attackflow.GetAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta S3Metadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}
