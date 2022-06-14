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
func (m *MockCodeRepository) ListGitleaks(ctx context.Context, projectID, codeDataSourceID, gitleaksID uint32) (*[]model.CodeGitleaks, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeGitleaks), args.Error(1)
}
func (m *MockCodeRepository) UpsertGitleaks(ctx context.Context, data *code.GitleaksForUpsert) (*model.CodeGitleaks, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaks), args.Error(1)
}
func (m *MockCodeRepository) DeleteGitleaks(ctx context.Context, projectID uint32, gitleaksID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockCodeRepository) GetGitleaks(ctx context.Context, projectID, gitleaksID uint32) (*model.CodeGitleaks, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeGitleaks), args.Error(1)
}
func (m *MockCodeRepository) ListEnterpriseOrg(ctx context.Context, projectID, gitleaksID uint32) (*[]model.CodeEnterpriseOrg, error) {
	args := m.Called()
	return args.Get(0).(*[]model.CodeEnterpriseOrg), args.Error(1)
}
func (m *MockCodeRepository) UpsertEnterpriseOrg(ctx context.Context, data *code.EnterpriseOrgForUpsert) (*model.CodeEnterpriseOrg, error) {
	args := m.Called()
	return args.Get(0).(*model.CodeEnterpriseOrg), args.Error(1)
}
func (m *MockCodeRepository) DeleteEnterpriseOrg(ctx context.Context, projectID, gitleaksID uint32, login string) error {
	args := m.Called()
	return args.Error(0)
}
