package attackflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

func (a *apiGatewayAnalyzer) analyzeV2(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	// https://docs.aws.amazon.com/apigatewayv2/latest/api-reference/apis-apiid.html
	apis, err := a.v2client.GetApis(ctx, &apigatewayv2.GetApisInput{})
	if err != nil && handleAPIGatewayError(err) != nil {
		return nil, err
	}
	for _, api := range apis.Items {
		if aws.ToString(api.ApiId) == a.resource.ShortName {
			a.metadata.ApiName = aws.ToString(api.Name)
			a.metadata.Description = aws.ToString(api.Description)
			a.metadata.DisableExecuteApiEndpoint = api.DisableExecuteApiEndpoint
			break
		}
	}

	// https://docs.aws.amazon.com/apigatewayv2/latest/api-reference/apis-apiid-routes.html#GetRoutes
	rs, err := a.v2client.GetRoutes(ctx, &apigatewayv2.GetRoutesInput{
		ApiId: aws.String(a.resource.ShortName),
	})
	if err != nil && handleAPIGatewayError(err) != nil {
		return nil, err
	}
	integrationMap := map[string]apiGatewayIntegration{}
	for _, r := range rs.Items {
		targetKey := aws.ToString(r.Target)
		if strings.HasPrefix(targetKey, "integrations/") {
			targetKey = strings.Split(targetKey, "/")[1]
		}
		integrationMap[targetKey] = apiGatewayIntegration{
			APIKeyRequired:    r.ApiKeyRequired,
			AuthorizationType: fmt.Sprint(r.AuthorizationType),
			RouteKey:          aws.ToString(r.RouteKey),
		}
	}

	// https://docs.aws.amazon.com/apigatewayv2/latest/api-reference/apis-apiid-integrations.html#GetIntegrations
	integration, err := a.v2client.GetIntegrations(ctx, &apigatewayv2.GetIntegrationsInput{
		ApiId: aws.String(a.resource.ShortName),
	})
	if err != nil && handleAPIGatewayError(err) != nil {
		return nil, err
	}
	for _, i := range integration.Items {
		if _, ok := integrationMap[aws.ToString(i.IntegrationId)]; ok {
			target := integrationMap[aws.ToString(i.IntegrationId)]
			target.Target = aws.ToString(i.IntegrationUri)
			integrationMap[aws.ToString(i.IntegrationId)] = target
		}
	}
	for _, i := range integrationMap {
		a.metadata.Destination = append(a.metadata.Destination, i)
		if !i.APIKeyRequired && i.AuthorizationType == "NONE" {
			a.metadata.IsPublic = true
		}
	}
	// https://docs.aws.amazon.com/apigatewayv2/latest/api-reference/apis-apiid-stages.html
	stages, err := a.v2client.GetStages(ctx, &apigatewayv2.GetStagesInput{
		ApiId: aws.String(a.resource.ShortName),
	})
	if err != nil && handleAPIGatewayError(err) != nil {
		return nil, err
	}
	for _, s := range stages.Items {
		if s.AccessLogSettings != nil && aws.ToString(s.AccessLogSettings.DestinationArn) != "" {
			a.metadata.Logging = true
		}
	}

	a.resource.MetaData, err = parseMetadata(*a.metadata)
	if err != nil {
		return nil, err
	}
	resp = setNode(a.metadata.IsPublic, "api", a.resource, resp)
	return resp, nil
}

func isV2ApiID(ctx context.Context, apiID string, cfg *aws.Config) (bool, error) {
	client := apigatewayv2.NewFromConfig(*cfg)
	// https://docs.aws.amazon.com/apigatewayv2/latest/api-reference/apis-apiid.html
	apis, err := client.GetApis(ctx, &apigatewayv2.GetApisInput{})
	if err != nil {
		return false, err
	}
	for _, api := range apis.Items {
		if aws.ToString(api.ApiId) == apiID {
			return true, nil
		}
	}
	return false, nil
}
