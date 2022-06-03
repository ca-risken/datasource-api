package message

import (
	"encoding/json"
	"time"

	"github.com/ca-risken/datasource-api/proto/google"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// GoogleAssetDataSource is the DataSource label for Cloud Asset Inventory.
	GoogleAssetDataSource = "google:asset"
	// GoogleCloudSploitDataSource is the DataSource label for Aqua Cloud Sploit.
	GoogleCloudSploitDataSource = "google:cloudsploit"
	// GoogleSCCDataSource is the DataSource label for Security Command Center.
	GoogleSCCDataSource = "google:scc"
	// GooglePortscanDataSource is the DataSource label for Portscan.
	GooglePortscanDataSource = "google:portscan"
)

// GCPQueueMessage is the message for SQS queue
type GCPQueueMessage struct {
	GCPID              uint32 `json:"gcp_id"`
	ProjectID          uint32 `json:"project_id"`
	GoogleDataSourceID uint32 `json:"google_data_source_id"`
	ScanOnly           bool   `json:"scan_only,string"`
}

// Validate is the validation to GuardDutyMessage
func (g *GCPQueueMessage) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GCPID, validation.Required),
		validation.Field(&g.ProjectID, validation.Required),
		validation.Field(&g.GoogleDataSourceID, validation.Required),
	)
}

// ParseMessage parse message & validation
func ParseMessageGCP(msg string) (*GCPQueueMessage, error) {
	message := &GCPQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}

// CutString returns cutting specific `cut` characters with ` ...` suffix from `input` string.
func CutString(input string, cut int) string {
	if len(input) > cut {
		return input[:cut] + " ..." // cut long text
	}
	return input
}

// InitScanStatusGCP return init AttachGCPDataSourceRequest data
func InitScanStatusGCP(g *google.GCPDataSource) *google.AttachGCPDataSourceRequest {
	return &google.AttachGCPDataSourceRequest{
		ProjectId: g.ProjectId,
		GcpDataSource: &google.GCPDataSourceForUpsert{
			GcpId:              g.GcpId,
			GoogleDataSourceId: g.GoogleDataSourceId,
			ProjectId:          g.ProjectId,
			ScanAt:             time.Now().Unix(),
			Status:             google.Status_UNKNOWN, // After scan, will be updated
			StatusDetail:       "",
		},
	}
}
