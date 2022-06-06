package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/common/pkg/profiler"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/server"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "datasource-api"
	serviceName = "datasource-api"
)

type AppConf struct {
	Port            string   `default:"8081"`
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// gRPC
	CoreSvcAddr string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`

	// queue
	AWSRegion   string `envconfig:"aws_region" default:"ap-northeast-1"`
	SQSEndpoint string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	GuardDutyQueueURL      string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-guardduty"`
	AccessAnalyzerQueueURL string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-accessanalyzer"`
	AdminCheckerQueueURL   string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-adminchecker"`
	CloudsploitQueueURL    string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-cloudsploit"`
	PortscanQueueURL       string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-portscan"`

	// db
	DBMasterHost     string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBMasterUser     string `split_words:"true" default:"hoge"`
	DBMasterPassword string `split_words:"true" default:"moge"`
	DBSlaveHost      string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBSlaveUser      string `split_words:"true" default:"hoge"`
	DBSlavePassword  string `split_words:"true" default:"moge"`
	DBSchema         string `required:"true"    default:"mimosa"`
	DBPort           int    `required:"true"    default:"3306"`
	DBLogMode        bool   `split_words:"true" default:"false"`
	DBMaxConnection  int    `split_words:"true" default:"10"`
}

func main() {
	ctx := context.Background()
	var logger = logging.NewLogger()
	var conf AppConf
	err := envconfig.Process("", &conf)
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: conf.EnvName,
		Debug:       conf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	dbConf := &db.Config{
		MasterHost:     conf.DBMasterHost,
		MasterUser:     conf.DBMasterUser,
		MasterPassword: conf.DBMasterPassword,
		SlaveHost:      conf.DBSlaveHost,
		SlaveUser:      conf.DBSlaveUser,
		SlavePassword:  conf.DBSlavePassword,
		Schema:         conf.DBSchema,
		Port:           conf.DBPort,
		LogMode:        conf.DBLogMode,
		MaxConnection:  conf.DBMaxConnection,
	}
	db := db.NewClient(dbConf, logger)
	server := server.NewServer(conf.Port, conf.CoreSvcAddr, db, logger)

	err = server.Run(ctx)
	if err != nil {
		logger.Fatalf(ctx, "failed to run server: %w", err)
	}
}

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}
