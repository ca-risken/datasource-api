// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/ca-risken/datasource-api/pkg/model"
)

// DiagnosisRepoInterface is an autogenerated mock type for the DiagnosisRepoInterface type
type DiagnosisRepoInterface struct {
	mock.Mock
}

// DeleteApplicationScan provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) DeleteApplicationScan(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) DeleteApplicationScanBasicSetting(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDiagnosisDataSource provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) DeleteDiagnosisDataSource(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePortscanSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) DeletePortscanSetting(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePortscanTarget provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) DeletePortscanTarget(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePortscanTargetByPortscanSettingID provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) DeletePortscanTargetByPortscanSettingID(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWpscanSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) DeleteWpscanSetting(_a0 context.Context, _a1 uint32, _a2 uint32) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetApplicationScan provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) GetApplicationScan(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.ApplicationScan, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.ApplicationScan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.ApplicationScan, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.ApplicationScan); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationScan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) GetApplicationScanBasicSetting(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.ApplicationScanBasicSetting, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.ApplicationScanBasicSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.ApplicationScanBasicSetting, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.ApplicationScanBasicSetting); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationScanBasicSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDiagnosisDataSource provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) GetDiagnosisDataSource(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.DiagnosisDataSource, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.DiagnosisDataSource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.DiagnosisDataSource, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.DiagnosisDataSource); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DiagnosisDataSource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortscanSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) GetPortscanSetting(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.PortscanSetting, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.PortscanSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.PortscanSetting, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.PortscanSetting); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PortscanSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortscanTarget provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) GetPortscanTarget(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.PortscanTarget, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.PortscanTarget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.PortscanTarget, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.PortscanTarget); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PortscanTarget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWpscanSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) GetWpscanSetting(_a0 context.Context, _a1 uint32, _a2 uint32) (*model.WpscanSetting, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.WpscanSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*model.WpscanSetting, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *model.WpscanSetting); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.WpscanSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllWpscanSetting provides a mock function with given fields: _a0
func (_m *DiagnosisRepoInterface) ListAllWpscanSetting(_a0 context.Context) (*[]model.WpscanSetting, error) {
	ret := _m.Called(_a0)

	var r0 *[]model.WpscanSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*[]model.WpscanSetting, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *[]model.WpscanSetting); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.WpscanSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListApplicationScan provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) ListApplicationScan(_a0 context.Context, _a1 uint32, _a2 uint32) (*[]model.ApplicationScan, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.ApplicationScan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*[]model.ApplicationScan, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *[]model.ApplicationScan); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.ApplicationScan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) ListApplicationScanBasicSetting(_a0 context.Context, _a1 uint32, _a2 uint32) (*[]model.ApplicationScanBasicSetting, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.ApplicationScanBasicSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*[]model.ApplicationScanBasicSetting, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *[]model.ApplicationScanBasicSetting); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.ApplicationScanBasicSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDiagnosisDataSource provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) ListDiagnosisDataSource(_a0 context.Context, _a1 uint32, _a2 string) (*[]model.DiagnosisDataSource, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.DiagnosisDataSource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) (*[]model.DiagnosisDataSource, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, string) *[]model.DiagnosisDataSource); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.DiagnosisDataSource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPortscanSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) ListPortscanSetting(_a0 context.Context, _a1 uint32, _a2 uint32) (*[]model.PortscanSetting, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.PortscanSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*[]model.PortscanSetting, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *[]model.PortscanSetting); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.PortscanSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPortscanTarget provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) ListPortscanTarget(_a0 context.Context, _a1 uint32, _a2 uint32) (*[]model.PortscanTarget, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.PortscanTarget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*[]model.PortscanTarget, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *[]model.PortscanTarget); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.PortscanTarget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListWpscanSetting provides a mock function with given fields: _a0, _a1, _a2
func (_m *DiagnosisRepoInterface) ListWpscanSetting(_a0 context.Context, _a1 uint32, _a2 uint32) (*[]model.WpscanSetting, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *[]model.WpscanSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) (*[]model.WpscanSetting, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32, uint32) *[]model.WpscanSetting); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]model.WpscanSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32, uint32) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertApplicationScan provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisRepoInterface) UpsertApplicationScan(_a0 context.Context, _a1 *model.ApplicationScan) (*model.ApplicationScan, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.ApplicationScan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.ApplicationScan) (*model.ApplicationScan, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.ApplicationScan) *model.ApplicationScan); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationScan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.ApplicationScan) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisRepoInterface) UpsertApplicationScanBasicSetting(_a0 context.Context, _a1 *model.ApplicationScanBasicSetting) (*model.ApplicationScanBasicSetting, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.ApplicationScanBasicSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.ApplicationScanBasicSetting) (*model.ApplicationScanBasicSetting, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.ApplicationScanBasicSetting) *model.ApplicationScanBasicSetting); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationScanBasicSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.ApplicationScanBasicSetting) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertDiagnosisDataSource provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisRepoInterface) UpsertDiagnosisDataSource(_a0 context.Context, _a1 *model.DiagnosisDataSource) (*model.DiagnosisDataSource, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.DiagnosisDataSource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DiagnosisDataSource) (*model.DiagnosisDataSource, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.DiagnosisDataSource) *model.DiagnosisDataSource); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DiagnosisDataSource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.DiagnosisDataSource) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertPortscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisRepoInterface) UpsertPortscanSetting(_a0 context.Context, _a1 *model.PortscanSetting) (*model.PortscanSetting, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.PortscanSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.PortscanSetting) (*model.PortscanSetting, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.PortscanSetting) *model.PortscanSetting); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PortscanSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.PortscanSetting) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertPortscanTarget provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisRepoInterface) UpsertPortscanTarget(_a0 context.Context, _a1 *model.PortscanTarget) (*model.PortscanTarget, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.PortscanTarget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.PortscanTarget) (*model.PortscanTarget, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.PortscanTarget) *model.PortscanTarget); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PortscanTarget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.PortscanTarget) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertWpscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisRepoInterface) UpsertWpscanSetting(_a0 context.Context, _a1 *model.WpscanSetting) (*model.WpscanSetting, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.WpscanSetting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.WpscanSetting) (*model.WpscanSetting, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.WpscanSetting) *model.WpscanSetting); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.WpscanSetting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.WpscanSetting) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDiagnosisRepoInterface creates a new instance of DiagnosisRepoInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDiagnosisRepoInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DiagnosisRepoInterface {
	mock := &DiagnosisRepoInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}