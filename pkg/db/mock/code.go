package mock

import (
	"context"

	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/stretchr/testify/mock"
)

type MockCodeRepository struct {
	mock.Mock
}

var _ db.CodeRepoInterface = (*MockCodeRepository)(nil) // verify interface compliance

func (m *MockCodeRepository) ListCodeDataSource(ctx context.Context, codeDataSourceID uint32, name string) (*[]model.CodeDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeDataSource), args.Error(1)
}
func (m *MockCodeRepository) ListGitHubSetting(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGitHubSetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeGitHubSetting), args.Error(1)
}
func (m *MockCodeRepository) UpsertGitHubSetting(ctx context.Context, data *code.GitHubSettingForUpsert) (*model.CodeGitHubSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitHubSetting), args.Error(1)
}
func (m *MockCodeRepository) DeleteGitHubSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockCodeRepository) GetGitHubSetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeGitHubSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitHubSetting), args.Error(1)
}
func (m *MockCodeRepository) ListGitleaksSetting(ctx context.Context, projectID uint32) (*[]model.CodeGitleaksSetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeGitleaksSetting), args.Error(1)
}
func (m *MockCodeRepository) GetGitleaksSetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeGitleaksSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaksSetting), args.Error(1)
}
func (m *MockCodeRepository) UpsertGitleaksSetting(ctx context.Context, data *code.GitleaksSettingForUpsert) (*model.CodeGitleaksSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaksSetting), args.Error(1)
}
func (m *MockCodeRepository) DeleteGitleaksSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockCodeRepository) GetGitleaksCache(ctx context.Context, projectID, githubSettingID uint32, repositoryFullName string, immediately bool) (*model.CodeGitleaksCache, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaksCache), args.Error(1)
}
func (m *MockCodeRepository) UpsertGitleaksCache(ctx context.Context, projectID uint32, data *code.GitleaksCacheForUpsert) (*model.CodeGitleaksCache, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaksCache), args.Error(1)
}
func (m *MockCodeRepository) DeleteGitleaksCache(ctx context.Context, githubSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockCodeRepository) ListDependencySetting(ctx context.Context, projectID uint32) (*[]model.CodeDependencySetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeDependencySetting), args.Error(1)
}
func (m *MockCodeRepository) GetDependencySetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeDependencySetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeDependencySetting), args.Error(1)
}
func (m *MockCodeRepository) UpsertDependencySetting(ctx context.Context, data *code.DependencySettingForUpsert) (*model.CodeDependencySetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeDependencySetting), args.Error(1)
}
func (m *MockCodeRepository) DeleteDependencySetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockCodeRepository) ListGitHubEnterpriseOrg(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGitHubEnterpriseOrg, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeGitHubEnterpriseOrg), args.Error(1)
}
func (m *MockCodeRepository) UpsertGitHubEnterpriseOrg(ctx context.Context, data *code.GitHubEnterpriseOrgForUpsert) (*model.CodeGitHubEnterpriseOrg, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitHubEnterpriseOrg), args.Error(1)
}
func (m *MockCodeRepository) DeleteGitHubEnterpriseOrg(ctx context.Context, projectID, githubSettingID uint32, organization string) error {
	args := m.Called()
	return args.Error(0)
}
