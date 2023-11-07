package gcp

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/gcp"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type GCP struct {
	cloudID string
	client  gcp.GcpServiceClient
	logger  logging.Logger
}

func NewGCP(
	ctx context.Context,
	req *datasource.AnalyzeAttackFlowRequest,
	repo db.GoogleRepoInterface,
	c gcp.GcpServiceClient,
	logger logging.Logger,
) (attackflow.CSP, error) {
	csp := &GCP{
		cloudID: req.CloudId,
		client:  c,
		logger:  logger,
	}
	return csp, nil
}

func (g *GCP) GetInitialServiceAnalyzer(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (
	attackflow.CloudServiceAnalyzer, error,
) {
	asset, err := g.client.GetAsset(ctx, req.CloudId, req.ResourceName)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, fmt.Errorf("asset not found: %s", req.ResourceName)
	}
	g.logger.Infof(ctx, "[ASSET] %+v", asset)
	r := &datasource.Resource{
		ResourceName: asset.Name,
		ShortName:    asset.DisplayName,
		CloudId:      req.CloudId,
		CloudType:    attackflow.CLOUD_TYPE_GCP,
		Region:       asset.Location,
		Service:      asset.AssetType,
	}
	// asset types: https://cloud.google.com/asset-inventory/docs/resource-name-format
	switch asset.AssetType {
	case "compute.googleapis.com/Instance":
		return newComputeAnalyzer(ctx, r, g.client, g.logger)
	default:
		return nil, fmt.Errorf("unsupported asset type: %s", asset.AssetType)
	}
}
