package attackflow

import (
	"context"

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
)

type AttackFlowAnalyzer interface {
	Analyze(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (*datasource.AnalyzeAttackFlowResponse, error)
}
