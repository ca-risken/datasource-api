package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/asset/apiv1/assetpb"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/api/iterator"
)

func (g *GcpClient) GetAsset(ctx context.Context, gcpProjectID, resourceName string) (*assetpb.ResourceSearchResult, error) {
	operation := func() (*assetpb.ResourceSearchResult, error) {
		return g.getAsset(ctx, gcpProjectID, resourceName)
	}
	return backoff.RetryNotifyWithData(operation, g.retryer, g.newRetryLogger(ctx, "GetAsset"))
}

func (g *GcpClient) getAsset(ctx context.Context, gcpProjectID, resourceName string) (*assetpb.ResourceSearchResult, error) {
	req := &assetpb.SearchAllResourcesRequest{
		Scope: fmt.Sprintf("projects/%s", gcpProjectID),
		Query: fmt.Sprintf("name:%s", resourceName),
	}
	it := g.asset.SearchAllResources(ctx, req)
	r, err := it.Next() // get only first result
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	return r, nil
}
