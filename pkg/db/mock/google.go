package mock

import (
	"context"

	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/google"
	"github.com/stretchr/testify/mock"
)

type MockGoogleRepository struct {
	mock.Mock
}

var _ db.GoogleRepoInterface = (*MockGoogleRepository)(nil) // verify interface compliance

func (m *MockGoogleRepository) ListGoogleDataSource(ctx context.Context, googleDataSourceID uint32, name string) (*[]model.GoogleDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.GoogleDataSource), args.Error(1)
}
func (m *MockGoogleRepository) ListGCP(ctx context.Context, projectID, gcpID uint32, gcpProjectID string) (*[]model.GCP, error) {
	args := m.Called()
	return args.Get(0).(*[]model.GCP), args.Error(1)
}
func (m *MockGoogleRepository) GetGCP(ctx context.Context, projectID, gcpID uint32) (*model.GCP, error) {
	args := m.Called()
	return args.Get(0).(*model.GCP), args.Error(1)
}
func (m *MockGoogleRepository) UpsertGCP(ctx context.Context, data *google.GCPForUpsert) (*model.GCP, error) {
	args := m.Called()
	return args.Get(0).(*model.GCP), args.Error(1)
}
func (m *MockGoogleRepository) DeleteGCP(ctx context.Context, projectID uint32, gcpID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockGoogleRepository) ListGCPDataSource(ctx context.Context, projectID, gcpID uint32) (*[]db.GCPDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]db.GCPDataSource), args.Error(1)
}
func (m *MockGoogleRepository) GetGCPDataSource(ctx context.Context, projectID, gcpID, googleDataSourceID uint32) (*db.GCPDataSource, error) {
	args := m.Called()
	return args.Get(0).(*db.GCPDataSource), args.Error(1)
}
func (m *MockGoogleRepository) UpsertGCPDataSource(ctx context.Context, _ *google.GCPDataSourceForUpsert) (*db.GCPDataSource, error) {
	args := m.Called()
	return args.Get(0).(*db.GCPDataSource), args.Error(1)
}
func (m *MockGoogleRepository) DeleteGCPDataSource(ctx context.Context, projectID, gcpID, googleDataSourceID uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockGoogleRepository) ListGCPDataSourceByDataSourceID(ctx context.Context, googleDataSourceID uint32) (*[]db.GCPDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]db.GCPDataSource), args.Error(1)
}
