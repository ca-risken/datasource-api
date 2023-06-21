// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	osint "github.com/ca-risken/datasource-api/proto/osint"
)

// OsintServiceServer is an autogenerated mock type for the OsintServiceServer type
type OsintServiceServer struct {
	mock.Mock
}

// DeleteOsint provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) DeleteOsint(_a0 context.Context, _a1 *osint.DeleteOsintRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteOsintRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) DeleteOsintDataSource(_a0 context.Context, _a1 *osint.DeleteOsintDataSourceRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDataSourceRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDataSourceRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOsintDetectWord provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) DeleteOsintDetectWord(_a0 context.Context, _a1 *osint.DeleteOsintDetectWordRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDetectWordRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteOsintDetectWordRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteOsintDetectWordRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRelOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) DeleteRelOsintDataSource(_a0 context.Context, _a1 *osint.DeleteRelOsintDataSourceRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteRelOsintDataSourceRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.DeleteRelOsintDataSourceRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.DeleteRelOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOsint provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) GetOsint(_a0 context.Context, _a1 *osint.GetOsintRequest) (*osint.GetOsintResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.GetOsintResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintRequest) (*osint.GetOsintResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintRequest) *osint.GetOsintResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetOsintResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetOsintRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) GetOsintDataSource(_a0 context.Context, _a1 *osint.GetOsintDataSourceRequest) (*osint.GetOsintDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.GetOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDataSourceRequest) (*osint.GetOsintDataSourceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDataSourceRequest) *osint.GetOsintDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOsintDetectWord provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) GetOsintDetectWord(_a0 context.Context, _a1 *osint.GetOsintDetectWordRequest) (*osint.GetOsintDetectWordResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.GetOsintDetectWordResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDetectWordRequest) (*osint.GetOsintDetectWordResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetOsintDetectWordRequest) *osint.GetOsintDetectWordResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetOsintDetectWordResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetOsintDetectWordRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRelOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) GetRelOsintDataSource(_a0 context.Context, _a1 *osint.GetRelOsintDataSourceRequest) (*osint.GetRelOsintDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.GetRelOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetRelOsintDataSourceRequest) (*osint.GetRelOsintDataSourceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.GetRelOsintDataSourceRequest) *osint.GetRelOsintDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.GetRelOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.GetRelOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScan provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) InvokeScan(_a0 context.Context, _a1 *osint.InvokeScanRequest) (*osint.InvokeScanResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.InvokeScanResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanRequest) (*osint.InvokeScanResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanRequest) *osint.InvokeScanResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.InvokeScanResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.InvokeScanRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanAll provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) InvokeScanAll(_a0 context.Context, _a1 *osint.InvokeScanAllRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanAllRequest) (*emptypb.Empty, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.InvokeScanAllRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.InvokeScanAllRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOsint provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) ListOsint(_a0 context.Context, _a1 *osint.ListOsintRequest) (*osint.ListOsintResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.ListOsintResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintRequest) (*osint.ListOsintResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintRequest) *osint.ListOsintResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListOsintResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListOsintRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) ListOsintDataSource(_a0 context.Context, _a1 *osint.ListOsintDataSourceRequest) (*osint.ListOsintDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.ListOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDataSourceRequest) (*osint.ListOsintDataSourceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDataSourceRequest) *osint.ListOsintDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOsintDetectWord provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) ListOsintDetectWord(_a0 context.Context, _a1 *osint.ListOsintDetectWordRequest) (*osint.ListOsintDetectWordResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.ListOsintDetectWordResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDetectWordRequest) (*osint.ListOsintDetectWordResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListOsintDetectWordRequest) *osint.ListOsintDetectWordResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListOsintDetectWordResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListOsintDetectWordRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRelOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) ListRelOsintDataSource(_a0 context.Context, _a1 *osint.ListRelOsintDataSourceRequest) (*osint.ListRelOsintDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.ListRelOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListRelOsintDataSourceRequest) (*osint.ListRelOsintDataSourceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.ListRelOsintDataSourceRequest) *osint.ListRelOsintDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.ListRelOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.ListRelOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutOsint provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) PutOsint(_a0 context.Context, _a1 *osint.PutOsintRequest) (*osint.PutOsintResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.PutOsintResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintRequest) (*osint.PutOsintResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintRequest) *osint.PutOsintResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutOsintResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutOsintRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) PutOsintDataSource(_a0 context.Context, _a1 *osint.PutOsintDataSourceRequest) (*osint.PutOsintDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.PutOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDataSourceRequest) (*osint.PutOsintDataSourceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDataSourceRequest) *osint.PutOsintDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutOsintDetectWord provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) PutOsintDetectWord(_a0 context.Context, _a1 *osint.PutOsintDetectWordRequest) (*osint.PutOsintDetectWordResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.PutOsintDetectWordResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDetectWordRequest) (*osint.PutOsintDetectWordResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutOsintDetectWordRequest) *osint.PutOsintDetectWordResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutOsintDetectWordResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutOsintDetectWordRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutRelOsintDataSource provides a mock function with given fields: _a0, _a1
func (_m *OsintServiceServer) PutRelOsintDataSource(_a0 context.Context, _a1 *osint.PutRelOsintDataSourceRequest) (*osint.PutRelOsintDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *osint.PutRelOsintDataSourceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutRelOsintDataSourceRequest) (*osint.PutRelOsintDataSourceResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *osint.PutRelOsintDataSourceRequest) *osint.PutRelOsintDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*osint.PutRelOsintDataSourceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *osint.PutRelOsintDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOsintServiceServer creates a new instance of OsintServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOsintServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *OsintServiceServer {
	mock := &OsintServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
