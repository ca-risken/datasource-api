package datasource

import (
	"context"

	"github.com/ca-risken/datasource-api/proto/datasource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *DataSourceService) InvokeAIRemediation(ctx context.Context, req *datasource.InvokeAIRemediationRequest) (*datasource.InvokeAIRemediationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "InvokeAIRemediation is not implemented")
}
