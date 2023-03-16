package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// AWSGuardDutyDataSource is the specific data_source label for guard-duty
	AWSGuardDutyDataSource = "aws:guard-duty"
	// AWSAccessAnalyzerDataSource is the specific data_source label for access-analyzer
	AWSAccessAnalyzerDataSource = "aws:access-analyzer"
	// AWSAdminCheckerDataSource is the specific data_source label for admin-checker
	AWSAdminCheckerDataSource = "aws:admin-checker"
	// AWSCloudSploitDataSource is the specific data_source label for cloudsploit
	AWSCloudSploitDataSource = "aws:cloudsploit"
	// AWSPortscanDataSource is the specific data_source label for portscan
	AWSPortscanDataSource = "aws:portscan"
)

// AWSQueueMessage is the message for SQS queue
type AWSQueueMessage struct {
	AWSID           uint32 `json:"aws_id"`
	AWSDataSourceID uint32 `json:"aws_data_source_id"`
	DataSource      string `json:"data_source"`
	ProjectID       uint32 `json:"project_id"`
	AccountID       string `json:"account_id"`
	AssumeRoleArn   string `json:"assume_role_arn"`
	ExternalID      string `json:"external_id"`
	ScanOnly        bool   `json:"scan_only,string"`
}

// Validate is the validation to GuardDutyMessage
func (g *AWSQueueMessage) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.AWSID, validation.Required),
		validation.Field(&g.AWSDataSourceID, validation.Required),
		validation.Field(&g.DataSource, validation.Required, validation.In(
			AWSGuardDutyDataSource,
			AWSAccessAnalyzerDataSource,
			AWSAdminCheckerDataSource,
			AWSCloudSploitDataSource,
			AWSPortscanDataSource,
		)),
		validation.Field(&g.ProjectID, validation.Required),
		validation.Field(&g.AccountID, validation.Required, validation.Length(12, 12)),
	)
}

// ParseMessage parse message & validation
func ParseMessageAWS(msg string) (*AWSQueueMessage, error) {
	message := &AWSQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
