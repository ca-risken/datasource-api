package aws

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type AWSService struct {
	repository    db.AWSRepoInterface
	sqs           *queue.Client
	projectClient project.ProjectServiceClient
	logger        logging.Logger
}

func NewAWSService(repo db.AWSRepoInterface, q *queue.Client, pj project.ProjectServiceClient, l logging.Logger) *AWSService {
	return &AWSService{
		repository:    repo,
		sqs:           q,
		projectClient: pj,
		logger:        l,
	}
}
