package ai

import (
	"context"

	aipb "github.com/ca-risken/datasource-api/proto/ai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*AIService) GenerateRemediationProposal(_ context.Context, _ *aipb.GenerateRemediationProposalRequest) (*aipb.GenerateRemediationProposalResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GenerateRemediationProposal is not implemented")
}
