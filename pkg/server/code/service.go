package code

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

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

func NewCodeService(dataKey string, repo db.CodeRepoInterface, q *queue.Client, pj project.ProjectServiceClient, l logging.Logger) (code.CodeServiceServer, error) {
	key := []byte(dataKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher, err=%w", err)
	}
	return &CodeService{
		repository:    repo,
		sqs:           q,
		cipherBlock:   block,
		projectClient: pj,
		logger:        l,
	}, nil
}
