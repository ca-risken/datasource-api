package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type AWSQueueAPI interface {
	SendAWSGuardDutyMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error)
	SendAWSAdminCheckerMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error)
	SendAWSAccessAnalyzerMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error)
	SendAWSCloudSploitMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error)
	SendAWSPortscanMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error)
}

var _ AWSQueueAPI = (*Client)(nil) // verify interface compliance

func (c *Client) SendAWSGuardDutyMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendAWSMessage(ctx, c.awsGuardDutyQueueURL, msg)
}

func (c *Client) SendAWSAdminCheckerMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendAWSMessage(ctx, c.awsAdminCheckerQueueURL, msg)
}

func (c *Client) SendAWSAccessAnalyzerMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendAWSMessage(ctx, c.awsAdminCheckerQueueURL, msg)
}

func (c *Client) SendAWSCloudSploitMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendAWSMessage(ctx, c.awsCloudSploitQueueURL, msg)
}

func (c *Client) SendAWSPortscanMessage(ctx context.Context, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendAWSMessage(ctx, c.awsPortscanQueueURL, msg)
}

func (c *Client) sendAWSMessage(ctx context.Context, url string, msg *message.AWSQueueMessage) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	return c.send(ctx, url, &buf)
}
