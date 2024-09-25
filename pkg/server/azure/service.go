package azure

import (
	"context"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/azure"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type AzureService struct {
	repository    db.AzureRepoInterface
	sqs           *queue.Client
	azureClient   azure.AzureServiceClient
	projectClient project.ProjectServiceClient
	logger        logging.Logger
}

func NewAzureService(
	ctx context.Context, a azure.AzureServiceClient, repo db.AzureRepoInterface, q *queue.Client, pj project.ProjectServiceClient, l logging.Logger,
) *AzureService {
	return &AzureService{
		repository:    repo,
		sqs:           q,
		azureClient:   a,
		projectClient: pj,
		logger:        l,
	}
}
