package gcp

import (
	"context"
	"fmt"
	"os"
	"time"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type GcpServiceClient interface {
	VerifyCode(ctx context.Context, gcpProjectID, verificationCode string) (bool, error)
	GetAsset(ctx context.Context, gcpProjectID, resourceName string) (*assetpb.ResourceSearchResult, error)
	DescribeInstance(ctx context.Context, projectID, zone, instanceName string) (*Compute, error)
}

type GcpClient struct {
	logger  logging.Logger
	retryer backoff.BackOff

	asset   *asset.Client
	crm     *cloudresourcemanager.Service
	compute *compute.Service
}

func NewGcpClient(ctx context.Context, credentialPath string, l logging.Logger) (GcpServiceClient, error) {
	as, err := asset.NewClient(ctx, option.WithCredentialsFile(credentialPath))
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate for Google Asset API client: %w", err)
	}
	crm, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(credentialPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create new Cloud Resource Manager service: err=%w", err)
	}
	compute, err := compute.NewService(ctx, option.WithCredentialsFile(credentialPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create new Compute service: err=%w", err)
	}

	// Remove credential file for Security
	if err := os.Remove(credentialPath); err != nil {
		return nil, fmt.Errorf("failed to remove file: path=%s, err=%w", credentialPath, err)
	}
	return &GcpClient{
		asset:   as,
		crm:     crm,
		compute: compute,
		logger:  l,
		retryer: backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10),
	}, nil
}

func (g *GcpClient) newRetryLogger(ctx context.Context, funcName string) func(error, time.Duration) {
	return func(err error, t time.Duration) {
		g.logger.Warnf(ctx, "[RetryLogger] %s error: duration=%+v, err=%+v", funcName, t, err)
	}
}
