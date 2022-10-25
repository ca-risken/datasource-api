// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	context "context"

	code "github.com/ca-risken/datasource-api/proto/code"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// CodeServiceServer is an autogenerated mock type for the CodeServiceServer type
type CodeServiceServer struct {
	mock.Mock
}

// DeleteDependencySetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) DeleteDependencySetting(_a0 context.Context, _a1 *code.DeleteDependencySettingRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.DeleteDependencySettingRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.DeleteDependencySettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGitHubEnterpriseOrg provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) DeleteGitHubEnterpriseOrg(_a0 context.Context, _a1 *code.DeleteGitHubEnterpriseOrgRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.DeleteGitHubEnterpriseOrgRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.DeleteGitHubEnterpriseOrgRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGitHubSetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) DeleteGitHubSetting(_a0 context.Context, _a1 *code.DeleteGitHubSettingRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.DeleteGitHubSettingRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.DeleteGitHubSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteGitleaksSetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) DeleteGitleaksSetting(_a0 context.Context, _a1 *code.DeleteGitleaksSettingRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.DeleteGitleaksSettingRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.DeleteGitleaksSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGitHubSetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) GetGitHubSetting(_a0 context.Context, _a1 *code.GetGitHubSettingRequest) (*code.GetGitHubSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.GetGitHubSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.GetGitHubSettingRequest) *code.GetGitHubSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.GetGitHubSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.GetGitHubSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanAll provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) InvokeScanAll(_a0 context.Context, _a1 *emptypb.Empty) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanDependency provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) InvokeScanDependency(_a0 context.Context, _a1 *code.InvokeScanDependencyRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.InvokeScanDependencyRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.InvokeScanDependencyRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanGitleaks provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) InvokeScanGitleaks(_a0 context.Context, _a1 *code.InvokeScanGitleaksRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *code.InvokeScanGitleaksRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.InvokeScanGitleaksRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDataSource provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) ListDataSource(_a0 context.Context, _a1 *code.ListDataSourceRequest) (*code.ListDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.ListDataSourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.ListDataSourceRequest) *code.ListDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.ListDataSourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.ListDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGitHubEnterpriseOrg provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) ListGitHubEnterpriseOrg(_a0 context.Context, _a1 *code.ListGitHubEnterpriseOrgRequest) (*code.ListGitHubEnterpriseOrgResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.ListGitHubEnterpriseOrgResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.ListGitHubEnterpriseOrgRequest) *code.ListGitHubEnterpriseOrgResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.ListGitHubEnterpriseOrgResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.ListGitHubEnterpriseOrgRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListGitHubSetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) ListGitHubSetting(_a0 context.Context, _a1 *code.ListGitHubSettingRequest) (*code.ListGitHubSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.ListGitHubSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.ListGitHubSettingRequest) *code.ListGitHubSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.ListGitHubSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.ListGitHubSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutDependencySetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) PutDependencySetting(_a0 context.Context, _a1 *code.PutDependencySettingRequest) (*code.PutDependencySettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.PutDependencySettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutDependencySettingRequest) *code.PutDependencySettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutDependencySettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutDependencySettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutGitHubEnterpriseOrg provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) PutGitHubEnterpriseOrg(_a0 context.Context, _a1 *code.PutGitHubEnterpriseOrgRequest) (*code.PutGitHubEnterpriseOrgResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.PutGitHubEnterpriseOrgResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutGitHubEnterpriseOrgRequest) *code.PutGitHubEnterpriseOrgResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutGitHubEnterpriseOrgResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutGitHubEnterpriseOrgRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutGitHubSetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) PutGitHubSetting(_a0 context.Context, _a1 *code.PutGitHubSettingRequest) (*code.PutGitHubSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.PutGitHubSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutGitHubSettingRequest) *code.PutGitHubSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutGitHubSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutGitHubSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutGitleaksSetting provides a mock function with given fields: _a0, _a1
func (_m *CodeServiceServer) PutGitleaksSetting(_a0 context.Context, _a1 *code.PutGitleaksSettingRequest) (*code.PutGitleaksSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *code.PutGitleaksSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *code.PutGitleaksSettingRequest) *code.PutGitleaksSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*code.PutGitleaksSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *code.PutGitleaksSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCodeServiceServer creates a new instance of CodeServiceServer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewCodeServiceServer(t testing.TB) *CodeServiceServer {
	mock := &CodeServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}