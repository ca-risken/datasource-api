package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// GitleaksDataSource is the specific data_source label for gitleaks
	GitleaksDataSource = "code:gitleaks"
	// DependencyDataSource is the specific data_source label for dependency
	DependencyDataSource = "code:dependency"
	// CodeScanDataSource is the specific data_source label for codescan
	CodeScanDataSource = "code:codescan"
)

type RepositoryMetadata struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Visibility    string `json:"visibility"`
	Archived      bool   `json:"archived"`
	Fork          bool   `json:"fork"`
	Disabled      bool   `json:"disabled"`
	Size          int64  `json:"size"`
	CreatedAt     int64  `json:"created_at"`
	PushedAt      int64  `json:"pushed_at"`
	HTMLURL       string `json:"html_url"`
}

// CodeQueueMessage is the message for SQS queue
type CodeQueueMessage struct {
	GitHubSettingID uint32              `json:"github_setting_id"`
	ProjectID       uint32              `json:"project_id"`
	ScanOnly        bool                `json:"scan_only,string"`
	FullScan        bool                `json:"full_scan,string"`
	RepositoryName  string              `json:"repository_name"`
	Repository      *RepositoryMetadata `json:"repository"`
}

// Validate is the validation to GuardDutyMessage
func (g *CodeQueueMessage) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GitHubSettingID, validation.Required),
		validation.Field(&g.ProjectID, validation.Required),
	)
}

// ParseMessage parse message & validation
func ParseMessageGitHub(msg string) (*CodeQueueMessage, error) {
	message := &CodeQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
