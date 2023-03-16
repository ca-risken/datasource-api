package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/common/pkg/profiler"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/queue"
	"github.com/ca-risken/datasource-api/pkg/server"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "core"
	serviceName = "datasource-api"
)

type AppConf struct {
	Port            string   `default:"8081"`
	Debug           bool     `default:"false"`
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// gRPC
	CoreSvcAddr string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`

	// queue
	AWSRegion   string `envconfig:"aws_region" default:"ap-northeast-1"`
	SQSEndpoint string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	AWSGuardDutyQueueURL             string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-guardduty"`
	AWSAccessAnalyzerQueueURL        string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-accessanalyzer"`
	AWSAdminCheckerQueueURL          string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-adminchecker"`
	AWSCloudSploitQueueURL           string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-cloudsploit"`
	AWSCloudSploitOldQueueURL        string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-cloudsploit-old"`
	AWSPortscanQueueURL              string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/aws-portscan"`
	GoogleAssetQueueURL              string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/google-asset"`
	GoogleCloudSploitQueueURL        string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/google-cloudsploit"`
	GoogleCloudSploitOldQueueURL     string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/google-cloudsploit-old"`
	GoogleSCCQueueURL                string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/google-scc"`
	GooglePortscanQueueURL           string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/google-portscan"`
	CodeGitleaksQueueURL             string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/code-gitleaks"`
	CodeDependencyQueueURL           string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/code-dependency"`
	OSINTSubdomainQueueURL           string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-subdomain"`
	OSINTWebsiteQueueURL             string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-website"`
	DiagnosisWpscanQueueURL          string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-wpscan"`
	DiagnosisPortscanQueueURL        string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-portscan"`
	DiagnosisApplicationScanQueueURL string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/diagnosis-applicationscan"`

	// datasource
	GoogleCredentialPath string `required:"true" split_words:"true" default:"/tmp/credential.json"` // google
	CodeDataKey          string `split_words:"true" required:"true"`                                // code

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
	if conf.Debug {
		logger.Level(logging.DebugLevel)
		logger.Debug(ctx, "Set debug logger")
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
	d, err := db.NewClient(dbConf, logger)
	if err != nil {
		logger.Fatalf(ctx, "Failed to create database client: %w", err)
	}
	queueConf := &queue.SQSConfig{
		AWSRegion:   conf.AWSRegion,
		SQSEndpoint: conf.SQSEndpoint,

		AWSGuardDutyQueueURL:             conf.AWSGuardDutyQueueURL,
		AWSAccessAnalyzerQueueURL:        conf.AWSAccessAnalyzerQueueURL,
		AWSAdminCheckerQueueURL:          conf.AWSAdminCheckerQueueURL,
		AWSCloudSploitQueueURL:           conf.AWSCloudSploitQueueURL,
		AWSCloudSploitOldQueueURL:        conf.AWSCloudSploitOldQueueURL,
		AWSPortscanQueueURL:              conf.AWSPortscanQueueURL,
		GoogleAssetQueueURL:              conf.GoogleAssetQueueURL,
		GoogleCloudSploitQueueURL:        conf.GoogleCloudSploitQueueURL,
		GoogleCloudSploitOldQueueURL:     conf.GoogleCloudSploitOldQueueURL,
		GoogleSCCQueueURL:                conf.GoogleSCCQueueURL,
		GooglePortscanQueueURL:           conf.GooglePortscanQueueURL,
		CodeGitleaksQueueURL:             conf.CodeGitleaksQueueURL,
		CodeDependencyQueueURL:           conf.CodeDependencyQueueURL,
		OSINTSubdomainQueueURL:           conf.OSINTSubdomainQueueURL,
		OSINTWebsiteQueueURL:             conf.OSINTWebsiteQueueURL,
		DiagnosisWpscanQueueURL:          conf.DiagnosisWpscanQueueURL,
		DiagnosisPortscanQueueURL:        conf.DiagnosisPortscanQueueURL,
		DiagnosisApplicationScanQueueURL: conf.DiagnosisApplicationScanQueueURL,
	}
	q, err := queue.NewClient(ctx, queueConf, logger)
	if err != nil {
		logger.Fatalf(ctx, "Failed to create sqs client: %w", err)
	}
	s := server.NewServer(
		conf.Port,
		conf.CoreSvcAddr,
		conf.AWSRegion,
		conf.GoogleCredentialPath,
		conf.CodeDataKey,
		d,
		q,
		logger,
	)

	err = s.Run(ctx)
	if err != nil {
		logger.Fatalf(ctx, "Failed to run server: %w", err)
	}
}

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}
