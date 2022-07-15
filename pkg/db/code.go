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

	// code_github_setting
	ListGithubSetting(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGithubSetting, error)
	GetGithubSetting(ctx context.Context, projectID, GithubSettingID uint32) (*model.CodeGithubSetting, error)
	UpsertGithubSetting(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGithubSetting, error)
	DeleteGithubSetting(ctx context.Context, projectID uint32, GithubSettingID uint32) error

	// code_gitleaks_setting
	ListGitleaksSetting(ctx context.Context, projectID uint32) (*[]model.CodeGitleaksSetting, error)
	GetGitleaksSetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeGitleaksSetting, error)
	UpsertGitleaksSetting(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaksSetting, error)
	DeleteGitleaksSetting(ctx context.Context, projectID uint32, GithubSettingID uint32) error

	// code_github_enterprise_org
	ListGithubEnterpriseOrg(ctx context.Context, projectID, GithubSettingID uint32) (*[]model.CodeGithubEnterpriseOrg, error)
	UpsertGithubEnterpriseOrg(ctx context.Context, data *code.EnterpriseOrgForUpsert) (*model.CodeGithubEnterpriseOrg, error)
	DeleteGithubEnterpriseOrg(ctx context.Context, projectID, githubSettingID uint32, organization string) error
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

func (c *Client) ListGithubSetting(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGithubSetting, error) {
	query := `select * from code_github_setting where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id = ?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(githubSettingID) {
		query += " and code_github_setting_id = ?"
		params = append(params, githubSettingID)
	}
	data := []model.CodeGithubSetting{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetGithubSetting(ctx context.Context, projectID uint32, githubSettingID uint32) (*model.CodeGithubSetting, error) {
	var data model.CodeGithubSetting
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertGithubSetting(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGithubSetting, error) {
	if data.PersonalAccessToken != "" {
		return c.UpsertGithubSettingWithToken(ctx, data)
	}
	return c.UpsertGithubSettingWithoutToken(ctx, data)
}

const upsertGithubWithToken = `
INSERT INTO code_github_setting (
  code_github_setting_id,
  name,
  project_id,
  type,
  base_url,
  target_resource,
  github_user,
  personal_access_token
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	code_github_setting_id=VALUES(code_github_setting_id),
	name=VALUES(name),
	project_id=VALUES(project_id),
	type=VALUES(type),
	base_url=VALUES(base_url),
	target_resource=VALUES(target_resource),
	github_user=VALUES(github_user),
	personal_access_token=VALUES(personal_access_token)
`

func (c *Client) UpsertGithubSettingWithToken(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGithubSetting, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(upsertGithubWithToken,
		data.GitleaksId,
		convertZeroValueToNull(data.Name),
		data.ProjectId,
		data.Type.String(),
		data.BaseUrl,
		data.TargetResource,
		convertZeroValueToNull(data.GithubUser),
		convertZeroValueToNull(data.PersonalAccessToken)).Error; err != nil {
		return nil, err
	}
	return c.GetGithubSettingByUniqueIndex(ctx, data.ProjectId, data.Name)
}

const upsertGithubSettingWithoutToken = `
INSERT INTO code_github_setting (
  code_github_setting_id,
  name,
  project_id,
  type,
  base_url,
  target_resource,
  github_user
)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	code_github_setting_id=VALUES(code_github_setting_id),
	name=VALUES(name),
	project_id=VALUES(project_id),
	type=VALUES(type),
	base_url=VALUES(base_url),
	target_resource=VALUES(target_resource),
	github_user=VALUES(github_user)
`

func (c *Client) UpsertGithubSettingWithoutToken(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGithubSetting, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(upsertGithubSettingWithoutToken,
		data.GitleaksId,
		convertZeroValueToNull(data.Name),
		data.ProjectId,
		data.Type.String(),
		data.BaseUrl,
		data.TargetResource,
		convertZeroValueToNull(data.GithubUser)).Error; err != nil {
		return nil, err
	}
	return c.GetGithubSettingByUniqueIndex(ctx, data.ProjectId, data.Name)
}

func (c *Client) DeleteGithubSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).Delete(&model.CodeGithubSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListGitleaksSetting(ctx context.Context, projectID uint32) (*[]model.CodeGitleaksSetting, error) {
	query := `select * from code_gitleaks_setting where and project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	data := []model.CodeGitleaksSetting{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetGitleaksSetting(ctx context.Context, projectID uint32, githubSettingID uint32) (*model.CodeGitleaksSetting, error) {
	var data model.CodeGitleaksSetting
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const upsertGitleaksWithToken = `
INSERT INTO code_gitleaks_setting (
  code_github_setting_id,
  code_data_source_id,
  project_id,
  repository_pattern,
  scan_public,
  scan_internal,
  scan_private,
  status,
  status_detail,
  scan_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	code_data_source_id=VALUES(code_data_source_id),
	project_id=VALUES(project_id),
	repository_pattern=VALUES(repository_pattern),
	scan_public=VALUES(scan_public),
	scan_internal=VALUES(scan_internal),
	scan_private=VALUES(scan_private),
	status=VALUES(status),
	status_detail=VALUES(status_detail),
	scan_at=VALUES(scan_at)
`

func (c *Client) UpsertGitleaksSetting(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaksSetting, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(upsertGitleaksWithToken,
		data.GitleaksId,
		data.CodeDataSourceId,
		data.ProjectId,
		convertZeroValueToNull(data.RepositoryPattern),
		fmt.Sprintf("%t", data.ScanPublic),
		fmt.Sprintf("%t", data.ScanInternal),
		fmt.Sprintf("%t", data.ScanPrivate),
		data.Status.String(),
		convertZeroValueToNull(data.StatusDetail),
		time.Unix(data.ScanAt, 0)).Error; err != nil {
		return nil, err
	}
	return c.GetGitleaksSetting(ctx, data.ProjectId, data.GitleaksId)

}

func (c *Client) DeleteGitleaksSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).Delete(&model.CodeGitleaksSetting{}).Error; err != nil {
		return err
	}
	return nil
}

const selectGetCodeGithubSettingByUniqueIndex = `select * from code_github_setting where project_id=? and name=?`

func (c *Client) GetGithubSettingByUniqueIndex(ctx context.Context, projectID uint32, name string) (*model.CodeGithubSetting, error) {
	data := model.CodeGithubSetting{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetCodeGithubSettingByUniqueIndex, projectID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListGithubEnterpriseOrg(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGithubEnterpriseOrg, error) {
	query := `select * from code_github_enterprise_org where 1=1`
	var params []interface{}
	if !zero.IsZeroVal(projectID) {
		query += " and project_id=?"
		params = append(params, projectID)
	}
	if !zero.IsZeroVal(githubSettingID) {
		query += " and code_github_setting_id=?"
		params = append(params, githubSettingID)
	}
	data := []model.CodeGithubEnterpriseOrg{}
	if err := c.MasterDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertGithubEnterpriseOrg(ctx context.Context, data *code.EnterpriseOrgForUpsert) (*model.CodeGithubEnterpriseOrg, error) {
	var updated model.CodeGithubEnterpriseOrg
	if err := c.MasterDB.WithContext(ctx).
		Where("code_github_setting_id=? and organization=? and project_id=?", data.GitleaksId, data.Login, data.ProjectId).
		Assign(map[string]interface{}{
			"code_github_setting_id": data.GitleaksId,
			"organization":           data.Login,
			"project_id":             data.ProjectId,
		}).
		FirstOrCreate(&updated).
		Error; err != nil {
		return nil, err
	}
	return &model.CodeGithubEnterpriseOrg{
		CodeGithubSettingID: updated.CodeGithubSettingID,
		Organization:        updated.Organization,
		ProjectID:           data.ProjectId,
		UpdatedAt:           updated.UpdatedAt,
		CreatedAt:           updated.CreatedAt,
	}, nil
}

func (c *Client) DeleteGithubEnterpriseOrg(ctx context.Context, projectID, githubSettingID uint32, organization string) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ? AND organization = ?", projectID, githubSettingID, organization).Delete(&model.CodeGithubEnterpriseOrg{}).Error; err != nil {
		return err
	}
	return nil
}
