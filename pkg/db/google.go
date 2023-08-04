package db

import (
	"context"
	"time"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/google"
	"github.com/vikyd/zero"
)

type GoogleRepoInterface interface {
	// google_data_source
	ListGoogleDataSource(ctx context.Context, googleDataSourceID uint32, name string) (*[]model.GoogleDataSource, error)

	// gcp
	ListGCP(ctx context.Context, projectID, gcpID uint32, gcpProjectID string) (*[]model.GCP, error)
	GetGCP(ctx context.Context, projectID, gcpID uint32) (*model.GCP, error)
	UpsertGCP(ctx context.Context, gcp *google.GCPForUpsert) (*model.GCP, error)
	DeleteGCP(ctx context.Context, projectID uint32, gcpID uint32) error

	// gcp_data_source
	ListGCPDataSource(ctx context.Context, projectID, gcpID uint32) (*[]GCPDataSource, error)
	GetGCPDataSource(ctx context.Context, projectID, gcpID, googleDataSourceID uint32) (*GCPDataSource, error)
	UpsertGCPDataSource(ctx context.Context, gcpDataSource *google.GCPDataSourceForUpsert) (*GCPDataSource, error)
	DeleteGCPDataSource(ctx context.Context, projectID, gcpID, googleDataSourceID uint32) error
	ListGCPDataSourceByDataSourceID(ctx context.Context, googleDataSourceID uint32) (*[]GCPDataSource, error)

	// scan_error
	ListGCPScanErrorForNotify(ctx context.Context) ([]*GCPScanError, error)
	UpdateGCPErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, gcpID, googleDataSourceID, projectID uint32) error
}

var _ GoogleRepoInterface = (*Client)(nil) // verify interface compliance

func (c *Client) ListGoogleDataSource(ctx context.Context, googleDataSourceID uint32, name string) (*[]model.GoogleDataSource, error) {
	query := `select * from google_data_source where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(googleDataSourceID) {
		query += " and google_data_source_id = ?"
		params = append(params, googleDataSourceID)
	}
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	data := []model.GoogleDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetGoogleDataSource string = "select * from google_data_source where google_data_source_id=?"

func (c *Client) GetGoogleDataSource(ctx context.Context, googleDataSourceID uint32) (*model.GCP, error) {
	data := model.GCP{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectGetGoogleDataSource, googleDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListGCP(ctx context.Context, projectID, gcpID uint32, gcpProjectID string) (*[]model.GCP, error) {
	query := `select * from gcp where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(gcpID) {
		query += " and gcp_id = ?"
		params = append(params, gcpID)
	}
	if !zero.IsZeroVal(gcpProjectID) {
		query += " and gcp_project_id = ?"
		params = append(params, gcpProjectID)
	}
	data := []model.GCP{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetGCP string = `select * from gcp where project_id=? and gcp_id=?`

func (c *Client) GetGCP(ctx context.Context, projectID, gcpID uint32) (*model.GCP, error) {
	data := model.GCP{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectGetGCP, projectID, gcpID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertGCP = `
INSERT INTO gcp (
  gcp_id,
  name,
  project_id,
  gcp_project_id,
  verification_code
)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  project_id=VALUES(project_id),
  gcp_project_id=VALUES(gcp_project_id),
  verification_code=VALUES(verification_code)
`

func (c *Client) UpsertGCP(ctx context.Context, gcp *google.GCPForUpsert) (*model.GCP, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertGCP,
		gcp.GcpId,
		convertZeroValueToNull(gcp.Name),
		gcp.ProjectId,
		gcp.GcpProjectId,
		gcp.VerificationCode,
	).Error; err != nil {
		return nil, err
	}
	return c.GetGCPByUniqueIndex(ctx, gcp.ProjectId, gcp.GcpProjectId)
}

const selectGetGCPByUniqueIndex string = `select * from gcp where project_id=? and gcp_project_id=?`

func (c *Client) GetGCPByUniqueIndex(ctx context.Context, projectID uint32, gcpProjectID string) (*model.GCP, error) {
	data := model.GCP{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetGCPByUniqueIndex, projectID, gcpProjectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteGCP string = `delete from gcp where project_id=? and gcp_id=?`

func (c *Client) DeleteGCP(ctx context.Context, projectID, gcpID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteGCP, projectID, gcpID).Error; err != nil {
		return err
	}
	return nil
}

type GCPDataSource struct {
	GCPID              uint32 `gorm:"primary_key column:gcp_id"`
	GoogleDataSourceID uint32 `gorm:"primary_key"`
	ProjectID          uint32
	SpecificVersion    string
	Status             string
	StatusDetail       string
	ScanAt             time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Name               string  // google_data_source.name
	Description        string  // google_data_source.description
	MaxScore           float32 // google_data_source.max_score
	GCPProjectID       string  // gcp.gcp_project_id
	ErrorNotifiedAt    time.Time
}

func (c *Client) ListGCPDataSource(ctx context.Context, projectID, gcpID uint32) (*[]GCPDataSource, error) {
	query := `
select
  gds.*, google.name, google.max_score, google.description, gcp.gcp_project_id
from
  gcp_data_source gds
  inner join google_data_source google using(google_data_source_id)
  inner join gcp using(gcp_id, project_id)
where
	1=1
`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and gds.project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(gcpID) {
		query += " and gds.gcp_id = ?"
		params = append(params, gcpID)
	}
	data := []GCPDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetGCPDataSource string = `
select
  gds.*, google.name, google.max_score, google.description, gcp.gcp_project_id, gds.error_notified_at
from
  gcp_data_source gds
  inner join google_data_source google using(google_data_source_id)
  inner join gcp using(gcp_id, project_id)
where
	gds.project_id=? and gds.gcp_id=? and gds.google_data_source_id=?
`

func (c *Client) GetGCPDataSource(ctx context.Context, projectID, gcpID, googleDataSourceID uint32) (*GCPDataSource, error) {
	data := GCPDataSource{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetGCPDataSource, projectID, gcpID, googleDataSourceID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertGCPDataSource string = `
INSERT INTO gcp_data_source (
  gcp_id,
  google_data_source_id,
  project_id,
  specific_version,
  status,
  status_detail,
  scan_at
)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  project_id=VALUES(project_id),
  specific_version=VALUES(specific_version),
  status=VALUES(status),
  status_detail=VALUES(status_detail),
  scan_at=VALUES(scan_at)
`

func (c *Client) UpsertGCPDataSource(ctx context.Context, gcpDataSource *google.GCPDataSourceForUpsert) (*GCPDataSource, error) {
	// Check master table exists
	if _, err := c.GetGoogleDataSource(ctx, gcpDataSource.GoogleDataSourceId); err != nil {
		c.logger.Errorf(ctx, "Not exists google_data_source or DB error: google_data_source_id=%d", gcpDataSource.GoogleDataSourceId)
		return nil, err
	}
	if _, err := c.GetGCP(ctx, gcpDataSource.ProjectId, gcpDataSource.GcpId); err != nil {
		c.logger.Errorf(ctx, "Not exists gcp or DB error: gcp_id=%d", gcpDataSource.GcpId)
		return nil, err
	}

	// Upsert
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertGCPDataSource,
		gcpDataSource.GcpId,
		gcpDataSource.GoogleDataSourceId,
		gcpDataSource.ProjectId,
		gcpDataSource.SpecificVersion,
		gcpDataSource.Status.String(),
		convertZeroValueToNull(gcpDataSource.StatusDetail),
		time.Unix(gcpDataSource.ScanAt, 0),
	).Error; err != nil {
		return nil, err
	}
	return c.GetGCPDataSource(ctx, gcpDataSource.ProjectId, gcpDataSource.GcpId, gcpDataSource.GoogleDataSourceId)
}

const deleteGCPDataSource string = `delete from gcp_data_source where project_id=? and gcp_id=? and google_data_source_id=?`

func (c *Client) DeleteGCPDataSource(ctx context.Context, projectID, gcpID, googleDataSourceID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteGCPDataSource, projectID, gcpID, googleDataSourceID).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListGCPDataSourceByDataSourceID(ctx context.Context, googleDataSourceID uint32) (*[]GCPDataSource, error) {
	query := `
select
  gds.*, google.name, google.max_score, google.description, gcp.gcp_project_id
from
  gcp_data_source gds
  inner join google_data_source google using(google_data_source_id)
  inner join gcp using(gcp_id, project_id)`
	var params []interface{}
	if !zero.IsZeroVal(googleDataSourceID) {
		query += " where google_data_source_id = ?"
		params = append(params, googleDataSourceID)
	}
	data := []GCPDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

type GCPScanError struct {
	GCPID              uint32
	GoogleDataSourceID uint32
	ProjectID          uint32
	DataSource         string
	StatusDetail       string
}

const selectListGCPScanError = `
select
  gds.gcp_id, gds.google_data_source_id, g.name as data_source, gds.project_id, gds.status_detail
from
  gcp_data_source gds 
  inner join google_data_source g using(google_data_source_id) 
where
  gds.status = 'ERROR'
  and gds.error_notified_at is null
`

func (c *Client) ListGCPScanErrorForNotify(ctx context.Context) ([]*GCPScanError, error) {
	data := []*GCPScanError{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectListGCPScanError).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const updateGCPErrorNotifiedAt = `update gcp_data_source set error_notified_at = ? where gcp_id = ? and google_data_source_id = ? and project_id = ?`

func (c *Client) UpdateGCPErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, gcpID, googleDataSourceID, projectID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(updateGCPErrorNotifiedAt, errNotifiedAt, gcpID, googleDataSourceID, projectID).Error; err != nil {
		return err
	}
	return nil
}
