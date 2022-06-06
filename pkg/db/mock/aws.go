package mock

import (
	"context"

	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/aws"
	"github.com/stretchr/testify/mock"
)

type MockAWSRepository struct {
	mock.Mock
}

func (m *MockAWSRepository) ListAWS(context.Context, uint32, uint32, string) (*[]model.AWS, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AWS), args.Error(1)
}
func (m *MockAWSRepository) GetAWSByAccountID(context.Context, uint32, string) (*model.AWS, error) {
	args := m.Called()
	return args.Get(0).(*model.AWS), args.Error(1)
}
func (m *MockAWSRepository) UpsertAWS(context.Context, *model.AWS) (*model.AWS, error) {
	args := m.Called()
	return args.Get(0).(*model.AWS), args.Error(1)
}
func (m *MockAWSRepository) DeleteAWS(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAWSRepository) ListAWSDataSource(context.Context, uint32, uint32, string) (*[]db.DataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]db.DataSource), args.Error(1)
}
func (m *MockAWSRepository) ListDataSourceByAWSDataSourceID(context.Context, uint32) (*[]db.DataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]db.DataSource), args.Error(1)
}
func (m *MockAWSRepository) ListAWSRelDataSource(ctx context.Context, projectID, awsID uint32) (*[]model.AWSRelDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.AWSRelDataSource), args.Error(1)
}
func (m *MockAWSRepository) UpsertAWSRelDataSource(context.Context, *aws.DataSourceForAttach) (*model.AWSRelDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.AWSRelDataSource), args.Error(1)
}
func (m *MockAWSRepository) GetAWSRelDataSourceByID(ctx context.Context, awsID, awsDataSourceID, projectID uint32) (*model.AWSRelDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.AWSRelDataSource), args.Error(1)
}
func (m *MockAWSRepository) DeleteAWSRelDataSource(context.Context, uint32, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockAWSRepository) GetAWSDataSourceForMessage(ctx context.Context, awsID, awsDataSourceID, projectID uint32) (*message.AWSQueueMessage, error) {
	args := m.Called()
	return args.Get(0).(*message.AWSQueueMessage), args.Error(1)
}
