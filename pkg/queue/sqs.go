package queue

import (
	"context"

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
	GitleaksQueueURL         string
	GitleaksFullScanQueueURL string

	// osint
	SubdomainQueueURL string
	WebsiteQueueURL   string

	// diagnosis
	DiagnosisWpscanQueueURL          string
	DiagnosisPortscanQueueURL        string
	DiagnosisApplicationScanQueueURL string
}

type Client struct {
	svc    *sqs.SQS
	logger logging.Logger

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

func NewSQSClient(ctx context.Context, conf *SQSConfig, l logging.Logger) *Client {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		l.Fatalf(ctx, "Failed to create sqs session, err=%w", err)
	}
	session := sqs.New(sess, &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.SQSEndpoint,
	})
	return &Client{
		svc:    session,
		logger: l,

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
		subdomainQueueURL:                conf.SubdomainQueueURL,
		websiteQueueURL:                  conf.WebsiteQueueURL,
		diagnosisWpscanQueueURL:          conf.DiagnosisWpscanQueueURL,
		diagnosisPortscanQueueURL:        conf.DiagnosisPortscanQueueURL,
		diagnosisApplicationScanQueueURL: conf.DiagnosisApplicationScanQueueURL,
	}
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
