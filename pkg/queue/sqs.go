package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
)

type SQSConfig struct {
	AWSRegion   string
	SQSEndpoint string

	// aws
	AWSGuardDutyQueueURL      string
	AWSAccessAnalyzerQueueURL string
	AWSAdminCheckerQueueURL   string
	AWSCloudSploitQueueURL    string
	AWSCloudSploitOldQueueURL string
	AWSPortscanQueueURL       string

	// google
	GoogleAssetQueueURL          string
	GoogleCloudSploitQueueURL    string
	GoogleCloudSploitOldQueueURL string
	GoogleSCCQueueURL            string
	GooglePortscanQueueURL       string

	// code
	CodeGitleaksQueueURL   string
	CodeDependencyQueueURL string

	// osint
	OSINTSubdomainQueueURL string
	OSINTWebsiteQueueURL   string

	// diagnosis
	DiagnosisWpscanQueueURL          string
	DiagnosisPortscanQueueURL        string
	DiagnosisApplicationScanQueueURL string
}

type Client struct {
	svc    *sqs.Client
	logger logging.Logger

	// aws
	AWSGuardDutyQueueURL      string
	AWSAccessAnalyzerQueueURL string
	AWSAdminCheckerQueueURL   string
	AWSCloudSploitQueueURL    string
	AWSCloudSploitOldQueueURL string
	AWSPortscanQueueURL       string

	// google
	GoogleAssetQueueURL          string
	GoogleCloudSploitQueueURL    string
	GoogleCloudSploitOldQueueURL string
	GoogleSCCQueueURL            string
	GooglePortscanQueueURL       string

	// code
	CodeGitleaksQueueURL   string
	CodeDependencyQueueURL string

	// osint
	OSINTSubdomainQueueURL string
	OSINTWebsiteQueueURL   string

	// diagnosis
	DiagnosisWpscanQueueURL          string
	DiagnosisPortscanQueueURL        string
	DiagnosisApplicationScanQueueURL string
}

func NewClient(ctx context.Context, conf *SQSConfig, l logging.Logger) (*Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service != sqs.ServiceID {
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		}

		ep := aws.Endpoint{
			PartitionID:   "aws",
			SigningRegion: region,
		}
		if conf.SQSEndpoint != "" {
			ep.URL = conf.SQSEndpoint
		}
		if conf.AWSRegion != "" {
			ep.SigningRegion = conf.AWSRegion
		}
		return ep, nil
	})
	l.Debugf(ctx, "SQS endpoint: region=%s, url=%s", conf.AWSRegion, conf.SQSEndpoint)
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, fmt.Errorf("load aws configuration error, err=%+w", err)
	}
	return &Client{
		svc:    sqs.NewFromConfig(cfg),
		logger: l,

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
	}, nil
}

func (c *Client) Send(ctx context.Context, url string, msg interface{}) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("SQS message parse error, err=%w", err)
	}
	resp, err := c.svc.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: 1,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
