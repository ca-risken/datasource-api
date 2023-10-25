package google

import (
	"context"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/gcp"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type GoogleService struct {
	repository    db.GoogleRepoInterface
	sqs           *queue.Client
	gcpClient     gcp.GcpServiceClient
	projectClient project.ProjectServiceClient
	logger        logging.Logger
}

func NewGoogleService(
	ctx context.Context, g gcp.GcpServiceClient, repo db.GoogleRepoInterface, q *queue.Client, pj project.ProjectServiceClient, l logging.Logger,
) (
	*GoogleService, error,
) {
	return &GoogleService{
		repository:    repo,
		sqs:           q,
		gcpClient:     g,
		projectClient: pj,
		logger:        l,
	}, nil
}
