package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// GitleaksDataSource is the specific data_source label for gitleaks
	GitleaksDataSource = "code:gitleaks"
)

// GitHubQueueMessage is the message for SQS queue
type GitHubQueueMessage struct {
	GitHubSettingID uint32 `json:"github_setting_id"`
	ProjectID       uint32 `json:"project_id"`
	ScanOnly        bool   `json:"scan_only,string"`
}

// Validate is the validation to GuardDutyMessage
func (g *GitHubQueueMessage) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GitHubSettingID, validation.Required),
		validation.Field(&g.ProjectID, validation.Required),
	)
}

// ParseMessage parse message & validation
func ParseMessageGitHub(msg string) (*GitHubQueueMessage, error) {
	message := &GitHubQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
