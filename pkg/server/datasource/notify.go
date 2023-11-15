package datasource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/cenkalti/backoff/v4"
	"github.com/slack-go/slack"
)

const (
	LOCALE_JA      = "ja"
	LOCALE_EN      = "en"
	DEFAULT_LOCALE = LOCALE_EN
	MESSAGE_JA     = `スキャンエラーが発生しました。
エラー原因は複数の可能性があります。エラーメッセージから判断がつかない場合はシステム管理者にお問い合わせください。

- 設定ミスの場合はスキャン設定を修正してください。
- スキャン対象の障害や一時的なネットワークエラーなどが発生している場合(サーバー側の問題)は、しばらく待ってから再度スキャンを実行するか次回のスケジュール実行までお待ち下さい。
`
	MESSAGE_EN = `Scan error has occurred.
There could be multiple possible reasons for the error. If you cannot determine the cause from the error message, please contact your system administrator.

- If there is a setting mistake, please correct your scan settings.
- If there are issues such as malfunctions of the scan target or temporary network errors (server-side issues), please wait for a while before running the scan again, or wait until the next scheduled scan.
`
)

type slackNotifySetting struct {
	WebhookURL string `json:"webhook_url"`
	ChannelID  string `json:"channel_id"`
	Locale     string `json:"locale"`
}

func (d *DataSourceService) notifyScanError(ctx context.Context, n *alert.Notification, scanErrors *ScanErrors) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(n.NotifySetting), &setting); err != nil {
		return err
	}
	if setting.WebhookURL == "" && setting.ChannelID == "" {
		d.logger.Warnf(ctx, "webhook_url and channel_id are empty: project_id=%d, notification_id=%d", n.ProjectId, n.NotificationId)
		return nil
	}
	locale := d.defaultLocale
	if setting.Locale != "" {
		locale = setting.Locale
	}

	if setting.WebhookURL != "" {
		if err := slack.PostWebhook(setting.WebhookURL, d.getScanErrorWebhookMessage(locale, n.ProjectId, scanErrors)); err != nil {
			return fmt.Errorf("failed to send slack(webhook): %w", err)
		}
	} else if setting.ChannelID != "" {
		if err := d.postMessageSlackWithRetry(ctx, setting.ChannelID, d.getScanErrorMessageOpt(locale, n.ProjectId, scanErrors)...); err != nil {
			return fmt.Errorf("failed to send slack(postmessage): %w", err)
		}
	}
	return nil
}

func (d *DataSourceService) getScanErrorWebhookMessage(locale string, projectID uint32, scanErrors *ScanErrors) *slack.WebhookMessage {
	text := MESSAGE_EN
	if locale == LOCALE_JA {
		text = MESSAGE_JA
	}
	msg := slack.WebhookMessage{
		Text:        text,
		Attachments: d.getSlackAttachments(projectID, scanErrors),
	}
	return &msg
}

func (d *DataSourceService) getScanErrorMessageOpt(locale string, projectID uint32, scanErrors *ScanErrors) []slack.MsgOption {
	text := MESSAGE_EN
	if locale == LOCALE_JA {
		text = MESSAGE_JA
	}
	msgOpt := []slack.MsgOption{}
	msgOpt = append(msgOpt, slack.MsgOptionText(text, false))
	msgOpt = append(msgOpt, slack.MsgOptionAttachments(d.getSlackAttachments(projectID, scanErrors)...))
	return msgOpt
}

func (d *DataSourceService) getSlackAttachments(projectID uint32, scanErrors *ScanErrors) []slack.Attachment {
	attachments := []slack.Attachment{}
	for _, aws := range scanErrors.awsErrors {
		attachments = append(attachments, generateSlackAttachment(d.baseURL, aws.DataSource, aws.StatusDetail, projectID))
	}
	for _, gcp := range scanErrors.gcpErrors {
		attachments = append(attachments, generateSlackAttachment(d.baseURL, gcp.DataSource, gcp.StatusDetail, projectID))
	}
	for _, g := range scanErrors.githubErrors {
		attachments = append(attachments, generateSlackAttachment(d.baseURL, g.DataSource, g.StatusDetail, projectID))
	}
	for _, diagnosis := range scanErrors.diagnosisErrors {
		attachments = append(attachments, generateSlackAttachment(d.baseURL, diagnosis.DataSource, diagnosis.StatusDetail, projectID))
	}
	for _, o := range scanErrors.osintErrors {
		attachments = append(attachments, generateSlackAttachment(d.baseURL, o.DataSource, o.StatusDetail, projectID))
	}
	return attachments
}

func generateSlackAttachment(baseURL, dataSource, errorMessage string, projectID uint32) slack.Attachment {
	return slack.Attachment{
		Color: "warning",
		Fields: []slack.AttachmentField{
			{
				Title: "DataSource",
				Value: fmt.Sprintf("<%s?project_id=%d&from=slack|%s>",
					getDataSourceSettingURL(baseURL, dataSource),
					projectID,
					dataSource,
				),
			},
			{
				Title: "ErrorMessage",
				Value: errorMessage,
			},
		},
	}
}

func getDataSourceSettingURL(baseURL, dataSource string) string {
	switch {
	case strings.HasPrefix(dataSource, "aws:"):
		return fmt.Sprintf("%s/aws/data-source", baseURL)
	case strings.HasPrefix(dataSource, "google:"):
		return fmt.Sprintf("%s/google/gcp-data-source", baseURL)
	case strings.HasPrefix(dataSource, "code:"):
		return fmt.Sprintf("%s/code/github", baseURL)
	case dataSource == message.DataSourceNameWPScan:
		return fmt.Sprintf("%s/diagnosis/wpscan", baseURL)
	case dataSource == message.DataSourceNamePortScan:
		return fmt.Sprintf("%s/diagnosis/portscan", baseURL)
	case dataSource == message.DataSourceNameApplicationScan:
		return fmt.Sprintf("%s/diagnosis/applicationscan", baseURL)
	case strings.HasPrefix(dataSource, "osint:"):
		return fmt.Sprintf("%s/osint/data-source", baseURL)
	default:
		return baseURL
	}
}

func (d *DataSourceService) postMessageSlack(channelID string, msg ...slack.MsgOption) error {
	if _, _, err := d.slackClient.PostMessage(channelID, msg...); err != nil {
		var rateLimitError *slack.RateLimitedError
		if errors.As(err, &rateLimitError) {
			time.Sleep(rateLimitError.RetryAfter)
		}
		return err
	}
	return nil
}

func (d *DataSourceService) postMessageSlackWithRetry(ctx context.Context, channelID string, msg ...slack.MsgOption) error {
	operation := func() error {
		return d.postMessageSlack(channelID, msg...)
	}
	return backoff.RetryNotify(operation, d.retryer, d.newRetryLogger(ctx, "postMessageSlack"))
}
