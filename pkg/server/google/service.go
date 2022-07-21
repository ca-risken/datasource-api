package google

import (
	"context"
	"fmt"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type GoogleService struct {
	repository      db.GoogleRepoInterface
	sqs             *queue.Client
	resourceManager ResourceManagerServiceClient
	projectClient   project.ProjectServiceClient
	logger          logging.Logger
}

func NewGoogleService(ctx context.Context, credentialPath string, repo db.GoogleRepoInterface, q *queue.Client, pj project.ProjectServiceClient, l logging.Logger) (*GoogleService, error) {
	r, err := newResourceManagerClient(ctx, credentialPath, l)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource manager client: err=%w", err)
	}
	return &GoogleService{
		repository:      repo,
		sqs:             q,
		resourceManager: r,
		projectClient:   pj,
		logger:          l,
	}, nil
}
