package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/cenkalti/backoff/v4"
)

type AzureServiceClient interface {
	VerifyCode(ctx context.Context, subscriptionID, verificationCode string) (bool, error)
}

type AzureClient struct {
	logger  logging.Logger
	retryer backoff.BackOff
	cred    *azidentity.DefaultAzureCredential
}

func NewAzureClient(ctx context.Context, l logging.Logger) (AzureServiceClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return &AzureClient{
		logger:  l,
		retryer: backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10),
		cred:    cred,
	}, nil
}
