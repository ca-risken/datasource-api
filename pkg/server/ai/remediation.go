package ai

import (
	"context"

	protoai "github.com/ca-risken/datasource-api/proto/ai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*AIService) InvokeAIRemediation(_ context.Context, _ *protoai.InvokeAIRemediationRequest) (*protoai.InvokeAIRemediationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "InvokeAIRemediation is not implemented")
}
