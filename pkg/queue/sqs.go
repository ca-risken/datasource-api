package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	AWSPortscanQueueURL       string

	// google
	GoogleAssetQueueURL       string
	GoogleCloudSploitQueueURL string
	GoogleSCCQueueURL         string
	GooglePortscanQueueURL    string

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
	svc    *sqs.SQS
	logger logging.Logger

	// aws
	AWSGuardDutyQueueURL      string
	AWSAccessAnalyzerQueueURL string
	AWSAdminCheckerQueueURL   string
	AWSCloudSploitQueueURL    string
	AWSPortscanQueueURL       string

	// google
	GoogleAssetQueueURL       string
	GoogleCloudSploitQueueURL string
	GoogleSCCQueueURL         string
	GooglePortscanQueueURL    string

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

func NewClient(conf *SQSConfig, l logging.Logger) (*Client, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sqs session, err=%w", err)
	}
	sqsClient := sqs.New(sess, &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.SQSEndpoint,
	})
	return &Client{
		svc:    sqsClient,
		logger: l,

		AWSGuardDutyQueueURL:             conf.AWSGuardDutyQueueURL,
		AWSAccessAnalyzerQueueURL:        conf.AWSAccessAnalyzerQueueURL,
		AWSAdminCheckerQueueURL:          conf.AWSAdminCheckerQueueURL,
		AWSCloudSploitQueueURL:           conf.AWSCloudSploitQueueURL,
		AWSPortscanQueueURL:              conf.AWSPortscanQueueURL,
		GoogleAssetQueueURL:              conf.GoogleAssetQueueURL,
		GoogleCloudSploitQueueURL:        conf.GoogleCloudSploitQueueURL,
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
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	resp, err := c.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
