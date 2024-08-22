package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// AzureProwler is the specific data_source label for prowler
	AzureProwlerDataSource = "azure:prowler"
)

// AzureQueueMessage is the message for SQS queue
type AzureQueueMessage struct {
	AzureID           uint32 `json:"azure_id"`
	AzureDataSourceID uint32 `json:"azure_data_source_id"`
	DataSource        string `json:"data_source"`
	ProjectID         uint32 `json:"project_id"`
	SubscriptionID    string `json:"account_id"`
	VerificationCode  string `json:"verification_code"`
	ScanOnly          bool   `json:"scan_only,string"`
}

// Validate is the validation to ProwlerMessage
func (g *AzureQueueMessage) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.AzureID, validation.Required),
		validation.Field(&g.AzureDataSourceID, validation.Required),
		validation.Field(&g.DataSource, validation.Required, validation.In(
			AzureProwlerDataSource,
		)),
		validation.Field(&g.ProjectID, validation.Required),
		validation.Field(&g.SubscriptionID, validation.Required, validation.Length(36, 36)),
	)
}

// ParseMessage parse message & validation
func ParseMessageAzure(msg string) (*AzureQueueMessage, error) {
	message := &AzureQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
