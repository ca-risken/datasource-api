package db

import (
	"context"
	"time"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/azure"
	"github.com/vikyd/zero"
)

type AzureRepoInterface interface {
	// azure_data_source
	ListAzureDataSource(ctx context.Context, azureDataSourceID uint32, name string) (*[]model.AzureDataSource, error)

	// azure
	ListAzure(ctx context.Context, projectID, azureID uint32, azureProjectID string) (*[]model.Azure, error)
	GetAzure(ctx context.Context, projectID, azureID uint32) (*model.Azure, error)
	UpsertAzure(ctx context.Context, azure *azure.AzureForUpsert) (*model.Azure, error)
	DeleteAzure(ctx context.Context, projectID uint32, azureID uint32) error

	// azure_data_source
	ListRelAzureDataSource(ctx context.Context, projectID, azureID uint32) (*[]RelAzureDataSource, error)
	GetRelAzureDataSource(ctx context.Context, projectID, azureID, azureDataSourceID uint32) (*RelAzureDataSource, error)
	UpsertRelAzureDataSource(ctx context.Context, relAzureDataSource *azure.RelAzureDataSourceForUpsert) (*RelAzureDataSource, error)
	DeleteRelAzureDataSource(ctx context.Context, projectID, azureID, azureDataSourceID uint32) error
	ListRelAzureDataSourceByDataSourceID(ctx context.Context, azureDataSourceID uint32) (*[]RelAzureDataSource, error)

	// scan_error
	ListAzureScanErrorForNotify(ctx context.Context) ([]*AzureScanError, error)
	UpdateAzureErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, azureID, azureDataSourceID, projectID uint32) error
}

var _ AzureRepoInterface = (*Client)(nil) // verify interface compliance

func (c *Client) ListAzureDataSource(ctx context.Context, azureDataSourceID uint32, name string) (*[]model.AzureDataSource, error) {
	query := `select * from azure_data_source where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(azureDataSourceID) {
		query += " and azure_data_source_id = ?"
		params = append(params, azureDataSourceID)
	}
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	data := []model.AzureDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAzureDataSource string = "select * from azure_data_source where azure_data_source_id=?"

func (c *Client) GetAzureDataSource(ctx context.Context, azureDataSourceID uint32) (*model.Azure, error) {
	data := model.Azure{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectGetAzureDataSource, azureDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListAzure(ctx context.Context, projectID, azureID uint32, azureProjectID string) (*[]model.Azure, error) {
	query := `select * from azure where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(azureID) {
		query += " and azure_id = ?"
		params = append(params, azureID)
	}
	if !zero.IsZeroVal(azureProjectID) {
		query += " and subscription_id = ?"
		params = append(params, azureProjectID)
	}
	data := []model.Azure{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAzure string = `select * from azure where project_id=? and azure_id=?`

func (c *Client) GetAzure(ctx context.Context, projectID, azureID uint32) (*model.Azure, error) {
	data := model.Azure{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectGetAzure, projectID, azureID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertAzure = `
INSERT INTO azure (
  azure_id,
  name,
  project_id,
  subscription_id,
  verification_code
)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  project_id=VALUES(project_id),
  subscription_id=VALUES(subscription_id),
  verification_code=VALUES(verification_code)
`

func (c *Client) UpsertAzure(ctx context.Context, azure *azure.AzureForUpsert) (*model.Azure, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertAzure,
		azure.AzureId,
		convertZeroValueToNull(azure.Name),
		azure.ProjectId,
		azure.SubscriptionId,
		azure.VerificationCode,
	).Error; err != nil {
		return nil, err
	}
	return c.GetAzureBySubscriptionID(ctx, azure.ProjectId, azure.SubscriptionId)
}

const selectGetAzureBySubscriptionID string = `select * from azure where project_id=? and subscription_id=?`

func (c *Client) GetAzureBySubscriptionID(ctx context.Context, projectID uint32, subscriptionID string) (*model.Azure, error) {
	data := model.Azure{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetAzureBySubscriptionID, projectID, subscriptionID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteAzure string = `delete from azure where project_id=? and azure_id=?`

func (c *Client) DeleteAzure(ctx context.Context, projectID, azureID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteAzure, projectID, azureID).Error; err != nil {
		return err
	}
	return nil
}

type RelAzureDataSource struct {
	AzureID           uint32 `gorm:"primary_key column:azure_id"`
	AzureDataSourceID uint32 `gorm:"primary_key"`
	ProjectID         uint32
	Status            string
	StatusDetail      string
	ScanAt            time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Name              string  // azure_data_source.name
	Description       string  // azure_data_source.description
	MaxScore          float32 // azure_data_source.max_score
	SubscriptionID    string  // azure.subscription_id
	ErrorNotifiedAt   time.Time
}

const selectListRelAzureDataSource string = `
select
  rads.*, azure.name, ads.max_score, ads.description, azure.subscription_id
from
  rel_azure_data_source rads
  inner join azure_data_source ads using(azure_data_source_id)
  inner join azure using(azure_id, project_id)
where
	1=1
`

func (c *Client) ListRelAzureDataSource(ctx context.Context, projectID, azureID uint32) (*[]RelAzureDataSource, error) {
	query := selectListRelAzureDataSource
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and rads.project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(azureID) {
		query += " and rads.azure_id = ?"
		params = append(params, azureID)
	}
	data := []RelAzureDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetRelAzureDataSource string = `
select
  rads.*, azure.name, ads.max_score, ads.description, azure.subscription_id, rads.error_notified_at
from
  rel_azure_data_source rads
  inner join azure_data_source ads using(azure_data_source_id)
  inner join azure using(azure_id, project_id)
where
	rads.project_id=? and rads.azure_id=? and rads.azure_data_source_id=?
`

func (c *Client) GetRelAzureDataSource(ctx context.Context, projectID, azureID, azureDataSourceID uint32) (*RelAzureDataSource, error) {
	data := RelAzureDataSource{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetRelAzureDataSource, projectID, azureID, azureDataSourceID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertRelAzureDataSource string = `
INSERT INTO rel_azure_data_source (
  azure_id,
  azure_data_source_id,
  project_id,
  status,
  status_detail,
  scan_at
)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  project_id=VALUES(project_id),
  status=VALUES(status),
  status_detail=VALUES(status_detail),
  scan_at=VALUES(scan_at)
`

func (c *Client) UpsertRelAzureDataSource(ctx context.Context, relAzureDataSource *azure.RelAzureDataSourceForUpsert) (*RelAzureDataSource, error) {
	// Check master table exists
	if _, err := c.GetAzureDataSource(ctx, relAzureDataSource.AzureDataSourceId); err != nil {
		c.logger.Errorf(ctx, "Not exists azure_data_source or DB error: azure_data_source_id=%d", relAzureDataSource.AzureDataSourceId)
		return nil, err
	}
	if _, err := c.GetAzure(ctx, relAzureDataSource.ProjectId, relAzureDataSource.AzureId); err != nil {
		c.logger.Errorf(ctx, "Not exists azure or DB error: azure_id=%d", relAzureDataSource.AzureId)
		return nil, err
	}

	// Upsert
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertRelAzureDataSource,
		relAzureDataSource.AzureId,
		relAzureDataSource.AzureDataSourceId,
		relAzureDataSource.ProjectId,
		relAzureDataSource.Status.String(),
		convertZeroValueToNull(relAzureDataSource.StatusDetail),
		time.Unix(relAzureDataSource.ScanAt, 0),
	).Error; err != nil {
		return nil, err
	}
	return c.GetRelAzureDataSource(ctx, relAzureDataSource.ProjectId, relAzureDataSource.AzureId, relAzureDataSource.AzureDataSourceId)
}

const deleteRelAzureDataSource string = `delete from rel_azure_data_source where project_id=? and azure_id=? and azure_data_source_id=?`

func (c *Client) DeleteRelAzureDataSource(ctx context.Context, projectID, azureID, azureDataSourceID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteRelAzureDataSource, projectID, azureID, azureDataSourceID).Error; err != nil {
		return err
	}
	return nil
}

const selectListRelAzureDataSourceByDataSourceID = `
select
  rads.*, ads.name, ads.max_score, ads.description, azure.subscription_id
from
  rel_azure_data_source rads
  inner join azure_data_source ads using(azure_data_source_id)
  inner join azure using(azure_id, project_id)
`

func (c *Client) ListRelAzureDataSourceByDataSourceID(ctx context.Context, azureDataSourceID uint32) (*[]RelAzureDataSource, error) {
	query := selectListRelAzureDataSourceByDataSourceID
	var params []interface{}
	if !zero.IsZeroVal(azureDataSourceID) {
		query += " where azure_data_source_id = ?"
		params = append(params, azureDataSourceID)
	}
	data := []RelAzureDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

type AzureScanError struct {
	AzureID           uint32
	AzureDataSourceID uint32
	ProjectID         uint32
	DataSource        string
	StatusDetail      string
}

const selectListAzureScanError = `
select
  rads.azure_id, rads.azure_data_source_id, ads.name as data_source, rads.project_id, rads.status_detail
from
  rel_azure_data_source rads 
  inner join azure_data_source ads using(azure_data_source_id) 
where
  rads.status = 'ERROR'
  and rads.error_notified_at is null
`

func (c *Client) ListAzureScanErrorForNotify(ctx context.Context) ([]*AzureScanError, error) {
	data := []*AzureScanError{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectListAzureScanError).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const updateAzureErrorNotifiedAt = `update rel_azure_data_source set error_notified_at = ? where azure_id = ? and azure_data_source_id = ? and project_id = ?`

func (c *Client) UpdateAzureErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, azureID, azureDataSourceID, projectID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(updateAzureErrorNotifiedAt, errNotifiedAt, azureID, azureDataSourceID, projectID).Error; err != nil {
		return err
	}
	return nil
}
