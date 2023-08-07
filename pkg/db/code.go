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
	ListGitHubSetting(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGitHubSetting, error)
	GetGitHubSetting(ctx context.Context, projectID, GitHubSettingID uint32) (*model.CodeGitHubSetting, error)
	UpsertGitHubSetting(ctx context.Context, data *code.GitHubSettingForUpsert) (*model.CodeGitHubSetting, error)
	DeleteGitHubSetting(ctx context.Context, projectID uint32, GitHubSettingID uint32) error

	// code_gitleaks_setting
	ListGitleaksSetting(ctx context.Context, projectID uint32) (*[]model.CodeGitleaksSetting, error)
	GetGitleaksSetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeGitleaksSetting, error)
	UpsertGitleaksSetting(ctx context.Context, data *code.GitleaksSettingForUpsert) (*model.CodeGitleaksSetting, error)
	DeleteGitleaksSetting(ctx context.Context, projectID uint32, GitHubSettingID uint32) error

	// code_gitleaks_cache
	ListGitleaksCache(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGitleaksCache, error)
	GetGitleaksCache(ctx context.Context, projectID, githubSettingID uint32, repositoryFullName string, immediately bool) (*model.CodeGitleaksCache, error)
	UpsertGitleaksCache(ctx context.Context, projectID uint32, data *code.GitleaksCacheForUpsert) (*model.CodeGitleaksCache, error)
	DeleteGitleaksCache(ctx context.Context, githubSettingID uint32) error

	// code_dependency_setting
	ListDependencySetting(ctx context.Context, projectID uint32) (*[]model.CodeDependencySetting, error)
	GetDependencySetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeDependencySetting, error)
	UpsertDependencySetting(ctx context.Context, data *code.DependencySettingForUpsert) (*model.CodeDependencySetting, error)
	DeleteDependencySetting(ctx context.Context, projectID uint32, GitHubSettingID uint32) error

	// scan error
	ListCodeGitHubScanErrorForNotify(ctx context.Context) ([]*GitHubScanError, error)
	UpdateCodeGitleaksErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, codeGithubSettingID, projectID uint32) error
	UpdateCodeDependencyErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, codeGithubSettingID, projectID uint32) error
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

func (c *Client) ListGitHubSetting(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGitHubSetting, error) {
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
	data := []model.CodeGitHubSetting{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetGitHubSetting(ctx context.Context, projectID uint32, githubSettingID uint32) (*model.CodeGitHubSetting, error) {
	var data model.CodeGitHubSetting
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertGitHubSetting(ctx context.Context, data *code.GitHubSettingForUpsert) (*model.CodeGitHubSetting, error) {
	if data.PersonalAccessToken != "" {
		return c.UpsertGitHubSettingWithToken(ctx, data)
	}
	return c.UpsertGitHubSettingWithoutToken(ctx, data)
}

const upsertGitHubWithToken = `
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

func (c *Client) UpsertGitHubSettingWithToken(ctx context.Context, data *code.GitHubSettingForUpsert) (*model.CodeGitHubSetting, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(upsertGitHubWithToken,
		data.GithubSettingId,
		convertZeroValueToNull(data.Name),
		data.ProjectId,
		data.Type.String(),
		data.BaseUrl,
		data.TargetResource,
		convertZeroValueToNull(data.GithubUser),
		convertZeroValueToNull(data.PersonalAccessToken)).Error; err != nil {
		return nil, err
	}
	return c.GetGitHubSettingByUniqueIndex(ctx, data.ProjectId, data.Name)
}

const upsertGitHubSettingWithoutToken = `
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

func (c *Client) UpsertGitHubSettingWithoutToken(ctx context.Context, data *code.GitHubSettingForUpsert) (*model.CodeGitHubSetting, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(upsertGitHubSettingWithoutToken,
		data.GithubSettingId,
		convertZeroValueToNull(data.Name),
		data.ProjectId,
		data.Type.String(),
		data.BaseUrl,
		data.TargetResource,
		convertZeroValueToNull(data.GithubUser)).Error; err != nil {
		return nil, err
	}
	return c.GetGitHubSettingByUniqueIndex(ctx, data.ProjectId, data.Name)
}

func (c *Client) DeleteGitHubSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).Delete(&model.CodeGitHubSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListGitleaksSetting(ctx context.Context, projectID uint32) (*[]model.CodeGitleaksSetting, error) {
	query := `select * from code_gitleaks_setting`
	var params []interface{}
	if projectID != 0 {
		query += " where project_id = ?"
		params = append(params, projectID)
	}
	data := []model.CodeGitleaksSetting{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetGitleaksSetting(ctx context.Context, projectID uint32, githubSettingID uint32) (*model.CodeGitleaksSetting, error) {
	var data model.CodeGitleaksSetting
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).First(&data).Error; err != nil {
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

func (c *Client) UpsertGitleaksSetting(ctx context.Context, data *code.GitleaksSettingForUpsert) (*model.CodeGitleaksSetting, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(upsertGitleaksWithToken,
		data.GithubSettingId,
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
	return c.GetGitleaksSetting(ctx, data.ProjectId, data.GithubSettingId)
}

func (c *Client) DeleteGitleaksSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).Delete(&model.CodeGitleaksSetting{}).Error; err != nil {
		return err
	}
	return nil
}

const selectListGitleaksCache = `
select
  cache.*
from 
  code_gitleaks_cache cache
  inner join code_github_setting github using(code_github_setting_id)
where 
  github.project_id = ?
  and cache.code_github_setting_id = ?
`

func (c *Client) ListGitleaksCache(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGitleaksCache, error) {
	var data []model.CodeGitleaksCache
	if err := c.SlaveDB.WithContext(ctx).Raw(selectListGitleaksCache, projectID, githubSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetGitleaksCache = `
select
  cache.*
from 
  code_gitleaks_cache cache
  inner join code_github_setting github using(code_github_setting_id)
where 
  github.project_id = ?
  and cache.code_github_setting_id = ? 
  and cache.repository_full_name = ?
`

func (c *Client) GetGitleaksCache(ctx context.Context, projectID, githubSettingID uint32, repositoryFullName string, immediately bool) (*model.CodeGitleaksCache, error) {
	db := c.SlaveDB
	if immediately {
		db = c.MasterDB
	}
	var data model.CodeGitleaksCache
	if err := db.WithContext(ctx).Raw(selectGetGitleaksCache, projectID, githubSettingID, repositoryFullName).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const upsertGitleaksCache = `
INSERT INTO code_gitleaks_cache (
  code_github_setting_id,
  repository_full_name,
  scan_at
)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	scan_at=VALUES(scan_at)
`

func (c *Client) UpsertGitleaksCache(ctx context.Context, projectID uint32, data *code.GitleaksCacheForUpsert) (*model.CodeGitleaksCache, error) {
	if err := c.MasterDB.WithContext(ctx).Exec(
		upsertGitleaksCache,
		data.GithubSettingId,
		data.RepositoryFullName,
		time.Unix(data.ScanAt, 0)).Error; err != nil {
		return nil, err
	}
	return c.GetGitleaksCache(ctx, projectID, data.GithubSettingId, data.RepositoryFullName, true)
}

func (c *Client) DeleteGitleaksCache(ctx context.Context, githubSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).
		Where("code_github_setting_id = ?", githubSettingID).
		Delete(&model.CodeGitleaksCache{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListDependencySetting(ctx context.Context, projectID uint32) (*[]model.CodeDependencySetting, error) {
	query := `select * from code_dependency_setting`
	var params []interface{}
	if projectID != 0 {
		query += " where project_id = ?"
		params = append(params, projectID)
	}
	data := []model.CodeDependencySetting{}
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetDependencySetting(ctx context.Context, projectID uint32, githubSettingID uint32) (*model.CodeDependencySetting, error) {
	var data model.CodeDependencySetting
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertDependencySetting(ctx context.Context, data *code.DependencySettingForUpsert) (*model.CodeDependencySetting, error) {
	// exclude ErrorNotifiedAt
	type codeDependencySetting struct {
		CodeGitHubSettingID uint32 `gorm:"primary_key;column:code_github_setting_id"`
		CodeDataSourceID    uint32
		ProjectID           uint32
		Status              string
		StatusDetail        string
		ScanAt              time.Time
		CreatedAt           time.Time
		UpdatedAt           time.Time
	}
	param := codeDependencySetting{
		CodeGitHubSettingID: data.GithubSettingId,
		CodeDataSourceID:    data.CodeDataSourceId,
		ProjectID:           data.ProjectId,
		Status:              data.Status.String(),
		StatusDetail:        data.StatusDetail,
		ScanAt:              time.Unix(data.ScanAt, 0),
	}
	var savedData codeDependencySetting
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", param.ProjectID, param.CodeGitHubSettingID).Assign(param).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return c.GetDependencySetting(ctx, param.ProjectID, param.CodeGitHubSettingID)
}

func (c *Client) DeleteDependencySetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND code_github_setting_id = ?", projectID, githubSettingID).Delete(&model.CodeDependencySetting{}).Error; err != nil {
		return err
	}
	return nil
}

const selectGetCodeGitHubSettingByUniqueIndex = `select * from code_github_setting where project_id=? and name=?`

func (c *Client) GetGitHubSettingByUniqueIndex(ctx context.Context, projectID uint32, name string) (*model.CodeGitHubSetting, error) {
	data := model.CodeGitHubSetting{}
	if err := c.MasterDB.WithContext(ctx).Raw(selectGetCodeGitHubSettingByUniqueIndex, projectID, name).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

type GitHubScanError struct {
	CodeGithubSettingID uint32
	ProjectID           uint32
	DataSource          string
	StatusDetail        string
}

const selectListCodeGitHubScanError = `
select
  gitleaks.code_github_setting_id, cds.name as data_source, gitleaks.project_id, gitleaks.status_detail
from
  code_gitleaks_setting gitleaks
  inner join code_data_source cds using(code_data_source_id) 
where
  gitleaks.status = 'ERROR'
  and gitleaks.error_notified_at is null 
union 
select 
  dependency.code_github_setting_id, cds.name as data_source, dependency.project_id, dependency.status_detail
from
  code_dependency_setting dependency
  inner join code_data_source cds using(code_data_source_id) 
where
  dependency.status = 'ERROR'
  and dependency.error_notified_at is null 
`

func (c *Client) ListCodeGitHubScanErrorForNotify(ctx context.Context) ([]*GitHubScanError, error) {
	data := []*GitHubScanError{}
	if err := c.SlaveDB.WithContext(ctx).Raw(selectListCodeGitHubScanError).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const updateCodeGitleaksErrorNotifiedAt = `update code_gitleaks_setting set error_notified_at = ? where code_github_setting_id = ? and project_id = ?`

func (c *Client) UpdateCodeGitleaksErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, codeGithubSettingID, projectID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(updateCodeGitleaksErrorNotifiedAt, errNotifiedAt, codeGithubSettingID, projectID).Error; err != nil {
		return err
	}
	return nil
}

const updateCodeDependencyErrorNotifiedAt = `update code_dependency_setting set error_notified_at = ? where code_github_setting_id = ? and project_id = ?`

func (c *Client) UpdateCodeDependencyErrorNotifiedAt(ctx context.Context, errNotifiedAt interface{}, codeGithubSettingID, projectID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Exec(updateCodeDependencyErrorNotifiedAt, errNotifiedAt, codeGithubSettingID, projectID).Error; err != nil {
		return err
	}
	return nil
}
