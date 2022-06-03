package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/vikyd/zero"
)

type CodeRepoInterface interface {
	// code_data_source
	ListCodeDataSource(ctx context.Context, codeDataSourceID uint32, name string) (*[]model.CodeDataSource, error)

	// code_gitleaks
	ListGitleaks(ctx context.Context, projectID, codeDataSourceID, gitleaksID uint32) (*[]model.CodeGitleaks, error)
	UpsertGitleaks(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaks, error)
	DeleteGitleaks(ctx context.Context, projectID uint32, gitleaksID uint32) error
	GetGitleaks(ctx context.Context, projectID, gitleaksID uint32) (*model.CodeGitleaks, error)

	// code_enterprise_org
	ListEnterpriseOrg(ctx context.Context, projectID, gitleaksID uint32) (*[]model.CodeEnterpriseOrg, error)
	UpsertEnterpriseOrg(ctx context.Context, data *code.EnterpriseOrgForUpsert) (*model.CodeEnterpriseOrg, error)
	DeleteEnterpriseOrg(ctx context.Context, projectID, gitleaksID uint32, login string) error
}

var _ CodeRepoInterface = (*Client)(nil) // verify interface compliance

func (c *Client) ListCodeDataSource(ctx context.Context, codeDataSourceID uint32, name string) (*[]model.CodeDataSource, error) {
	query := `select * from code_data_source where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(codeDataSourceID) {
		query += " and code_data_source_id = ?"
		params = append(params, codeDataSourceID)
	}
	if !zero.IsZeroVal(name) {
		query += " and name = ?"
		params = append(params, name)
	}
	data := []model.CodeDataSource{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListGitleaks(ctx context.Context, projectID, codeDataSourceID, gitleaksID uint32) (*[]model.CodeGitleaks, error) {
	query := `select * from code_gitleaks where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(codeDataSourceID) {
		query += " and code_data_source_id = ?"
		params = append(params, codeDataSourceID)
	}
	if !zero.IsZeroVal(gitleaksID) {
		query += " and gitleaks_id = ?"
		params = append(params, gitleaksID)
	}
	data := []model.CodeGitleaks{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertGitleaks(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaks, error) {
	if data.PersonalAccessToken != "" {
		return c.UpsertGitleaksWithToken(ctx, data)
	}
	return c.UpsertGitleaksWithoutToken(ctx, data)
}

const insertUpsertGitleaksWithToken = `
INSERT INTO code_gitleaks (
  gitleaks_id,
  code_data_source_id,
  name,
  project_id,
  type,
  base_url,
  target_resource,
  repository_pattern,
  github_user,
  personal_access_token,
  scan_public,
  scan_internal,
  scan_private,
  gitleaks_config,
  status,
  status_detail,
  scan_at,
  scan_succeeded_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	code_data_source_id=VALUES(code_data_source_id),
	name=VALUES(name),
	project_id=VALUES(project_id),
	type=VALUES(type),
	base_url=VALUES(base_url),
	target_resource=VALUES(target_resource),
	repository_pattern=VALUES(repository_pattern),
	github_user=VALUES(github_user),
	personal_access_token=VALUES(personal_access_token),
	scan_public=VALUES(scan_public),
	scan_internal=VALUES(scan_internal),
	scan_private=VALUES(scan_private),
	gitleaks_config=VALUES(gitleaks_config),
	status=VALUES(status),
	status_detail=VALUES(status_detail),
	scan_at=VALUES(scan_at),
	scan_succeeded_at=VALUES(scan_succeeded_at)
`

func (c *Client) UpsertGitleaksWithToken(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaks, error) {
	var scanSucceededAtTime time.Time
	if !zero.IsZeroVal(data.ScanSucceededAt) {
		scanSucceededAtTime = time.Unix(data.ScanSucceededAt, 0)
	}
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertGitleaksWithToken,
		data.GitleaksId,
		data.CodeDataSourceId,
		convertZeroValueToNull(data.Name),
		data.ProjectId,
		data.Type.String(),
		data.BaseUrl,
		data.TargetResource,
		convertZeroValueToNull(data.RepositoryPattern),
		convertZeroValueToNull(data.GithubUser),
		convertZeroValueToNull(data.PersonalAccessToken),
		fmt.Sprintf("%t", data.ScanPublic),
		fmt.Sprintf("%t", data.ScanInternal),
		fmt.Sprintf("%t", data.ScanPrivate),
		convertZeroValueToNull(data.GitleaksConfig),
		data.Status.String(),
		convertZeroValueToNull(data.StatusDetail),
		time.Unix(data.ScanAt, 0),
		convertZeroValueToNull(scanSucceededAtTime)).Error; err != nil {
		return nil, err
	}
	return c.GetGitleaksByUniqueIndex(ctx, data.ProjectId, data.CodeDataSourceId, data.Name)
}

const insertUpsertGitleaksWithoutToken = `
INSERT INTO code_gitleaks (
  gitleaks_id,
  code_data_source_id,
  name,
  project_id,
  type,
  base_url,
  target_resource,
  repository_pattern,
  github_user,
  scan_public,
  scan_internal,
  scan_private,
  gitleaks_config,
  status,
  status_detail,
  scan_at,
  scan_succeeded_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	code_data_source_id=VALUES(code_data_source_id),
	name=VALUES(name),
	project_id=VALUES(project_id),
	type=VALUES(type),
	base_url=VALUES(base_url),
	target_resource=VALUES(target_resource),
	repository_pattern=VALUES(repository_pattern),
	github_user=VALUES(github_user),
	scan_public=VALUES(scan_public),
	scan_internal=VALUES(scan_internal),
	scan_private=VALUES(scan_private),
	gitleaks_config=VALUES(gitleaks_config),
	status=VALUES(status),
	status_detail=VALUES(status_detail),
	scan_at=VALUES(scan_at),
	scan_succeeded_at=VALUES(scan_succeeded_at)
`

func (c *Client) UpsertGitleaksWithoutToken(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaks, error) {
	var scanSucceededAtTime time.Time
	if !zero.IsZeroVal(data.ScanSucceededAt) {
		scanSucceededAtTime = time.Unix(data.ScanSucceededAt, 0)
	}
	if err := c.MasterDB.WithContext(ctx).Exec(insertUpsertGitleaksWithoutToken,
		data.GitleaksId,
		data.CodeDataSourceId,
		convertZeroValueToNull(data.Name),
		data.ProjectId,
		data.Type.String(),
		data.BaseUrl,
		data.TargetResource,
		convertZeroValueToNull(data.RepositoryPattern),
		convertZeroValueToNull(data.GithubUser),
		fmt.Sprintf("%t", data.ScanPublic),
		fmt.Sprintf("%t", data.ScanInternal),
		fmt.Sprintf("%t", data.ScanPrivate),
		convertZeroValueToNull(data.GitleaksConfig),
		data.Status.String(),
		convertZeroValueToNull(data.StatusDetail),
		time.Unix(data.ScanAt, 0),
		convertZeroValueToNull(scanSucceededAtTime)).Error; err != nil {
		return nil, err
	}
	return c.GetGitleaksByUniqueIndex(ctx, data.ProjectId, data.CodeDataSourceId, data.Name)
}

const deleteGitleaks = `delete from code_gitleaks where project_id=? and gitleaks_id=?`

func (c *Client) DeleteGitleaks(ctx context.Context, projectID, gitleaksID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteGitleaks, projectID, gitleaksID).Error; err != nil {
		return err
	}
	return nil
}

const selectGetCodeGitleaks = `select * from code_gitleaks where project_id=? and gitleaks_id=?`

func (c *Client) GetGitleaks(ctx context.Context, projectID, gitleaksID uint32) (*model.CodeGitleaks, error) {
	data := model.CodeGitleaks{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectGetCodeGitleaks, projectID, gitleaksID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetCodeGitleaksByUniqueIndex = `select * from code_gitleaks where project_id=? and code_data_source_id=? and name=?`

func (c *Client) GetGitleaksByUniqueIndex(ctx context.Context, projectID, codeDataSourceID uint32, name string) (*model.CodeGitleaks, error) {
	data := model.CodeGitleaks{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetCodeGitleaksByUniqueIndex, projectID, codeDataSourceID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListEnterpriseOrg(ctx context.Context, projectID, gitleaksID uint32) (*[]model.CodeEnterpriseOrg, error) {
	query := `select * from code_enterprise_org where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id=?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(gitleaksID) {
		query += " and gitleaks_id=?"
		params = append(params, gitleaksID)
	}
	data := []model.CodeEnterpriseOrg{}
	if err := c.MasterDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertEnterpriseOrg(ctx context.Context, data *code.EnterpriseOrgForUpsert) (*model.CodeEnterpriseOrg, error) {
	var updated model.CodeEnterpriseOrg
	if err := c.MasterDB.WithContext(ctx).
		Where("gitleaks_id=? and login=? and project_id=?", data.GitleaksId, data.Login, data.ProjectId).
		Assign(map[string]interface{}{
			"gitleaks_id": data.GitleaksId,
			"login":       data.Login,
			"project_id":  data.ProjectId,
		}).
		FirstOrCreate(&updated).
		Error; err != nil {
		return nil, err
	}
	return &model.CodeEnterpriseOrg{
		GitleaksID: updated.GitleaksID,
		Login:      updated.Login,
		ProjectID:  data.ProjectId,
		UpdatedAt:  updated.UpdatedAt,
		CreatedAt:  updated.CreatedAt,
	}, nil
}

const deleteEnterpriseOrg = `delete from code_enterprise_org where project_id=? and gitleaks_id=? and login=?`

func (c *Client) DeleteEnterpriseOrg(ctx context.Context, projectID, gitleaksID uint32, login string) error {
	if err := c.MasterDB.WithContext(ctx).Exec(deleteEnterpriseOrg, projectID, gitleaksID, login).Error; err != nil {
		return err
	}
	return nil
}
