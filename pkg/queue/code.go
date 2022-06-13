package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type CodeQueueAPI interface {
	SendGitleaksMessage(ctx context.Context, msg *message.GitleaksQueueMessage, fullScan bool) (*sqs.SendMessageOutput, error)
}

var _ CodeQueueAPI = (*Client)(nil) // verify interface compliance

func (c *Client) SendGitleaksMessage(ctx context.Context, msg *message.GitleaksQueueMessage, fullScan bool) (*sqs.SendMessageOutput, error) {
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("parse error, err=%w", err)
	}
	url := c.codeGitleaksQueueURL
	if fullScan && c.codeGitleaksFullScanQueueURL != "" {
		url = c.codeGitleaksFullScanQueueURL
	}
	return c.send(ctx, url, &buf)
}
