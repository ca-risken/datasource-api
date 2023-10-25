package gcp

import (
	"context"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/gcp"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

type GCP struct {
	cloudID        string
	resource       *datasource.Resource
	initialService string
	client         gcp.GcpServiceClient
	logger         logging.Logger
}

func NewGCP(
	ctx context.Context,
	req *datasource.AnalyzeAttackFlowRequest,
	repo db.GoogleRepoInterface,
	c gcp.GcpServiceClient,
	logger logging.Logger,
) (attackflow.CSP, error) {
	r, err := c.GetAsset(ctx, req.CloudId, req.ResourceName)
	if err != nil {
		return nil, err
	}
	csp := &GCP{
		cloudID: req.CloudId,
		resource: &datasource.Resource{
			ResourceName: r.Name,
			ShortName:    r.DisplayName,
			CloudId:      req.CloudId,
			CloudType:    attackflow.CLOUD_TYPE_GCP,
			Region:       r.Location,
			Service:      r.AssetType,
		},
		initialService: r.AssetType,
		client:         c,
		logger:         logger,
	}
	return csp, nil
}

// func (g *GCP) getGCPInfoFromResourceName(resourceName string) (*datasource.Resource, error) {
// 	r, err := g.client.GetAsset(context.Background(), g.cloudID, resourceName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &datasource.Resource{
// 		ResourceName: resourceName,
// 		ShortName:    r.DisplayName,
// 		CloudId:      g.cloudID,
// 		CloudType:    "gcp",
// 		Region:       r.Location,
// 		Service:      r.AssetType,
// 	}, nil
// }

func (g *GCP) GetInitialServiceAnalyzer(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (attackflow.CloudServiceAnalyzer, error) {
	g.logger.Infof(ctx, "[REQUEST] %+v", req)
	g.logger.Infof(ctx, "[RESOURCE] %+v", g.resource)
	return nil, nil
}
