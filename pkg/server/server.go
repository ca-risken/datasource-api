package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
	awsServer "github.com/ca-risken/datasource-api/pkg/server/aws"
	"github.com/ca-risken/datasource-api/proto/aws"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type Server struct {
	port        string
	coreSvcAddr string
	db          *db.Client
	queue       *queue.Client
	logger      logging.Logger
}

func NewServer(port string, coreSvcAddr string, db *db.Client, logger logging.Logger) *Server {
	return &Server{
		port:        port,
		coreSvcAddr: coreSvcAddr,
		db:          db,
		logger:      logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	pjClient := s.newProjectClient(s.coreSvcAddr)
	awsSvc := awsServer.NewAWSSerevice(ctx, s.db, s.queue, pjClient, s.logger)
	hsvc := health.NewServer()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpctrace.UnaryServerInterceptor(),
				mimosarpc.LoggingUnaryServerInterceptor(s.logger))))
	aws.RegisterAWSServiceServer(server, awsSvc)
	grpc_health_v1.RegisterHealthServer(server, hsvc)

	reflection.Register(server) // enable reflection API

	localServerAddr := fmt.Sprintf(":%s", s.port)
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

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := healthCheck(ctx, localServerAddr); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			s.logger.Errorf(ctx, "health check is failed: %w", err)
		} else {
			fmt.Fprintln(w, "ok")
		}
	})

	go func() {
		if err := http.ListenAndServe(":3000", mux); err != http.ErrServerClosed {
			s.logger.Errorf(ctx, "failed to start http server: %w", err)
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

func healthCheck(ctx context.Context, addr string) error {
	conn, err := getGRPCConn(context.Background(), addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	res, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if res.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("returned status is '%v'", res.Status)
	}
	return nil
}

func (s *Server) newProjectClient(svcAddr string) project.ProjectServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		s.logger.Fatalf(ctx, "failed to get grpc connection: err=%+v", err)
	}
	return project.NewProjectServiceClient(conn)
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
