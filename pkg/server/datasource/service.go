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

func (d *DataSourceService) AnalyzeAttackFlow(ctx context.Context, req *datasource.AnalyzeAttackFlowRequest) (
	*datasource.AnalyzeAttackFlowResponse, error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var csp attackflow.CSP
	var err error
	switch req.CloudType {
	case attackflow.CLOUD_TYPE_AWS:
		csp, err = attackflow.NewAWS(ctx, req, d.dbClient, d.logger)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to create aws: %s", err.Error())
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid cloud type: %s", req.CloudType)
	}
	serviceAnalyzer, err := csp.GetInitialServiceAnalyzer(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get initial analyzer: %s", err.Error())
	}

	resp := &datasource.AnalyzeAttackFlowResponse{}
	nextAnalyzerList := []attackflow.CloudServiceAnalyzer{}
	analyzeCounter := 0
	for {
		resp, err = serviceAnalyzer.Analyze(ctx, resp)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to analyze attack flow: %s", err.Error())
		}
		var next []attackflow.CloudServiceAnalyzer
		resp, next, err = serviceAnalyzer.Next(ctx, resp)
		if err != nil {
			return nil, err
		}
		if len(next) > 0 {
			// push next analyzer
			nextAnalyzerList = append(nextAnalyzerList, next...)
		}
		if len(nextAnalyzerList) == 0 {
			break
		}

		analyzeCounter++
		if attackflow.MAX_ANALYZE_NUM < analyzeCounter {
			d.logger.Warnf(ctx, "analyze num exceeded: %d", analyzeCounter)
			break
		}

		// pop next analyzer
		serviceAnalyzer = nextAnalyzerList[0]
		nextAnalyzerList = nextAnalyzerList[1:]
	}

	return resp, nil
}
