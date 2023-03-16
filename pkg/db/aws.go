package db

import (
	"context"
	"time"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/aws"
	"github.com/vikyd/zero"
)

type AWSRepoInterface interface {
	ListAWS(ctx context.Context, projectID, awsID uint32, awsAccountID string) (*[]model.AWS, error)
	GetAWSByAccountID(ctx context.Context, projectID uint32, awsAccountID string) (*model.AWS, error)
	UpsertAWS(ctx context.Context, data *model.AWS) (*model.AWS, error)
	DeleteAWS(ctx context.Context, projectID, awsID uint32) error
	ListAWSDataSource(ctx context.Context, projectID, awsID uint32, ds string) (*[]DataSource, error)
	ListDataSourceByAWSDataSourceID(ctx context.Context, awsDataSourceID uint32) (*[]DataSource, error)
	ListAWSRelDataSource(ctx context.Context, projectID, awsID uint32) (*[]model.AWSRelDataSource, error)
	UpsertAWSRelDataSource(ctx context.Context, data *aws.DataSourceForAttach) (*model.AWSRelDataSource, error)
	GetAWSRelDataSourceByID(ctx context.Context, awsID, awsDataSourceID, projectID uint32) (*model.AWSRelDataSource, error)
	DeleteAWSRelDataSource(ctx context.Context, projectID, awsID, awsDataSourceID uint32) error
	GetAWSDataSourceForMessage(ctx context.Context, awsID, awsDataSourceID, projectID uint32) (*DataSource, error)
}

var _ AWSRepoInterface = (*Client)(nil) // verify interface compliance

func (c *Client) ListAWS(ctx context.Context, projectID, awsID uint32, awsAccountID string) (*[]model.AWS, error) {
	query := `
select
  *
from
  aws
where
  project_id = ?
`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(awsID) {
		query += " and aws_id = ?"
		params = append(params, awsID)
	}
	if !zero.IsZeroVal(awsAccountID) {
		query += " and aws_account_id = ?"
		params = append(params, awsAccountID)
	}

	data := []model.AWS{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetAWSByAccountID = `select * from aws where project_id = ? and aws_account_id = ?`

func (c *Client) GetAWSByAccountID(ctx context.Context, projectID uint32, awsAccountID string) (*model.AWS, error) {
	data := model.AWS{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetAWSByAccountID, projectID, awsAccountID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertAWS = `
INSERT INTO aws
  (aws_id, name, project_id, aws_account_id)
VALUES
  (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name)
`

func (c *Client) UpsertAWS(ctx context.Context, data *model.AWS) (*model.AWS, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertAWS,
		data.AWSID, data.Name, data.ProjectID, data.AWSAccountID).Error; err != nil {
		return nil, err
	}

	updated, err := c.GetAWSByAccountID(ctx, data.ProjectID, data.AWSAccountID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

const deleteAws = `delete from aws where project_id = ? and aws_id = ?`

func (c *Client) DeleteAWS(ctx context.Context, projectID, awsID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteAws, projectID, awsID).Error; err != nil {
		return err
	}
	return nil
}

type DataSource struct {
	AWSDataSourceID uint32
	DataSource      string
	MaxScore        float32
	AWSID           uint32 `gorm:"column:aws_id"`
	ProjectID       uint32
	AWSAccountID    string
	AssumeRoleArn   string
	ExternalID      string
	SpecificVersion string
	Status          string
	StatusDetail    string
	ScanAt          time.Time
}

func (c *Client) ListAWSDataSource(ctx context.Context, projectID, awsID uint32, ds string) (*[]DataSource, error) {
	var params []interface{}
	query := `
select
  ads.aws_data_source_id
  , ads.data_source
  , ads.max_score
  , ards.aws_id
  , ards.project_id
  , ards.assume_role_arn
  , ards.external_id
  , ards.specific_version
  , ards.status
  , ards.status_detail
  , ards.scan_at
from
  aws_data_source ads
  left outer join (
    select * from aws_rel_data_source where 1=1 `
	if !zero.IsZeroVal(projectID) {
		query += " and project_id = ? "
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(awsID) {
		query += " and aws_id = ?"
		params = append(params, awsID)
	}
	query += `
  ) ards using(aws_data_source_id)`
	if !zero.IsZeroVal(ds) {
		query += `
where
  ads.data_source = ?`
		params = append(params, ds)
	}
	query += `
order by
  ards.project_id
  , ards.aws_id
  , ads.aws_data_source_id
`
	data := []DataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListDataSourceByAWSDataSourceID(ctx context.Context, dataSourceID uint32) (*[]DataSource, error) {
	var params []interface{}
	query := `select ads.aws_data_source_id
	, ads.data_source
	, ads.max_score
	, ards.aws_id
	, ards.project_id
	, ards.assume_role_arn
	, ards.external_id
	, ards.specific_version
	, ards.status
	, ards.status_detail
	, ards.scan_at 
	from aws_data_source ads
	inner join aws_rel_data_source ards using(aws_data_source_id)`
	if !zero.IsZeroVal(dataSourceID) {
		query += " where aws_data_source_id = ?"
		params = append(params, dataSourceID)
	}
	data := []DataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectListAWSRelDataSource = "select * from aws_rel_data_source where project_id=? and aws_id=?"

func (c *Client) ListAWSRelDataSource(ctx context.Context, projectID, awsID uint32) (*[]model.AWSRelDataSource, error) {
	data := []model.AWSRelDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectListAWSRelDataSource, projectID, awsID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertUpsertAWSRelDataSource = `
INSERT INTO aws_rel_data_source
  (aws_id, aws_data_source_id, project_id, assume_role_arn, external_id, specific_version, status, status_detail, scan_at)
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  project_id=VALUES(project_id),
  assume_role_arn=VALUES(assume_role_arn),
  external_id=VALUES(external_id),
  specific_version=VALUES(specific_version),
  status=VALUES(status),
  status_detail=VALUES(status_detail),
  scan_at=VALUES(scan_at)
`

func (c *Client) UpsertAWSRelDataSource(ctx context.Context, data *aws.DataSourceForAttach) (*model.AWSRelDataSource, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertAWSRelDataSource,
		data.AwsId, data.AwsDataSourceId, data.ProjectId,
		data.AssumeRoleArn, data.ExternalId, data.SpecificVersion,
		data.Status.String(), data.StatusDetail, time.Unix(data.ScanAt, 0),
	).Error; err != nil {
		return nil, err
	}
	return c.GetAWSRelDataSourceByID(ctx, data.AwsId, data.AwsDataSourceId, data.ProjectId)
}

const selectGetAWSRelDataSourceByID = `select * from aws_rel_data_source where aws_id = ? and aws_data_source_id = ? and project_id = ?`

func (c *Client) GetAWSRelDataSourceByID(ctx context.Context, awsID, awsDataSourceID, projectID uint32) (*model.AWSRelDataSource, error) {
	data := model.AWSRelDataSource{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetAWSRelDataSourceByID, awsID, awsDataSourceID, projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteAWSRelDataSource = `delete from aws_rel_data_source where project_id = ? and aws_id = ? and aws_data_source_id = ?`

func (c *Client) DeleteAWSRelDataSource(ctx context.Context, projectID, awsID, awsDataSourceID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteAWSRelDataSource, projectID, awsID, awsDataSourceID).Error; err != nil {
		return err
	}
	return nil
}

const selectAWSDataSourceForMessage = `
select 
    a.aws_id                as aws_id
  , ards.aws_data_source_id as aws_data_source_id
  , ads.data_source         as data_source
  , ards.project_id         as project_id
  , a.aws_account_id        as aws_account_id
  , ards.assume_role_arn    as assume_role_arn
  , ards.external_id        as external_id
  , ards.specific_version   as specific_version
from
  aws_rel_data_source ards
  inner join aws a using(aws_id)
  inner join aws_data_source ads using(aws_data_source_id)
where
  ards.aws_id = ?
  and ards.aws_data_source_id = ?
  and ards.project_id = ? 
`

func (c *Client) GetAWSDataSourceForMessage(ctx context.Context, awsID, awsDataSourceID, projectID uint32) (*DataSource, error) {
	data := DataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectAWSDataSourceForMessage, awsID, awsDataSourceID, projectID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
