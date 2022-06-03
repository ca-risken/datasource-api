package queue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	GitleaksQueueURL         string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/code-gitleaks"`
	GitleaksFullScanQueueURL string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/code-gitleaks"`

	// osint
	SubdomainQueueURL string
	WebsiteQueueURL   string

	// diagnosis
	DiagnosisWpscanQueueURL          string
	DiagnosisPortscanQueueURL        string
	DiagnosisApplicationScanQueueURL string
}

type Client struct {
	svc *sqs.SQS

	// aws
	awsGuardDutyQueueURL      string
	awsAccessAnalyzerQueueURL string
	awsAdminCheckerQueueURL   string
	awsCloudSploitQueueURL    string
	awsPortscanQueueURL       string

	// google
	googleAssetQueueURL       string
	googleCloudSploitQueueURL string
	googleSCCQueueURL         string
	googlePortscanQueueURL    string

	// code
	gitleaksQueueURL         string
	gitleaksFullScanQueueURL string

	// osint
	subdomainQueueURL string
	websiteQueueURL   string

	// diagnosis
	diagnosisWpscanQueueURL          string
	diagnosisPortscanQueueURL        string
	diagnosisApplicationScanQueueURL string
}

func NewSQSClient(ctx context.Context, conf *SQSConfig) (*Client, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("create session error, %w", err)
	}
	session := sqs.New(sess, &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.SQSEndpoint,
	})
	return &Client{
		svc:                              session,
		awsGuardDutyQueueURL:             conf.AWSGuardDutyQueueURL,
		awsAccessAnalyzerQueueURL:        conf.AWSAccessAnalyzerQueueURL,
		awsAdminCheckerQueueURL:          conf.AWSAdminCheckerQueueURL,
		awsCloudSploitQueueURL:           conf.AWSCloudSploitQueueURL,
		awsPortscanQueueURL:              conf.AWSPortscanQueueURL,
		googleAssetQueueURL:              conf.GoogleAssetQueueURL,
		googleCloudSploitQueueURL:        conf.GoogleCloudSploitQueueURL,
		googleSCCQueueURL:                conf.GoogleSCCQueueURL,
		googlePortscanQueueURL:           conf.GooglePortscanQueueURL,
		gitleaksQueueURL:                 conf.GitleaksQueueURL,
		gitleaksFullScanQueueURL:         conf.GitleaksFullScanQueueURL,
		diagnosisWpscanQueueURL:          conf.DiagnosisWpscanQueueURL,
		diagnosisPortscanQueueURL:        conf.DiagnosisPortscanQueueURL,
		diagnosisApplicationScanQueueURL: conf.DiagnosisApplicationScanQueueURL,
	}, nil
}

func (c *Client) send(ctx context.Context, url string, buf *[]byte) (*sqs.SendMessageOutput, error) {
	resp, err := c.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody:  aws.String(string(*buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
