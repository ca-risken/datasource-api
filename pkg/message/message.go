package message

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
