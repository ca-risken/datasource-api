package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

// RemediationProposalQueueMessage is the message for SQS queue
type RemediationProposalQueueMessage struct {
	RemediationProposalID uint32 `json:"remediation_proposal_id"`
	FindingID             uint64 `json:"finding_id"`
	ProjectID             uint32 `json:"project_id"`
	AssumeRoleArn         string `json:"assume_role_arn"`
	ExternalID            string `json:"external_id"`
}

// Validate is the validation to RemediationProposalQueueMessage
func (a *RemediationProposalQueueMessage) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.RemediationProposalID, validation.Required),
		validation.Field(&a.FindingID, validation.Required),
		validation.Field(&a.ProjectID, validation.Required),
		validation.Field(&a.AssumeRoleArn, validation.Required),
		validation.Field(&a.ExternalID, validation.Required),
	)
}

// ParseMessageRemediationProposal parse message & validation
func ParseMessageRemediationProposal(msg string) (*RemediationProposalQueueMessage, error) {
	message := &RemediationProposalQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
