package attackflow

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type lambdaAnalyzer struct {
	resource  *datasource.Resource
	metadata  *lambdaMetadata
	awsConfig *aws.Config
	client    *lambda.Client
	logger    logging.Logger
}
type lambdaMetadata struct {
	Architectures  []string `json:"architectures"`
	Description    string   `json:"description"`
	EnvironmentKey []string `json:"environment_key"`
	MemorySize     int32    `json:"memory_size"`
	Role           string   `json:"role"`
	Runtime        string   `json:"runtime"`
	State          string   `json:"state"`
	Vpc            string   `json:"vpc"`
	IsPublic       bool     `json:"is_public"`
	FunctionURL    string   `json:"function_url"`
	Destination    []string `json:"destination"`
}

type lambdaTrigger struct {
	FunctionArn string `json:"function_arn"`
	State       string `json:"state"`
}

func newLambdaAnalyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	resource := getAWSInfoFromARN(arn)
	var err error
	if cfg.Region != resource.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, resource.Region)
		if err != nil {
			return nil, err
		}
	}
	return &lambdaAnalyzer{
		resource:  resource,
		metadata:  &lambdaMetadata{},
		awsConfig: cfg,
		client:    lambda.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (l *lambdaAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// cache
	cachedResource, cachedMeta, err := getLambdaAttackFlowCache(l.resource.CloudId, l.resource.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil && cachedMeta != nil {
		l.resource = cachedResource
		l.metadata = cachedMeta
		resp = setNode(cachedMeta.IsPublic, "function URL", cachedResource, resp)
		return resp, nil
	}

	// https://docs.aws.amazon.com/lambda/latest/dg/API_GetFunction.html
	conf, err := l.client.GetFunctionConfiguration(ctx, &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(l.resource.ResourceName),
	})
	if err != nil {
		return nil, err
	}
	// https://docs.aws.amazon.com/lambda/latest/dg/API_ListFunctionUrlConfigs.html
	urls, err := l.client.ListFunctionUrlConfigs(ctx, &lambda.ListFunctionUrlConfigsInput{
		FunctionName: aws.String(l.resource.ResourceName),
	})
	if err != nil {
		return nil, err
	}
	// https://docs.aws.amazon.com/lambda/latest/dg/API_GetFunctionEventInvokeConfig.html
	destination, err := l.client.ListFunctionEventInvokeConfigs(ctx, &lambda.ListFunctionEventInvokeConfigsInput{
		FunctionName: aws.String(l.resource.ResourceName),
	})
	if err != nil {
		return nil, err
	}

	l.metadata.Description = aws.ToString(conf.Description)
	for _, arch := range conf.Architectures {
		l.metadata.Architectures = append(l.metadata.Architectures, string(arch))
	}
	if conf.Environment != nil {
		for key := range conf.Environment.Variables {
			l.metadata.EnvironmentKey = append(l.metadata.EnvironmentKey, key)
		}
	}
	l.metadata.MemorySize = aws.ToInt32(conf.MemorySize)
	l.metadata.Role = aws.ToString(conf.Role)
	l.metadata.Runtime = string(conf.Runtime)
	l.metadata.State = string(conf.State)
	if conf.VpcConfig != nil {
		l.metadata.Vpc = aws.ToString(conf.VpcConfig.VpcId)
	}
	for _, url := range urls.FunctionUrlConfigs {
		l.metadata.FunctionURL = aws.ToString(url.FunctionUrl)
		l.metadata.IsPublic = url.AuthType == types.FunctionUrlAuthTypeNone
		break // check only first url
	}
	for _, dest := range destination.FunctionEventInvokeConfigs {
		destConf := dest.DestinationConfig
		if destConf.OnSuccess != nil && aws.ToString(destConf.OnSuccess.Destination) != "" {
			l.metadata.Destination = append(l.metadata.Destination,
				aws.ToString(destConf.OnSuccess.Destination))
		}
		if destConf.OnFailure != nil && aws.ToString(destConf.OnFailure.Destination) != "" {
			l.metadata.Destination = append(l.metadata.Destination,
				aws.ToString(destConf.OnFailure.Destination))
		}
	}
	l.resource.MetaData, err = parseMetadata(l.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(l.metadata.IsPublic, "function URL", l.resource, resp)

	// cache
	if err := setAttackFlowCache(l.resource.CloudId, l.resource.ResourceName, l.resource); err != nil {
		return nil, err
	}
	return resp, nil
}

func getLambdaTrigger(ctx context.Context, sourceArn string, lambdaClient *lambda.Client) ([]lambdaTrigger, error) {
	eventSourceMappings, err := lambdaClient.ListEventSourceMappings(ctx, &lambda.ListEventSourceMappingsInput{
		EventSourceArn: &sourceArn,
	})
	if err != nil {
		return nil, err
	}
	var triggers []lambdaTrigger
	for _, eventSourceMapping := range eventSourceMappings.EventSourceMappings {
		triggers = append(triggers, lambdaTrigger{
			FunctionArn: aws.ToString(eventSourceMapping.FunctionArn),
			State:       aws.ToString(eventSourceMapping.State),
		})
	}
	return triggers, nil
}

func (l *lambdaAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	// IAM role
	iamRole := getAWSInfoFromARN(l.metadata.Role)
	resp.Edges = append(resp.Edges, getEdge(l.resource.ResourceName, iamRole.ResourceName, "iam role"))
	iamAnalyzer, err := newIAMAnalyzer(iamRole.ResourceName, l.awsConfig, l.logger)
	if err != nil {
		return nil, nil, err
	}
	analyzers = append(analyzers, iamAnalyzer)

	// Destinations
	for _, dest := range l.metadata.Destination {
		r := getAWSInfoFromARN(dest)
		switch r.Service {
		case SERVICE_LAMBDA:
			resp.Edges = append(resp.Edges, getEdge(l.resource.ResourceName, r.ResourceName, "destination"))
			lambdaAnalyzer, err := newLambdaAnalyzer(ctx, r.ResourceName, l.awsConfig, l.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, lambdaAnalyzer)
		case SERVICE_SNS:
			resp.Edges = append(resp.Edges, getEdge(l.resource.ResourceName, r.ResourceName, "destination"))
			snsAnalyzer, err := newSnsAnalyzer(ctx, r.ResourceName, l.awsConfig, l.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, snsAnalyzer)
		case SERVICE_SQS:
			resp.Edges = append(resp.Edges, getEdge(l.resource.ResourceName, r.ResourceName, "destination"))
			sqsAnalyzer, err := newSqsAnalyzer(ctx, r.ResourceName, l.awsConfig, l.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, sqsAnalyzer)
		default:
			// TODO: support for EventBridge
			resp.Edges = append(resp.Edges, getEdge(l.resource.ResourceName, r.ResourceName, "destination"))
			resp.Nodes = append(resp.Nodes, r)
		}
	}
	return resp, analyzers, nil
}
