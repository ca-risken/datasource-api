package attackflow

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	v1types "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	v2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type apiGatewayAnalyzer struct {
	resource  *datasource.Resource
	metadata  *apiGatewayMetadata
	awsConfig *aws.Config
	v1client  *apigateway.Client   // REST API
	v2client  *apigatewayv2.Client // HTTP API and WebSocket API
	logger    logging.Logger
}
type apiGatewayMetadata struct {
	ApiName                   string                  `json:"api_name"`
	Description               string                  `json:"description"`
	IsPublic                  bool                    `json:"is_public"`
	Destination               []apiGatewayIntegration `json:"destination"`
	DisableExecuteApiEndpoint bool                    `json:"disable_execute_api_endpoint"`
	HasAPIResourcePolicy      bool                    `json:"has_api_resource_policy"`
	Logging                   bool                    `json:"logging"`
	WafEnabled                bool                    `json:"waf_enabled"`
}

type apiGatewayIntegration struct {
	APIKeyRequired    bool   `json:"api_key_required"`
	AuthorizationType string `json:"authorization_type"`
	RouteKey          string `json:"route_key"`
	Target            string `json:"target"`
}

func newAPIGatewayAnalyzer(ctx context.Context, arn string, cfg *aws.Config, logger logging.Logger) (CloudServiceAnalyzer, error) {
	r := getAWSInfoFromARN(arn)
	apiID, err := extractApiID(arn)
	if err != nil {
		return nil, err
	}
	r.ShortName = apiID // update short name to API ID
	if cfg.Region != r.Region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, r.Region)
		if err != nil {
			return nil, err
		}
	}
	return &apiGatewayAnalyzer{
		resource:  r,
		metadata:  &apiGatewayMetadata{},
		awsConfig: cfg,
		v1client:  apigateway.NewFromConfig(*cfg),
		v2client:  apigatewayv2.NewFromConfig(*cfg),
		logger:    logger,
	}, nil
}

func (a *apiGatewayAnalyzer) Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	if strings.Contains(a.resource.ResourceName, "/restapis/") {
		return a.analyzeV1(ctx, resp)
	}
	if strings.Contains(a.resource.ResourceName, "/apis/") {
		return a.analyzeV2(ctx, resp)
	}
	return nil, fmt.Errorf("invalid resource name: resource_name=%s", a.resource.ResourceName)
}

func (a *apiGatewayAnalyzer) Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error,
) {
	analyzers := []CloudServiceAnalyzer{}
	for _, destination := range a.metadata.Destination {

		// external service
		if strings.HasPrefix(destination.Target, "http") {
			r := getExternalServiceNode(destination.Target)
			resp.Edges = append(resp.Edges, getEdge(a.resource.ResourceName, destination.Target, "integration"))
			resp.Nodes = append(resp.Nodes, r)
			continue
		}

		r := getAWSInfoFromARN(destination.Target)
		switch r.Service {
		case SERVICE_S3:
			resp.Edges = append(resp.Edges, getEdge(a.resource.ResourceName, r.ResourceName, "integration"))
			s3Analyzer, err := newS3Analyzer(r.ResourceName, a.awsConfig, a.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, s3Analyzer)
		case SERVICE_LAMBDA:
			resp.Edges = append(resp.Edges, getEdge(a.resource.ResourceName, r.ResourceName, "integration"))
			lambdaAnalyzer, err := newLambdaAnalyzer(ctx, r.ResourceName, a.awsConfig, a.logger)
			if err != nil {
				return nil, nil, err
			}
			analyzers = append(analyzers, lambdaAnalyzer)
		}
	}

	return resp, analyzers, nil
}

func handleAPIGatewayError(err error) error {
	var v1NotFound *v1types.NotFoundException
	var v2NotFound *v2types.NotFoundException
	if errors.As(err, &v1NotFound) || errors.As(err, &v2NotFound) {
		return nil
	}
	return err
}

func extractApiID(arn string) (string, error) {
	// https://docs.aws.amazon.com/apigateway/latest/developerguide/arn-format-reference.html
	splitArn := strings.Split(arn, ":")
	if len(splitArn) < 5 {
		return "", fmt.Errorf("invalid arn: arn=%s", arn)
	}
	resource := splitArn[5]
	if !strings.HasPrefix(resource, "/restapis/") && !strings.HasPrefix(resource, "/apis/") {
		return "", fmt.Errorf("no api-id: resource=%s", resource)
	}
	return strings.Split(resource, "/")[2], nil
}

// extractArnFromMethodIntegration extracts target arn from integration URI
// URI format: arn:aws:apigateway:{region}:{subdomain.service|service}:path|action/{service_api}
// https://docs.aws.amazon.com/apigateway/latest/api/API_Integration.html
func extractArnFromMethodIntegration(integrationURI string) string {
	if strings.Contains(integrationURI, "arn:aws:lambda") {
		// lambda
		// e.g.) arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:my-function/invocations -> arn:aws:lambda:us-east-1:123456789012:function:my-function
		split := strings.Split(integrationURI, "arn:aws:")
		if len(split) < 3 {
			return ""
		}
		return strings.Split("arn:aws:"+split[2], "/invocations")[0]
	} else if strings.Contains(integrationURI, ":s3:action/") {
		// s3:action
		// e.g.) arn:aws:apigateway:us-west-2:s3:action/GetObject&Bucket={bucket}&Key={key} -> arn:aws:s3:::{bucket}
		split := strings.Split(integrationURI, "Bucket=")
		if len(split) < 2 {
			return ""
		}
		return "arn:aws:s3:::" + strings.Split(split[1], "&")[0]
	} else if strings.Contains(integrationURI, ":s3:path/") {
		// s3:path
		// e.g.) arn:aws:apigateway:us-west-2:s3:path/{bucket}/{key} -> arn:aws:s3:::{bucket}
		split := strings.Split(integrationURI, "s3:path/")
		if len(split) < 2 {
			return ""
		}
		return "arn:aws:s3:::" + strings.Split(split[1], "/")[0]
	}
	// unsupported uri
	return ""
}

func getAPIGatewayARNFromPublicDomain(ctx context.Context, domain, accountID string, cfg *aws.Config) (string, error) {
	if !domainPatternAPIGateway.MatchString(domain) {
		return "", fmt.Errorf("invalid domain: %s", domain)
	}
	// domain format: {api-id}.execute-api.{region}.amazonaws.com
	apiID := strings.Split(domain, ".")[0]
	region := strings.Split(domain, ".")[2]
	arn := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s", region, apiID) // v1 (REST API)

	var err error
	if cfg.Region != region {
		cfg, err = retrieveAWSCredentialWithRegion(ctx, *cfg, region)
		if err != nil {
			return "", err
		}
	}
	isV2Api, err := isV2ApiID(ctx, apiID, cfg)
	if err != nil {
		return "", err
	}
	if isV2Api {
		arn = fmt.Sprintf("arn:aws:apigateway:%s::/apis/%s", region, apiID) // v2 (HTTP API)
	}
	return arn, nil
}
