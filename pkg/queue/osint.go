package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type OSINTQueueAPI interface {
	SendOSINTSubdomainMessage(ctx context.Context, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error)
	SendOSINTWebsiteMessage(ctx context.Context, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error)
}

var _ OSINTQueueAPI = (*Client)(nil) // verify interface compliance

func (c *Client) SendOSINTSubdomainMessage(ctx context.Context, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendOSINTMessage(ctx, c.osintSubdomainQueueURL, msg)
}

func (c *Client) SendOSINTWebsiteMessage(ctx context.Context, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error) {
	return c.sendOSINTMessage(ctx, c.osintWebsiteQueueURL, msg)
}

func (c *Client) sendOSINTMessage(ctx context.Context, url string, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	return c.send(ctx, url, &buf)
}
