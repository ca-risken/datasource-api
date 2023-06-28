// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"

	osint "github.com/ca-risken/datasource-api/proto/osint"
)

// OsintServiceClient is an autogenerated mock type for the OsintServiceClient type
type OsintServiceClient struct {
	mock.Mock
}

// DeleteOsint provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) DeleteOsint(ctx context.Context, in *osint.DeleteOsintRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteOsintRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) DeleteOsintDataSource(ctx context.Context, in *osint.DeleteOsintDataSourceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDataSourceRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDataSourceRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOsintDetectWord provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) DeleteOsintDetectWord(ctx context.Context, in *osint.DeleteOsintDetectWordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDetectWordRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDetectWordRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteOsintDetectWordRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRelOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) DeleteRelOsintDataSource(ctx context.Context, in *osint.DeleteRelOsintDataSourceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteRelOsintDataSourceRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteRelOsintDataSourceRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteRelOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOsint provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) GetOsint(ctx context.Context, in *osint.GetOsintRequest, opts ...grpc.CallOption) (*osint.GetOsintResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.GetOsintResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintRequest, ...grpc.CallOption) (*osint.GetOsintResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintRequest, ...grpc.CallOption) *osint.GetOsintResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetOsintResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetOsintRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) GetOsintDataSource(ctx context.Context, in *osint.GetOsintDataSourceRequest, opts ...grpc.CallOption) (*osint.GetOsintDataSourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.GetOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDataSourceRequest, ...grpc.CallOption) (*osint.GetOsintDataSourceResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDataSourceRequest, ...grpc.CallOption) *osint.GetOsintDataSourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOsintDetectWord provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) GetOsintDetectWord(ctx context.Context, in *osint.GetOsintDetectWordRequest, opts ...grpc.CallOption) (*osint.GetOsintDetectWordResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.GetOsintDetectWordResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDetectWordRequest, ...grpc.CallOption) (*osint.GetOsintDetectWordResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDetectWordRequest, ...grpc.CallOption) *osint.GetOsintDetectWordResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetOsintDetectWordResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetOsintDetectWordRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRelOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) GetRelOsintDataSource(ctx context.Context, in *osint.GetRelOsintDataSourceRequest, opts ...grpc.CallOption) (*osint.GetRelOsintDataSourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.GetRelOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetRelOsintDataSourceRequest, ...grpc.CallOption) (*osint.GetRelOsintDataSourceResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetRelOsintDataSourceRequest, ...grpc.CallOption) *osint.GetRelOsintDataSourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetRelOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetRelOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScan provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) InvokeScan(ctx context.Context, in *osint.InvokeScanRequest, opts ...grpc.CallOption) (*osint.InvokeScanResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.InvokeScanResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanRequest, ...grpc.CallOption) (*osint.InvokeScanResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanRequest, ...grpc.CallOption) *osint.InvokeScanResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.InvokeScanResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.InvokeScanRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanAll provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) InvokeScanAll(ctx context.Context, in *osint.InvokeScanAllRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanAllRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanAllRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.InvokeScanAllRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOsint provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) ListOsint(ctx context.Context, in *osint.ListOsintRequest, opts ...grpc.CallOption) (*osint.ListOsintResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.ListOsintResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintRequest, ...grpc.CallOption) (*osint.ListOsintResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintRequest, ...grpc.CallOption) *osint.ListOsintResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListOsintResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListOsintRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) ListOsintDataSource(ctx context.Context, in *osint.ListOsintDataSourceRequest, opts ...grpc.CallOption) (*osint.ListOsintDataSourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.ListOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDataSourceRequest, ...grpc.CallOption) (*osint.ListOsintDataSourceResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDataSourceRequest, ...grpc.CallOption) *osint.ListOsintDataSourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOsintDetectWord provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) ListOsintDetectWord(ctx context.Context, in *osint.ListOsintDetectWordRequest, opts ...grpc.CallOption) (*osint.ListOsintDetectWordResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.ListOsintDetectWordResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDetectWordRequest, ...grpc.CallOption) (*osint.ListOsintDetectWordResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDetectWordRequest, ...grpc.CallOption) *osint.ListOsintDetectWordResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListOsintDetectWordResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListOsintDetectWordRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRelOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) ListRelOsintDataSource(ctx context.Context, in *osint.ListRelOsintDataSourceRequest, opts ...grpc.CallOption) (*osint.ListRelOsintDataSourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.ListRelOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListRelOsintDataSourceRequest, ...grpc.CallOption) (*osint.ListRelOsintDataSourceResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListRelOsintDataSourceRequest, ...grpc.CallOption) *osint.ListRelOsintDataSourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListRelOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListRelOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutOsint provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) PutOsint(ctx context.Context, in *osint.PutOsintRequest, opts ...grpc.CallOption) (*osint.PutOsintResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.PutOsintResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintRequest, ...grpc.CallOption) (*osint.PutOsintResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintRequest, ...grpc.CallOption) *osint.PutOsintResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutOsintResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutOsintRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) PutOsintDataSource(ctx context.Context, in *osint.PutOsintDataSourceRequest, opts ...grpc.CallOption) (*osint.PutOsintDataSourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.PutOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDataSourceRequest, ...grpc.CallOption) (*osint.PutOsintDataSourceResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDataSourceRequest, ...grpc.CallOption) *osint.PutOsintDataSourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutOsintDetectWord provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) PutOsintDetectWord(ctx context.Context, in *osint.PutOsintDetectWordRequest, opts ...grpc.CallOption) (*osint.PutOsintDetectWordResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.PutOsintDetectWordResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDetectWordRequest, ...grpc.CallOption) (*osint.PutOsintDetectWordResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDetectWordRequest, ...grpc.CallOption) *osint.PutOsintDetectWordResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutOsintDetectWordResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutOsintDetectWordRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRelOsintDataSource provides a mock function with given fields: ctx, in, opts
func (_m *OsintServiceClient) PutRelOsintDataSource(ctx context.Context, in *osint.PutRelOsintDataSourceRequest, opts ...grpc.CallOption) (*osint.PutRelOsintDataSourceResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *osint.PutRelOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutRelOsintDataSourceRequest, ...grpc.CallOption) (*osint.PutRelOsintDataSourceResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutRelOsintDataSourceRequest, ...grpc.CallOption) *osint.PutRelOsintDataSourceResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutRelOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutRelOsintDataSourceRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOsintServiceClient creates a new instance of OsintServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOsintServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *OsintServiceClient {
	mock := &OsintServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
