package attackflow

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type eventBridgeAnalyzer struct {
	resource  *datasource.Resource
	metadata  *eventBridgeMetadata
	awsConfig *aws.Config
	client    *eventbridge.Client
	logger    logging.Logger
}
type eventBridgeMetadata struct {
	Name   string            `json:"name"`
	Policy string            `json:"policy"`
	Rules  []eventBridgeRule `json:"rules"`
}

type eventBridgeRule struct {
	Name    string              `json:"name"`
	State   string              `json:"state"`
	Arn     string              `json:"arn"`
	Targets []eventBridgeTarget `json:"targets"`
}

type eventBridgeTarget struct {
	Arn     string `json:"arn"`
	RoleArn string `json:"role_arn"`
}

func newEventBridgeAnalyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	resource := getAWSInfoFromARN(arn)
	var err error
	if cfg.Region != resource.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, resource.Region)
		if err != nil {
			return nil, err
		}
	}
	return &eventBridgeAnalyzer{
		resource:  resource,
		metadata:  &eventBridgeMetadata{},
		awsConfig: cfg,
		client:    eventbridge.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (e *eventBridgeAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// cache
	cachedResource, cachedMeta, err := getEventBridgeAttackFlowCache(e.resource.CloudId, e.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil && cachedMeta != nil {
		e.logger.Infof(ctx, "cache hit: %+v", cachedResource)
		e.resource = cachedResource
		e.metadata = cachedMeta
		resp = setNode(false, "eventBridge", cachedResource, resp)
		return resp, nil
	}
	awsInfo := getAWSInfoFromARN(e.resource.ResourceName)

	e.metadata.Name = awsInfo.ShortName
	buses, err := e.client.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{
		NamePrefix: &awsInfo.ShortName,
	})
	if err != nil {
		return nil, err
	}
	for _, bus := range buses.EventBuses {
		if *bus.Name == awsInfo.ShortName {
			e.metadata.Policy = aws.ToString(bus.Policy)
			break
		}
	}
	rules, err := e.client.ListRules(ctx, &eventbridge.ListRulesInput{
		EventBusName: &awsInfo.ShortName,
	})
	if err != nil {
		return nil, err
	}
	for _, r := range rules.Rules {
		if r.State == types.RuleStateDisabled {
			continue
		}
		rule := eventBridgeRule{
			Name:  aws.ToString(r.Name),
			State: string(r.State),
			Arn:   aws.ToString(r.Arn),
		}
		targets, err := e.client.ListTargetsByRule(ctx, &eventbridge.ListTargetsByRuleInput{
			Rule:         r.Name,
			EventBusName: &awsInfo.ShortName,
		})
		if err != nil {
			return nil, err
		}
		for _, target := range targets.Targets {
			rule.Targets = append(rule.Targets, eventBridgeTarget{
				Arn:     aws.ToString(target.Arn),
				RoleArn: aws.ToString(target.RoleArn),
			})
		}
		e.metadata.Rules = append(e.metadata.Rules, rule)
	}

	e.resource.MetaData, err = parseMetadata(e.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(false, "eventBridge", e.resource, resp)
	// cache
	if err := setAttackFlowCache(e.resource.CloudId, e.resource.ResourceName, e.resource); err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *eventBridgeAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	for _, rule := range e.metadata.Rules {
		for _, target := range rule.Targets {
			targetInfo := getAWSInfoFromARN(target.Arn)
			switch targetInfo.Service {
			case "lambda":
				resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.Arn, "target"))
				analyzer, err := newLambdaAnalyzer(ctx, target.Arn, e.awsConfig, e.logger)
				if err != nil {
					return nil, nil, err
				}
				analyzers = append(analyzers, analyzer)
			case "sqs":
				resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.Arn, "target"))
				analyzer, err := newSqsAnalyzer(ctx, target.Arn, e.awsConfig, e.logger)
				if err != nil {
					return nil, nil, err
				}
				analyzers = append(analyzers, analyzer)
			case "sns":
				resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.Arn, "target"))
				analyzer, err := newSnsAnalyzer(ctx, target.Arn, e.awsConfig, e.logger)
				if err != nil {
					return nil, nil, err
				}
				analyzers = append(analyzers, analyzer)
			case "events":
				resourceType := getEventBridgeResourceType(target.Arn)
				resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.Arn, "target"))
				// APIの送信はanalyzerを作成せず、event-busがバックエンドにある場合のみanalyzerを作成する
				if resourceType == "event-bus" {
					analyzer, err := newEventBridgeAnalyzer(ctx, target.Arn, e.awsConfig, e.logger)
					if err != nil {
						return nil, nil, err
					}
					analyzers = append(analyzers, analyzer)
				}
			default:
				resp.Edges = append(resp.Edges, getEdge(e.resource.ResourceName, target.Arn, "target"))
			}
			if target.RoleArn != "" {
				resp.Edges = append(resp.Edges, getEdge(target.Arn, target.RoleArn, "role"))
				analyzer, err := newIAMAnalyzer(target.RoleArn, e.awsConfig, e.logger)
				if err != nil {
					return nil, nil, err
				}
				analyzers = append(analyzers, analyzer)
			}

		}
	}

	return resp, analyzers, nil
}

func getEventBridgeResourceType(arn string) string {
	// arn:aws:events:ap-northeast-1:123456789012:api-destination/resource-name/uuid -> api-destination
	// arn:aws:events:ap-northeast-1:123456789012:event-bus/resource-name -> event-bus
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return ""
	}
	resourceParts := strings.Split(parts[5], "/")
	if len(resourceParts) < 2 {
		return ""
	}
	return resourceParts[0]
}
