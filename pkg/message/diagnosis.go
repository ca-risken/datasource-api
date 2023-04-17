package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	DataSourceNameWPScan          = "diagnosis:wpscan"
	DataSourceNamePortScan        = "diagnosis:portscan"
	DataSourceNameApplicationScan = "diagnosis:application-scan"
)

// WpscanQueueMessage is the message for SQS queue for Wpscan
type WpscanQueueMessage struct {
	DataSource      string `json:"data_source"`
	WpscanSettingID uint32 `json:"wpscan_setting_id"`
	ProjectID       uint32 `json:"project_id"`
	TargetURL       string `json:"target_url"`
	Options         string `json:"options"`
	ScanOnly        bool   `json:"scan_only,string"`
}

// PortscanQueueMessage is the message for SQS queue for Portscan
type PortscanQueueMessage struct {
	DataSource        string `json:"data_source"`
	PortscanSettingID uint32 `json:"portscan_setting_id"`
	PortscanTargetID  uint32 `json:"portscan_target_id"`
	ProjectID         uint32 `json:"project_id"`
	Target            string `json:"target"`
	ScanOnly          bool   `json:"scan_only,string"`
}

// ApplicationScanQueueMessage is the message for SQS queue for ApplicationScan
type ApplicationScanQueueMessage struct {
	DataSource          string `json:"data_source"`
	ApplicationScanID   uint32 `json:"application_scan_id"`
	ProjectID           uint32 `json:"project_id"`
	Name                string `json:"name"`
	ApplicationScanType string `json:"application_scan_type"`
	ScanOnly            bool   `json:"scan_only,string"`
}

func (m *WpscanQueueMessage) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.WpscanSettingID, validation.Required),
		validation.Field(&m.ProjectID, validation.Required),
		validation.Field(&m.TargetURL, validation.Required),
	)
}

func (m *PortscanQueueMessage) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.PortscanSettingID, validation.Required),
		validation.Field(&m.PortscanTargetID, validation.Required),
		validation.Field(&m.ProjectID, validation.Required),
		validation.Field(&m.Target, validation.Required),
	)
}

func (m *ApplicationScanQueueMessage) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.ApplicationScanID, validation.Required),
		validation.Field(&m.ProjectID, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.ApplicationScanType, validation.Required),
	)
}

// ParseWpscanMessage parse wpscan message
func ParseWpscanMessage(msg string) (*WpscanQueueMessage, error) {
	message := &WpscanQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}

// ParsePortscanMessage parse portscan message
func ParsePortscanMessage(msg string) (*PortscanQueueMessage, error) {
	message := &PortscanQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}

// ParseApplicationScanMessage parse applicationscan message
func ParseApplicationScanMessage(msg string) (*ApplicationScanQueueMessage, error) {
	message := &ApplicationScanQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
