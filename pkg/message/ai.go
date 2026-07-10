package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

// AIRemediationQueueMessage is the message for SQS queue
type AIRemediationQueueMessage struct {
	RemediationProposalID uint32 `json:"remediation_proposal_id"`
	FindingID             uint64 `json:"finding_id"`
	ProjectID             uint32 `json:"project_id"`
	DataSource            string `json:"data_source"`
	AWSID                 uint32 `json:"aws_id"`
	AccountID             string `json:"account_id"`
	AssumeRoleArn         string `json:"assume_role_arn"`
	ExternalID            string `json:"external_id"`
}

// Validate is the validation to AIRemediationQueueMessage
func (a *AIRemediationQueueMessage) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.RemediationProposalID, validation.Required),
		validation.Field(&a.FindingID, validation.Required),
		validation.Field(&a.ProjectID, validation.Required),
		validation.Field(&a.DataSource, validation.Required, validation.In(
			AWSAccessAnalyzerDataSource,
			AWSAdminCheckerDataSource,
			AWSCloudSploitDataSource,
			AWSPortscanDataSource,
		)),
		validation.Field(&a.AWSID, validation.Required),
		validation.Field(&a.AccountID, validation.Required, validation.Length(12, 12)),
		validation.Field(&a.AssumeRoleArn, validation.Required),
	)
}

// ParseMessageAIRemediation parse message & validation
func ParseMessageAIRemediation(msg string) (*AIRemediationQueueMessage, error) {
	message := &AIRemediationQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
