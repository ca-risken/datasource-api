package diagnosis

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type DiagnosisService struct {
	repository    db.DiagnosisRepoInterface
	sqs           queue.DiagnosisQueueAPI
	projectClient project.ProjectServiceClient
	logger        logging.Logger
}

func NewDiagnosisService(repo db.DiagnosisRepoInterface, q queue.DiagnosisQueueAPI, pj project.ProjectServiceClient, l logging.Logger) *DiagnosisService {
	return &DiagnosisService{
		repository:    repo,
		sqs:           q,
		projectClient: pj,
		logger:        l,
	}
}
