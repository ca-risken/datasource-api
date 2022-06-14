package code

import (
	"context"
	"crypto/aes"
	"crypto/cipher"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
	"github.com/ca-risken/datasource-api/proto/code"
)

type CodeService struct {
	repository    db.CodeRepoInterface
	sqs           *queue.Client
	cipherBlock   cipher.Block
	projectClient project.ProjectServiceClient
	logger        logging.Logger
}

func NewCodeService(coreSvcAddr, dataKey string, repo db.CodeRepoInterface, q *queue.Client, pj project.ProjectServiceClient, l logging.Logger) code.CodeServiceServer {
	ctx := context.Background()
	key := []byte(dataKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		l.Fatal(ctx, err.Error())
	}
	return &CodeService{
		repository:    repo,
		sqs:           q,
		cipherBlock:   block,
		projectClient: pj,
		logger:        l,
	}
}
