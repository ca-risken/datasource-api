package gcp

import (
	"context"
	"fmt"
	"strings"

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
	req *datasource.AnalyzeAttackFlowRequest,
	repo db.GoogleRepoInterface,
	c gcp.GcpServiceClient,
	logger logging.Logger,
) attackflow.CSP {
	return &GCP{
		cloudID: req.CloudId,
		client:  c,
		logger:  logger,
	}
}

func (g *GCP) GetInitialServiceAnalyzer(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (
	attackflow.CloudServiceAnalyzer, error,
) {
	r, err := g.getResource(ctx, req)
	if err != nil {
		return nil, err
	}

	// asset types: https://cloud.google.com/asset-inventory/docs/resource-name-format
	switch r.Service {
	case "compute.googleapis.com/Instance":
		return newComputeAnalyzer(ctx, r, g.client, g.logger)
	default:
		return nil, fmt.Errorf("unsupported asset type: %s", r.Service)
	}
}

func (g *GCP) getResource(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (*datasource.Resource, error) {
	var r *datasource.Resource
	// read cache
	cachedResource, err := attackflow.GetAttackFlowCache(req.CloudId, req.ResourceName)
	if err != nil {
		return nil, err
	}
	if cachedResource != nil {
		return cachedResource, nil
	}

	asset, err := g.client.GetAsset(ctx, req.CloudId, req.ResourceName)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, fmt.Errorf("asset not found: %s", req.ResourceName)
	}
	r = &datasource.Resource{
		ResourceName: asset.Name,
		ShortName:    asset.DisplayName,
		CloudId:      req.CloudId,
		CloudType:    attackflow.CLOUD_TYPE_GCP,
		Region:       asset.Location,
		Service:      asset.AssetType,
	}
	return r, nil
}

func getShortNameFromURL(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}
