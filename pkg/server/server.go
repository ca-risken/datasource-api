package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/project"
	azureClient "github.com/ca-risken/datasource-api/pkg/azure"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/gcp"
	"github.com/ca-risken/datasource-api/pkg/queue"
	awsServer "github.com/ca-risken/datasource-api/pkg/server/aws"
	azureServer "github.com/ca-risken/datasource-api/pkg/server/azure"
	codeServer "github.com/ca-risken/datasource-api/pkg/server/code"
	dsServer "github.com/ca-risken/datasource-api/pkg/server/datasource"
	diagnosisServer "github.com/ca-risken/datasource-api/pkg/server/diagnosis"
	googleServer "github.com/ca-risken/datasource-api/pkg/server/google"
	osintServer "github.com/ca-risken/datasource-api/pkg/server/osint"
	"github.com/ca-risken/datasource-api/proto/aws"
	"github.com/ca-risken/datasource-api/proto/azure"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/ca-risken/datasource-api/proto/datasource"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
	"github.com/ca-risken/datasource-api/proto/google"
	"github.com/ca-risken/datasource-api/proto/osint"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/slack-go/slack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type Server struct {
	port                 string
	coreSvcAddr          string
	awsRegion            string
	googleCredentialPath string
	dataKey              string
	db                   *db.Client
	queue                *queue.Client
	baseURL              string
	defaultLocale        string
	slackAPIToken        string
	logger               logging.Logger
}

func NewServer(port, coreSvcAddr, awsRegion, googleCredentialPath, dataKey string, db *db.Client, q *queue.Client, url, defaultLocale, slackAPIToken string, logger logging.Logger) *Server {
	return &Server{
		port:                 port,
		coreSvcAddr:          coreSvcAddr,
		awsRegion:            awsRegion,
		googleCredentialPath: googleCredentialPath,
		dataKey:              dataKey,
		db:                   db,
		queue:                q,
		baseURL:              url,
		defaultLocale:        defaultLocale,
		slackAPIToken:        slackAPIToken,
		logger:               logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	localServerAddr := fmt.Sprintf(":%s", s.port)
	pjClient, err := newProjectClient(s.coreSvcAddr)
	if err != nil {
		return fmt.Errorf("failed to create project client: %w", err)
	}
	alertClient, err := newAlertClient(s.coreSvcAddr)
	if err != nil {
		return fmt.Errorf("failed to create alert client: %w", err)
	}
	gcpClient, err := gcp.NewGcpClient(ctx, s.googleCredentialPath, s.logger)
	if err != nil {
		return fmt.Errorf("failed to create gcp client: %w", err)
	}
	if gcpClient == nil {
		s.logger.Warnf(ctx, "Google credential file not exists at %s", s.googleCredentialPath)
		s.logger.Warn(ctx, "Google service will not be available")
	}
	azureClient, err := azureClient.NewAzureClient(ctx, s.logger)
	if err != nil {
		return fmt.Errorf("failed to create azure client: %w", err)
	}
	slackClient := slack.New(s.slackAPIToken)

	awsSvc := awsServer.NewAWSService(s.db, s.queue, pjClient, s.logger)
	googleSvc := googleServer.NewGoogleService(ctx, gcpClient, s.db, s.queue, pjClient, s.logger)
	codeSvc, err := codeServer.NewCodeService(s.dataKey, s.db, s.queue, pjClient, s.logger)
	if err != nil {
		return fmt.Errorf("failed to create code service: %w", err)
	}
	osintSvc := osintServer.NewOsintService(s.db, s.queue, pjClient, s.logger)
	diagnosisSvc := diagnosisServer.NewDiagnosisService(s.db, s.queue, pjClient, s.logger)
	azureSvc := azureServer.NewAzureService(ctx, azureClient, s.db, s.queue, pjClient, s.logger)
	dsSvc := dsServer.NewDataSourceService(s.db, alertClient, gcpClient, slackClient, s.baseURL, s.defaultLocale, s.logger)
	hsvc := health.NewServer()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpctrace.UnaryServerInterceptor(),
				mimosarpc.LoggingUnaryServerInterceptor(s.logger))))
	aws.RegisterAWSServiceServer(server, awsSvc)
	google.RegisterGoogleServiceServer(server, googleSvc)
	code.RegisterCodeServiceServer(server, codeSvc)
	osint.RegisterOsintServiceServer(server, osintSvc)
	diagnosis.RegisterDiagnosisServiceServer(server, diagnosisSvc)
	azure.RegisterAzureServiceServer(server, azureSvc)
	datasource.RegisterDataSourceServiceServer(server, dsSvc)
	grpc_health_v1.RegisterHealthServer(server, hsvc)

	reflection.Register(server) // enable reflection API

	s.logger.Infof(ctx, "Starting gRPC server at %s", localServerAddr)
	l, err := net.Listen("tcp", localServerAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	errChan := make(chan error)
	go func() {
		if err := server.Serve(l); err != nil && err != grpc.ErrServerStopped {
			s.logger.Errorf(ctx, "failed to serve grpc: %w", err)
			errChan <- err
		}
	}()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		s.logger.Info(ctx, "Shutdown gRPC server...")
		server.GracefulStop()
	}
	return nil
}

func newProjectClient(svcAddr string) (project.ProjectServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: err=%w", err)
	}
	return project.NewProjectServiceClient(conn), nil
}

func newAlertClient(svcAddr string) (alert.AlertServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: err=%w", err)
	}
	return alert.NewAlertServiceClient(conn), nil
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(grpctrace.UnaryClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
