package code

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/github"
	"github.com/ca-risken/datasource-api/pkg/queue"
	"github.com/ca-risken/datasource-api/proto/code"
)

type CodeService struct {
	repository             db.CodeRepoInterface
	sqs                    CodeQueue
	cipherBlock            cipher.Block
	projectClient          project.ProjectServiceClient
	logger                 logging.Logger
	codeGitleaksQueueURL   string
	codeDependencyQueueURL string
	codeCodeScanQueueURL   string
	githubClient           github.GithubServiceClient
	limitRepositorySizeKb  int
}

type CodeQueue interface {
	Send(ctx context.Context, url string, msg interface{}) (*sqs.SendMessageOutput, error)
}

func NewCodeService(dataKey string, repo db.CodeRepoInterface, q *queue.Client, pj project.ProjectServiceClient, limitRepositorySizeKb int, l logging.Logger) (code.CodeServiceServer, error) {
	key := []byte(dataKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher, err=%w", err)
	}

	githubClient := github.NewGithubClient("", l)

	return &CodeService{
		repository:             repo,
		sqs:                    q,
		cipherBlock:            block,
		projectClient:          pj,
		logger:                 l,
		codeGitleaksQueueURL:   q.CodeGitleaksQueueURL,
		codeDependencyQueueURL: q.CodeDependencyQueueURL,
		codeCodeScanQueueURL:   q.CodeCodeScanQueueURL,
		githubClient:           githubClient,
		limitRepositorySizeKb: limitRepositorySizeKb,
	}, nil
}
