package datasource

import (
	"context"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/proto/datasource"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CLOUD_TYPE_AWS = "aws"
)

type dsDBClient interface {
	db.DataSourceRepoInterface
	db.AWSRepoInterface
}

type DataSourceService struct {
	dbClient dsDBClient
	logger   logging.Logger
}

func NewDataSourceService(dbClient dsDBClient, l logging.Logger) *DataSourceService {
	return &DataSourceService{
		dbClient: dbClient,
		logger:   l,
	}
}

func (d *DataSourceService) CleanDataSource(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	if err := d.dbClient.CleanWithNoProject(ctx); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (d *DataSourceService) AnalyzeAttackFlow(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (*datasource.AnalyzeAttackFlowResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var analyzer attackflow.AttackFlowAnalyzer
	var err error
	switch req.CloudType {
	case CLOUD_TYPE_AWS:
		analyzer, err = attackflow.NewAWSAttackFlowAnalyzer(ctx, req, d.dbClient, d.logger)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to create aws analyzer: %s", err.Error())
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid cloud type: %s", req.CloudType)
	}

	resp, err := analyzer.Analyze(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze attack flow: %s", err.Error())
	}
	return resp, nil
}
