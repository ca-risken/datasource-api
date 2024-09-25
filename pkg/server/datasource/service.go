package datasource

import (
	"context"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/pkg/attackflow/aws"
	"github.com/ca-risken/datasource-api/pkg/attackflow/gcp"
	"github.com/ca-risken/datasource-api/pkg/db"
	gcpsvc "github.com/ca-risken/datasource-api/pkg/gcp"
	"github.com/ca-risken/datasource-api/proto/datasource"
	"github.com/cenkalti/backoff/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/slack-go/slack"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dsDBClient interface {
	db.DataSourceRepoInterface
	db.AWSRepoInterface
	db.GoogleRepoInterface
	db.CodeRepoInterface
	db.DiagnosisRepoInterface
	db.OSINTRepoInterface
	db.AzureRepoInterface
}

type DataSourceService struct {
	dbClient      dsDBClient
	alertClient   alert.AlertServiceClient
	baseURL       string
	defaultLocale string
	gcpClient     gcpsvc.GcpServiceClient
	slackClient   *slack.Client
	logger        logging.Logger
	retryer       backoff.BackOff
}

func NewDataSourceService(
	dbClient dsDBClient, alertClient alert.AlertServiceClient, gcpClient gcpsvc.GcpServiceClient, slackClient *slack.Client, url, defaultLocale string, l logging.Logger,
) *DataSourceService {
	local := defaultLocale
	if local == "" {
		local = DEFAULT_LOCALE
	}
	return &DataSourceService{
		dbClient:      dbClient,
		alertClient:   alertClient,
		baseURL:       url,
		defaultLocale: local,
		gcpClient:     gcpClient,
		slackClient:   slackClient,
		logger:        l,
		retryer:       backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10),
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
		csp, err = aws.NewAWS(ctx, req, d.dbClient, d.logger)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "failed to create aws: %s", err.Error())
		}
	case attackflow.CLOUD_TYPE_GCP:
		if d.gcpClient == nil {
			return nil, status.Errorf(codes.FailedPrecondition, "gcp service is not available")
		}
		csp = gcp.NewGCP(req, d.gcpClient, d.logger)
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

func (d *DataSourceService) NotifyScanError(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	// Get scan errors for all projects
	scanErrors, err := d.getScanErrors(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to getScanErrors: %s", err.Error())
	}

	// Notify error per project
	for projectID, errs := range scanErrors {
		// notify
		resp, err := d.alertClient.ListNotificationForInternal(ctx, &alert.ListNotificationForInternalRequest{
			ProjectId: projectID,
			Type:      "slack",
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to ListNotification: %s", err.Error())
		}
		if len(resp.Notification) == 0 {
			continue
		}
		for _, n := range resp.Notification {
			if err := d.notifyScanError(ctx, n, errs); err != nil {
				d.logger.Notifyf(ctx, logging.WarnLevel, "Failed to notify scan error: project_id=%d, notification_id=%d, err=%s",
					n.ProjectId, n.NotificationId, err.Error())
			}
		}

		// update err_notified_at
		if err := d.updateScanErrorNotifiedAt(ctx, projectID, errs); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to updateScanErrorNotifiedAt: %s", err.Error())
		}
	}
	return &empty.Empty{}, nil
}

func (d *DataSourceService) newRetryLogger(ctx context.Context, funcName string) func(error, time.Duration) {
	return func(err error, t time.Duration) {
		d.logger.Warnf(ctx, "[RetryLogger] %s error: duration=%+v, err=%+v", funcName, t, err)
	}
}
