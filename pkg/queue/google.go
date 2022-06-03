package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type GoogleQueueAPI interface {
	SendGoogleAssetMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error)
	SendGoogleCloudSploitMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error)
	SendGoogleSCCMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error)
	SendGooglePortscanMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error)
}

var _ GoogleQueueAPI = (*Client)(nil) // verify interface compliance

func (c *Client) SendGoogleAssetMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendGoogleMessage(ctx, c.googleAssetQueueURL, msg)
}

func (c *Client) SendGoogleCloudSploitMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendGoogleMessage(ctx, c.googleCloudSploitQueueURL, msg)
}

func (c *Client) SendGoogleSCCMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendGoogleMessage(ctx, c.googleSCCQueueURL, msg)
}

func (c *Client) SendGooglePortscanMessage(ctx context.Context, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendGoogleMessage(ctx, c.googlePortscanQueueURL, msg)
}

func (c *Client) sendGoogleMessage(ctx context.Context, url string, msg *message.GCPQueueMessage) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	return c.send(ctx, url, &buf)
}
