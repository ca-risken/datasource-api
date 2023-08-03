package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ca-risken/core/proto/alert"
	"github.com/slack-go/slack"
)

const (
	LOCALE_JA      = "ja"
	LOCALE_EN      = "en"
	DEFAULT_LOCALE = LOCALE_EN
	MESSAGE_JA     = `スキャンエラーが発生しました。
エラー原因は複数の可能性があります。エラーメッセージから判断がつかない場合はシステム管理者にお問い合わせください。

- 設定ミスの場合はスキャン設定を修正してください。
- スキャン対象の障害や一時的なネットワークエラーなどが発生している場合(サーバー側の問題)は、しばらく待ってから再度スキャンを実行するか次回のスケール実行までお待ち下さい。
`
	MESSAGE_EN = `Scan error has occurred.
There could be multiple possible reasons for the error. If you cannot determine the cause from the error message, please contact your system administrator.

- If there is a setting mistake, please correct your scan settings.
- If there are issues such as malfunctions of the scan target or temporary network errors (server-side issues), please wait for a while before running the scan again, or wait until the next scale execution.
`
)

type slackNotifySetting struct {
	WebhookURL string `json:"webhook_url"`
	Locale     string `json:"locale"`
}

func (d *DataSourceService) notifyScanError(ctx context.Context, n *alert.Notification, scanErrors *ScanErrors) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(n.NotifySetting), &setting); err != nil {
		return err
	}
	if setting.WebhookURL == "" {
		d.logger.Warnf(ctx, "webhook url is empty: project_id=%d, notification_id=%d", n.ProjectId, n.NotificationId)
		return nil
	}
	locale := DEFAULT_LOCALE
	if setting.Locale != "" {
		locale = setting.Locale
	}
	if err := slack.PostWebhook(setting.WebhookURL, d.getScanErrorPayload(locale, n.ProjectId, scanErrors)); err != nil {
		return fmt.Errorf("failed to send slack: %w", err)
	}
	return nil
}

func (d *DataSourceService) getScanErrorPayload(locale string, projectID uint32, scanErrors *ScanErrors) *slack.WebhookMessage {
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

func (d *DataSourceService) getSlackAttachments(projectID uint32, scanErrors *ScanErrors) []slack.Attachment {
	attachments := []slack.Attachment{}
	for _, aws := range scanErrors.awsErrors {
		attachments = append(attachments, slack.Attachment{
			Color: "warning",
			Fields: []slack.AttachmentField{
				{
					Title: "DataSource",
					Value: fmt.Sprintf("<%s?project_id=%d&from=slack|%s>",
						getDataSourceSettingURL(d.baseURL, aws.DataSource),
						projectID,
						aws.DataSource,
					),
				},
				{
					Title: "ErrorMessage",
					Value: aws.StatusDetail,
				},
			},
		})
	}
	return attachments
}

func getDataSourceSettingURL(baseURL, dataSource string) string {
	switch {
	case strings.HasPrefix(dataSource, "aws:"):
		return fmt.Sprintf("%s/#/aws/data-source", baseURL)
	default:
		return baseURL
	}
}