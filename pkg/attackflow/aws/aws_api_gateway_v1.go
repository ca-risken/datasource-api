package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	v1types "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

func (a *apiGatewayAnalyzer) analyzeV1(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// https://docs.aws.amazon.com/apigateway/latest/api/API_GetRestApis.html
	apis, err := a.v1client.GetRestApis(ctx, &apigateway.GetRestApisInput{})
	if err != nil && handleAPIGatewayError(err) != nil {
		return nil, err
	}
	for _, api := range apis.Items {
		if aws.ToString(api.Id) == a.resource.ShortName {
			a.metadata.ApiName = aws.ToString(api.Name)
			a.metadata.Description = aws.ToString(api.Description)
			a.metadata.DisableExecuteApiEndpoint = api.DisableExecuteApiEndpoint
			a.metadata.HasAPIResourcePolicy = aws.ToString(api.Policy) != ""
			break
		}
	}

	// https://docs.aws.amazon.com/apigateway/latest/api/API_GetResources.html
	resources, err := a.v1client.GetResources(ctx, &apigateway.GetResourcesInput{
		RestApiId: aws.String(a.resource.ShortName),
		Embed:     []string{"methods"},
	})
	if err != nil && handleAPIGatewayError(err) != nil {
		return nil, err
	}

	for _, r := range resources.Items {
		for _, m := range r.ResourceMethods {
			if m.MethodIntegration == nil {
				continue
			}
			// https://docs.aws.amazon.com/apigateway/latest/api/API_Integration.html
			target := aws.ToString(m.MethodIntegration.Uri)
			if m.MethodIntegration.Type == v1types.IntegrationTypeAws ||
				m.MethodIntegration.Type == v1types.IntegrationTypeAwsProxy {
				target = extractArnFromMethodIntegration(target)
			}
			if target == "" {
				continue
			}

			dest := apiGatewayIntegration{
				APIKeyRequired:    aws.ToBool(m.ApiKeyRequired),
				AuthorizationType: aws.ToString(m.AuthorizationType),
				RouteKey:          aws.ToString(m.HttpMethod) + " " + aws.ToString(r.Path),
				Target:            target,
			}
			a.metadata.Destination = append(a.metadata.Destination, dest)
			if !a.metadata.HasAPIResourcePolicy && !dest.APIKeyRequired && dest.AuthorizationType == "NONE" {
				a.metadata.IsPublic = true
			}
		}
	}

	// https://docs.aws.amazon.com/apigateway/latest/api/API_GetStages.html
	stages, err := a.v1client.GetStages(ctx, &apigateway.GetStagesInput{
		RestApiId: aws.String(a.resource.ShortName),
	})
	if err != nil && handleAPIGatewayError(err) != nil {
		return nil, err
	}
	for _, s := range stages.Item {
		if s.AccessLogSettings != nil && aws.ToString(s.AccessLogSettings.DestinationArn) != "" {
			a.metadata.Logging = true
		}
		if s.WebAclArn != nil && aws.ToString(s.WebAclArn) != "" {
			a.metadata.WafEnabled = true
		}
	}

	a.resource.MetaData, err = attackflow.ParseMetadata(a.metadata)
	if err != nil {
		return nil, err
	}
	resp = attackflow.SetNode(a.metadata.IsPublic, "api", a.resource, resp)
	return resp, nil
}
