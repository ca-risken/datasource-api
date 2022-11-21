// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	code "github.com/ca-risken/datasource-api/proto/code"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// CodeServiceClient is an autogenerated mock type for the CodeServiceClient type
type CodeServiceClient struct {
	mock.Mock
}

// DeleteDependencySetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) DeleteDependencySetting(ctx context.Context, in *code.DeleteDependencySettingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.DeleteDependencySettingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.DeleteDependencySettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGitHubSetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) DeleteGitHubSetting(ctx context.Context, in *code.DeleteGitHubSettingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.DeleteGitHubSettingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.DeleteGitHubSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGitleaksSetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) DeleteGitleaksSetting(ctx context.Context, in *code.DeleteGitleaksSettingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.DeleteGitleaksSettingRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.DeleteGitleaksSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGitHubSetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) GetGitHubSetting(ctx context.Context, in *code.GetGitHubSettingRequest, opts ...grpc.CallOption) (*code.GetGitHubSettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.GetGitHubSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.GetGitHubSettingRequest, ...grpc.CallOption) *code.GetGitHubSettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.GetGitHubSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.GetGitHubSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGitleaksCache provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) GetGitleaksCache(ctx context.Context, in *code.GetGitleaksCacheRequest, opts ...grpc.CallOption) (*code.GetGitleaksCacheResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.GetGitleaksCacheResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.GetGitleaksCacheRequest, ...grpc.CallOption) *code.GetGitleaksCacheResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.GetGitleaksCacheResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.GetGitleaksCacheRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanAll provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) InvokeScanAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanDependency provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) InvokeScanDependency(ctx context.Context, in *code.InvokeScanDependencyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.InvokeScanDependencyRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.InvokeScanDependencyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanGitleaks provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) InvokeScanGitleaks(ctx context.Context, in *code.InvokeScanGitleaksRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.InvokeScanGitleaksRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.InvokeScanGitleaksRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDataSource provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) ListDataSource(ctx context.Context, in *code.ListDataSourceRequest, opts ...grpc.CallOption) (*code.ListDataSourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.ListDataSourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.ListDataSourceRequest, ...grpc.CallOption) *code.ListDataSourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.ListDataSourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.ListDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGitHubSetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) ListGitHubSetting(ctx context.Context, in *code.ListGitHubSettingRequest, opts ...grpc.CallOption) (*code.ListGitHubSettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.ListGitHubSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.ListGitHubSettingRequest, ...grpc.CallOption) *code.ListGitHubSettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.ListGitHubSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.ListGitHubSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutDependencySetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) PutDependencySetting(ctx context.Context, in *code.PutDependencySettingRequest, opts ...grpc.CallOption) (*code.PutDependencySettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.PutDependencySettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutDependencySettingRequest, ...grpc.CallOption) *code.PutDependencySettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutDependencySettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutDependencySettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutGitHubSetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) PutGitHubSetting(ctx context.Context, in *code.PutGitHubSettingRequest, opts ...grpc.CallOption) (*code.PutGitHubSettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.PutGitHubSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutGitHubSettingRequest, ...grpc.CallOption) *code.PutGitHubSettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutGitHubSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutGitHubSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutGitleaksCache provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) PutGitleaksCache(ctx context.Context, in *code.PutGitleaksCacheRequest, opts ...grpc.CallOption) (*code.PutGitleaksCacheResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.PutGitleaksCacheResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutGitleaksCacheRequest, ...grpc.CallOption) *code.PutGitleaksCacheResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutGitleaksCacheResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutGitleaksCacheRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutGitleaksSetting provides a mock function with given fields: ctx, in, opts
func (_m *CodeServiceClient) PutGitleaksSetting(ctx context.Context, in *code.PutGitleaksSettingRequest, opts ...grpc.CallOption) (*code.PutGitleaksSettingResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *code.PutGitleaksSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutGitleaksSettingRequest, ...grpc.CallOption) *code.PutGitleaksSettingResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutGitleaksSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutGitleaksSettingRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCodeServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewCodeServiceClient creates a new instance of CodeServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCodeServiceClient(t mockConstructorTestingTNewCodeServiceClient) *CodeServiceClient {
	mock := &CodeServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
