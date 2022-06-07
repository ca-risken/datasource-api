package osint

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type OsintService struct {
	repository    db.OSINTRepoInterface
	sqs           queue.OSINTQueueAPI
	projectClient project.ProjectServiceClient
	logger        logging.Logger
}

func NewOsintService(repo db.OSINTRepoInterface, q queue.OSINTQueueAPI, pj project.ProjectServiceClient, l logging.Logger) *OsintService {
	return &OsintService{
		repository:    repo,
		sqs:           q,
		projectClient: pj,
		logger:        l,
	}
}
