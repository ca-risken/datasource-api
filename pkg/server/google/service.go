package google

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type GoogleService struct {
	repository      db.GoogleRepoInterface
	sqs             queue.GoogleQueueAPI
	resourceManager ResourceManagerServiceClient
	projectClient   project.ProjectServiceClient
	logger          logging.Logger
}

func NewGoogleService(credentialPath string, repo db.GoogleRepoInterface, q queue.GoogleQueueAPI, pj project.ProjectServiceClient, l logging.Logger) *GoogleService {
	r := newResourceManagerClient(credentialPath, l)
	return &GoogleService{
		repository:      repo,
		sqs:             q,
		resourceManager: r,
		projectClient:   pj,
		logger:          l,
	}
}
