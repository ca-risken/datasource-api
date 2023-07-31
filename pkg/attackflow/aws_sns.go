package attackflow

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type snsAnalyzer struct {
	resource  *datasource.Resource
	metadata  *snsMetadata
	awsConfig *aws.Config
	client    *sns.Client
	logger    logging.Logger
}
type snsMetadata struct {
	Name          string            `json:"name"`
	Policy        string            `json:"policy"`
	Owner         string            `json:"owner"`
	Subscriptions []SnsSubscription `json:"subscriptions"`
}

type SnsSubscription struct {
	Endpoint        string `json:"endpoint"`
	Owner           string `json:"owner"`
	Protocol        string `json:"protocol"`
	SubscriptionArn string `json:"subscription_arn"`
}

func newSnsAnalyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	resource := getAWSInfoFromARN(arn)
	var err error
	if cfg.Region != resource.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, resource.Region)
		if err != nil {
			return nil, err
		}
	}
	return &snsAnalyzer{
		resource:  resource,
		metadata:  &snsMetadata{},
		awsConfig: cfg,
		client:    sns.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (s *snsAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// cache
	cachedResource, cachedMeta, err := getSnsAttackFlowCache(s.resource.CloudId, s.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil && cachedMeta != nil {
		s.logger.Infof(ctx, "cache hit: %+v", cachedResource)
		s.resource = cachedResource
		s.metadata = cachedMeta
		resp = setNode(false, "sns", cachedResource, resp)
		return resp, nil
	}

	// https://docs.aws.amazon.com/sns/latest/api/API_GetTopicAttributes.html
	topic, err := s.client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: aws.String(s.resource.ResourceName),
	})
	if err != nil {
		return nil, err
	}
	sliceArn := strings.Split(s.resource.ResourceName, ":")
	s.metadata.Name = sliceArn[len(sliceArn)-1]
	s.metadata.Policy = topic.Attributes["Policy"]
	s.metadata.Owner = topic.Attributes["Owner"]

	// https://docs.aws.amazon.com/sns/latest/api/API_ListSubscriptionsByTopic.html
	subscriptions, err := s.client.ListSubscriptionsByTopic(ctx, &sns.ListSubscriptionsByTopicInput{
		TopicArn: &s.resource.ResourceName,
	})
	if err != nil {
		return nil, err
	}
	for _, subscription := range subscriptions.Subscriptions {
		sub := SnsSubscription{
			Endpoint:        aws.ToString(subscription.Endpoint),
			Owner:           aws.ToString(subscription.Owner),
			Protocol:        aws.ToString(subscription.Protocol),
			SubscriptionArn: aws.ToString(subscription.SubscriptionArn),
		}
		s.metadata.Subscriptions = append(s.metadata.Subscriptions, sub)
	}
	s.resource.MetaData, err = parseMetadata(s.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(false, "sns", s.resource, resp)
	// cache
	if err := setAttackFlowCache(s.resource.CloudId, s.resource.ResourceName, s.resource); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *snsAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	for _, subscription := range s.metadata.Subscriptions {
		switch subscription.Protocol {
		case "lambda":
			resp.Edges = append(resp.Edges, getEdge(s.resource.ResourceName, subscription.Endpoint, "subscription"))
			lambdaAnalyzer, err := newLambdaAnalyzer(ctx, subscription.Endpoint, s.awsConfig, s.logger)
			analyzers = append(analyzers, lambdaAnalyzer)
			if err != nil {
				return nil, nil, err
			}
		case "sqs":
			resp.Edges = append(resp.Edges, getEdge(s.resource.ResourceName, subscription.Endpoint, "subscription"))
			sqsAnalyzer, err := newSqsAnalyzer(ctx, subscription.Endpoint, s.awsConfig, s.logger)
			analyzers = append(analyzers, sqsAnalyzer)
			if err != nil {
				return nil, nil, err
			}
		default:
			r := getExternalServiceNode(subscription.Endpoint)
			resp.Edges = append(resp.Edges, getEdge(s.resource.ResourceName, subscription.Endpoint, "subscription"))
			resp.Nodes = append(resp.Nodes, r)
		}
	}
	return resp, analyzers, nil
}
