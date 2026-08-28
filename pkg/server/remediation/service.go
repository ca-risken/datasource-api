package remediation

import (
	"context"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	coreai "github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type aiDBClient interface {
	db.AWSRepoInterface
}

type sqsAPI interface {
	Send(ctx context.Context, url string, msg interface{}) (*awssqs.SendMessageOutput, error)
}

type RemediationService struct {
	dbClient                    aiDBClient
	findingClient               finding.FindingServiceClient
	coreAIClient                coreai.AIServiceClient
	sqs                         sqsAPI
	remediationProposalQueueURL string
	logger                      logging.Logger
}

func NewRemediationService(dbClient aiDBClient, findingClient finding.FindingServiceClient, coreAIClient coreai.AIServiceClient, q *queue.Client, l logging.Logger) *RemediationService {
	return &RemediationService{
		dbClient:                    dbClient,
		findingClient:               findingClient,
		coreAIClient:                coreAIClient,
		sqs:                         q,
		remediationProposalQueueURL: q.AWSRemediationProposalQueueURL,
		logger:                      l,
	}
}
