package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type sqsAnalyzer struct {
	resource  *datasource.Resource
	metadata  *sqsMetadata
	awsConfig *aws.Config
	client    *sqs.Client
	logger    logging.Logger
}
type sqsMetadata struct {
	Name                 string          `json:"name"`
	Policy               string          `json:"policy"`
	VisibilityTimeout    string          `json:"visibility_timeout"`
	KmsMasterKeyId       string          `json:"kms_master_key_id"`
	SqsManagedSseEnabled bool            `json:"sqs_managed_sse_enabled"`
	LambdaTrigger        []lambdaTrigger `json:"lambda_trigger"`
}

func newSqsAnalyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (attackflow.CloudServiceAnalyzer, error) {
	resource := getAWSInfoFromARN(arn)
	var err error
	if cfg.Region != resource.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, resource.Region)
		if err != nil {
			return nil, err
		}
	}
	return &sqsAnalyzer{
		resource:  resource,
		metadata:  &sqsMetadata{},
		awsConfig: cfg,
		client:    sqs.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (s *sqsAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// cache
	cachedResource, cachedMeta, err := getSqsAttackFlowCache(s.resource.CloudId, s.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil && cachedMeta != nil {
		s.logger.Infof(ctx, "cache hit: %+v", cachedResource)
		s.resource = cachedResource
		s.metadata = cachedMeta
		resp = attackflow.SetNode(false, "sqs", cachedResource, resp)
		return resp, nil
	}
	awsInfo := getAWSInfoFromARN(s.resource.ResourceName)

	queueUrl := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", awsInfo.Region, awsInfo.CloudId, s.resource.ShortName)
	// get queue attributes
	attributes, err := s.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: &queueUrl,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
	})
	if err != nil {
		return nil, err
	}
	s.metadata.Name = awsInfo.ShortName
	s.metadata.VisibilityTimeout = attributes.Attributes["VisibilityTimeout"]
	s.metadata.Policy = attributes.Attributes["Policy"]
	s.metadata.KmsMasterKeyId = attributes.Attributes["KmsMasterKeyId"]
	s.metadata.SqsManagedSseEnabled = attributes.Attributes["SqsManagedSseEnabled"] == "true"

	s.metadata.LambdaTrigger, err = getLambdaTrigger(ctx, s.resource.ResourceName, s.awsConfig)
	if err != nil {
		return nil, err
	}

	s.resource.MetaData, err = attackflow.ParseMetadata(s.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(false, "sqs", s.resource, resp)
	// cache
	if err := attackflow.SetAttackFlowCache(s.resource.CloudId, s.resource.ResourceName, s.resource); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *sqsAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []attackflow.CloudServiceAnalyzer, error,
) {
	analyzers := []attackflow.CloudServiceAnalyzer{}
	for _, trigger := range s.metadata.LambdaTrigger {
		resp.Edges = append(resp.Edges, attackflow.GetEdge(s.resource.ResourceName, trigger.FunctionArn, "trigger"))
		lambdaAnalyzer, err := newLambdaAnalyzer(ctx, trigger.FunctionArn, s.awsConfig, s.logger)
		analyzers = append(analyzers, lambdaAnalyzer)
		if err != nil {
			return nil, nil, err
		}
	}
	return resp, analyzers, nil
}

func getSqsAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *sqsMetadata, error) {
	resource, err := attackflow.GetAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta sqsMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}
