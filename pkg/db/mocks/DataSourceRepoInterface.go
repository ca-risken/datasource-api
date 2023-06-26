// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// DataSourceRepoInterface is an autogenerated mock type for the DataSourceRepoInterface type
type DataSourceRepoInterface struct {
	mock.Mock
}

// CleanWithNoProject provides a mock function with given fields: ctx
func (_m *DataSourceRepoInterface) CleanWithNoProject(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDataSourceRepoInterface creates a new instance of DataSourceRepoInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataSourceRepoInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataSourceRepoInterface {
	mock := &DataSourceRepoInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
