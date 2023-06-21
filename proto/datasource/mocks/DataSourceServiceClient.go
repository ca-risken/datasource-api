// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	datasource "github.com/ca-risken/datasource-api/proto/datasource"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// DataSourceServiceClient is an autogenerated mock type for the DataSourceServiceClient type
type DataSourceServiceClient struct {
	mock.Mock
}

// AnalyzeAttackFlow provides a mock function with given fields: ctx, in, opts
func (_m *DataSourceServiceClient) AnalyzeAttackFlow(ctx context.Context, in *datasource.AnalyzeAttackFlowRequest, opts ...grpc.CallOption) (*datasource.AnalyzeAttackFlowResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *datasource.AnalyzeAttackFlowResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *datasource.AnalyzeAttackFlowRequest, ...grpc.CallOption) (*datasource.AnalyzeAttackFlowResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *datasource.AnalyzeAttackFlowRequest, ...grpc.CallOption) *datasource.AnalyzeAttackFlowResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*datasource.AnalyzeAttackFlowResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *datasource.AnalyzeAttackFlowRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CleanDataSource provides a mock function with given fields: ctx, in, opts
func (_m *DataSourceServiceClient) CleanDataSource(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDataSourceServiceClient creates a new instance of DataSourceServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataSourceServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataSourceServiceClient {
	mock := &DataSourceServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
