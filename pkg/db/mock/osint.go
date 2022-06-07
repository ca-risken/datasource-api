package mock

import (
	"context"

	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/stretchr/testify/mock"
)

type MockOsintRepository struct {
	mock.Mock
}

var _ db.OSINTRepoInterface = (*MockOsintRepository)(nil) // verify interface compliance

func (m *MockOsintRepository) ListOsint(context.Context, uint32) (*[]model.Osint, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Osint), args.Error(1)
}
func (m *MockOsintRepository) GetOsint(context.Context, uint32, uint32) (*model.Osint, error) {
	args := m.Called()
	return args.Get(0).(*model.Osint), args.Error(1)
}
func (m *MockOsintRepository) UpsertOsint(context.Context, *model.Osint) (*model.Osint, error) {
	args := m.Called()
	return args.Get(0).(*model.Osint), args.Error(1)
}
func (m *MockOsintRepository) DeleteOsint(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockOsintRepository) ListOsintDataSource(context.Context, uint32, string) (*[]model.OsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.OsintDataSource), args.Error(1)
}
func (m *MockOsintRepository) GetOsintDataSource(context.Context, uint32, uint32) (*model.OsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDataSource), args.Error(1)
}
func (m *MockOsintRepository) UpsertOsintDataSource(context.Context, *model.OsintDataSource) (*model.OsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDataSource), args.Error(1)
}
func (m *MockOsintRepository) DeleteOsintDataSource(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockOsintRepository) ListRelOsintDataSource(context.Context, uint32, uint32, uint32) (*[]model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelOsintDataSource), args.Error(1)
}
func (m *MockOsintRepository) GetRelOsintDataSource(context.Context, uint32, uint32) (*model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.RelOsintDataSource), args.Error(1)
}
func (m *MockOsintRepository) UpsertRelOsintDataSource(context.Context, *model.RelOsintDataSource) (*model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*model.RelOsintDataSource), args.Error(1)
}
func (m *MockOsintRepository) DeleteRelOsintDataSource(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockOsintRepository) ListOsintDetectWord(context.Context, uint32, uint32) (*[]model.OsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*[]model.OsintDetectWord), args.Error(1)
}
func (m *MockOsintRepository) GetOsintDetectWord(context.Context, uint32, uint32) (*model.OsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDetectWord), args.Error(1)
}
func (m *MockOsintRepository) UpsertOsintDetectWord(context.Context, *model.OsintDetectWord) (*model.OsintDetectWord, error) {
	args := m.Called()
	return args.Get(0).(*model.OsintDetectWord), args.Error(1)
}
func (m *MockOsintRepository) DeleteOsintDetectWord(context.Context, uint32, uint32) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockOsintRepository) ListAllRelOsintDataSource(context.Context, uint32) (*[]model.RelOsintDataSource, error) {
	args := m.Called()
	return args.Get(0).(*[]model.RelOsintDataSource), args.Error(1)
}
