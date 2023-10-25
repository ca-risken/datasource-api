// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	assetpb "cloud.google.com/go/asset/apiv1/assetpb"

	mock "github.com/stretchr/testify/mock"
)

// GcpServiceClient is an autogenerated mock type for the GcpServiceClient type
type GcpServiceClient struct {
	mock.Mock
}

// GetAsset provides a mock function with given fields: ctx, gcpProjectID, resourceName
func (_m *GcpServiceClient) GetAsset(ctx context.Context, gcpProjectID string, resourceName string) (*assetpb.ResourceSearchResult, error) {
	ret := _m.Called(ctx, gcpProjectID, resourceName)

	var r0 *assetpb.ResourceSearchResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*assetpb.ResourceSearchResult, error)); ok {
		return rf(ctx, gcpProjectID, resourceName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *assetpb.ResourceSearchResult); ok {
		r0 = rf(ctx, gcpProjectID, resourceName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assetpb.ResourceSearchResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, gcpProjectID, resourceName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyCode provides a mock function with given fields: ctx, gcpProjectID, verificationCode
func (_m *GcpServiceClient) VerifyCode(ctx context.Context, gcpProjectID string, verificationCode string) (bool, error) {
	ret := _m.Called(ctx, gcpProjectID, verificationCode)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (bool, error)); ok {
		return rf(ctx, gcpProjectID, verificationCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, gcpProjectID, verificationCode)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, gcpProjectID, verificationCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGcpServiceClient creates a new instance of GcpServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGcpServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *GcpServiceClient {
	mock := &GcpServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
