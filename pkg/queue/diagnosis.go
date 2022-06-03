package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type DiagnosisQueueAPI interface {
	SendWpscanMessage(ctx context.Context, msg *message.WpscanQueueMessage) (*sqs.SendMessageOutput, error)
	SendPortscanMessage(ctx context.Context, msg *message.PortscanQueueMessage) (*sqs.SendMessageOutput, error)
	SendApplicationScanMessage(ctx context.Context, msg *message.ApplicationScanQueueMessage) (*sqs.SendMessageOutput, error)
}

var _ DiagnosisQueueAPI = (*Client)(nil) // verify interface compliance

func (c *Client) SendWpscanMessage(ctx context.Context, msg *message.WpscanQueueMessage) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	return c.send(ctx, c.diagnosisWpscanQueueURL, &buf)
}

func (c *Client) SendPortscanMessage(ctx context.Context, msg *message.PortscanQueueMessage) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	return c.send(ctx, c.diagnosisPortscanQueueURL, &buf)
}

func (c *Client) SendApplicationScanMessage(ctx context.Context, msg *message.ApplicationScanQueueMessage) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	return c.send(ctx, c.diagnosisApplicationScanQueueURL, &buf)
}
