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
func (m *MockCodeRepository) ListGithubSetting(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGithubSetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeGithubSetting), args.Error(1)
}
func (m *MockCodeRepository) UpsertGithubSetting(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGithubSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGithubSetting), args.Error(1)
}
func (m *MockCodeRepository) DeleteGithubSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockCodeRepository) GetGithubSetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeGithubSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGithubSetting), args.Error(1)
}
func (m *MockCodeRepository) ListGitleaksSetting(ctx context.Context, projectID uint32) (*[]model.CodeGitleaksSetting, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeGitleaksSetting), args.Error(1)
}
func (m *MockCodeRepository) GetGitleaksSetting(ctx context.Context, projectID, githubSettingID uint32) (*model.CodeGitleaksSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaksSetting), args.Error(1)
}
func (m *MockCodeRepository) UpsertGitleaksSetting(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaksSetting, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaksSetting), args.Error(1)
}
func (m *MockCodeRepository) DeleteGitleaksSetting(ctx context.Context, projectID uint32, githubSettingID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockCodeRepository) ListGithubEnterpriseOrg(ctx context.Context, projectID, githubSettingID uint32) (*[]model.CodeGithubEnterpriseOrg, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeGithubEnterpriseOrg), args.Error(1)
}
func (m *MockCodeRepository) UpsertGithubEnterpriseOrg(ctx context.Context, data *code.EnterpriseOrgForUpsert) (*model.CodeGithubEnterpriseOrg, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGithubEnterpriseOrg), args.Error(1)
}
func (m *MockCodeRepository) DeleteGithubEnterpriseOrg(ctx context.Context, projectID, githubSettingID uint32, organization string) error {
	args := m.Called()
	return args.Error(0)
}
