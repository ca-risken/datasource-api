package attackflow

import (
	"context"
	"fmt"

	"github.com/ca-risken/datasource-api/proto/datasource"
)

const (
	// cloud type
	CLOUD_TYPE_AWS = "aws"

	// service layer
	LAYER_INTERNET  = "INTERNET"
	LAYER_CDN       = "CDN"
	LAYER_DATASTORE = "DATASTORE"

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
	}
}

func getEdge(source, target, edgeLabel string) *datasource.ResourceRelationship {
	return &datasource.ResourceRelationship{
		RelationId:         fmt.Sprintf("ed-[%s]-[%s]", source, target),
		SourceResourceName: source,
		TargetResourceName: target,
		RelationLabel:      edgeLabel,
	}
}
