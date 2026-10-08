package remediation

import (
	"context"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

type sqsAPI interface {
	Send(ctx context.Context, url string, msg interface{}) (*awssqs.SendMessageOutput, error)
}

type RemediationService struct {
	dbClient                    db.AWSRepoInterface
	findingClient               finding.FindingServiceClient
	aiClient                    ai.AIServiceClient
	sqs                         sqsAPI
	remediationProposalQueueURL string
	logger                      logging.Logger
}

func NewRemediationService(dbClient db.AWSRepoInterface, findingClient finding.FindingServiceClient, aiClient ai.AIServiceClient, q *queue.Client, l logging.Logger) *RemediationService {
	return &RemediationService{
		dbClient:                    dbClient,
		findingClient:               findingClient,
		aiClient:                    aiClient,
		sqs:                         q,
		remediationProposalQueueURL: q.AWSRemediationProposalQueueURL,
		logger:                      l,
	}
}
