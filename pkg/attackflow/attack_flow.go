package attackflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ca-risken/datasource-api/proto/datasource"
)

const (
	// region
	REGION_GLOBAL = "global"

	// cloud type
	CLOUD_TYPE_AWS = "aws"
	CLOUD_TYPE_GCP = "gcp"

	// service layer
	LAYER_INTERNET         = "INTERNET"
	LAYER_CDN              = "CDN"
	LAYER_LB               = "LB"
	LAYER_GATEWAY          = "GATEWAY"
	LAYER_DATASTORE        = "DATASTORE"
	LAYER_COMPUTE          = "COMPUTE"
	LAYER_LATERAL_MOVEMENT = "LATERAL_MOVEMENT"
	LAYER_EXTERNAL_SERVICE = "EXTERNAL_SERVICE"
	LAYER_INTERNAL_SERVICE = "INTERNAL_SERVICE"
	LAYER_CODE_REPOSITORY  = "CODE_REPOSITORY"

	// common resource
	RESOURCE_INTERNET = "Internet"

	// hard limit
	MAX_ANALYZE_NUM = 100
)

type CSP interface {
	GetInitialServiceAnalyzer(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (CloudServiceAnalyzer, error)
}

type CloudServiceAnalyzer interface {
	Analyze(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (*datasource.AnalyzeAttackFlowResponse, error)
	Next(ctx context.Context, resp *datasource.AnalyzeAttackFlowResponse) (*datasource.AnalyzeAttackFlowResponse, []CloudServiceAnalyzer, error)
}

func existsInternetNode(nodes []*datasource.Resource) bool {
	for _, node := range nodes {
		if node.ResourceName == RESOURCE_INTERNET {
			return true
		}
	}
	return false
}

func getInternetNode() *datasource.Resource {
	return &datasource.Resource{
		ResourceName: RESOURCE_INTERNET,
		ShortName:    RESOURCE_INTERNET,
		Layer:        LAYER_INTERNET,
		Region:       REGION_GLOBAL,
		Service:      "internet",
	}
}

func GetEdge(source, target, edgeLabel string) *datasource.ResourceRelationship {
	return &datasource.ResourceRelationship{
		RelationId:         fmt.Sprintf("ed-[%s]-[%s]", source, target),
		SourceResourceName: source,
		TargetResourceName: target,
		RelationLabel:      edgeLabel,
	}
}

func ParseMetadata(metadata interface{}) (string, error) {
	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(metaJSON), nil
}

func SetNode(isPublic bool, internetEdgeLabel string, resource *datasource.Resource, resp *datasource.AnalyzeAttackFlowResponse) *datasource.AnalyzeAttackFlowResponse {
	if isPublic {
		internet := getInternetNode()
		if !existsInternetNode(resp.Nodes) {
			resp.Nodes = append(resp.Nodes, internet)
		}
		resp.Edges = append(resp.Edges, GetEdge(internet.ResourceName, resource.ResourceName, internetEdgeLabel))
	}
	resp.Nodes = append(resp.Nodes, resource)
	return resp
}

func GetExternalServiceNode(target string) *datasource.Resource {
	return &datasource.Resource{
		ResourceName: target,
		ShortName:    target,
		Layer:        LAYER_EXTERNAL_SERVICE,
		Region:       REGION_GLOBAL,
		Service:      "external-service",
	}
}

func GetInternalServiceNode(target, region string) *datasource.Resource {
	return &datasource.Resource{
		ResourceName: target,
		ShortName:    target,
		Layer:        LAYER_INTERNAL_SERVICE,
		Region:       region,
		Service:      "internal-service",
	}
}

func GetCodeRepositoryNode(repository, service string) *datasource.Resource {
	return &datasource.Resource{
		ResourceName: repository,
		ShortName:    repository,
		Layer:        LAYER_CODE_REPOSITORY,
		Region:       REGION_GLOBAL,
		Service:      service,
	}
}
