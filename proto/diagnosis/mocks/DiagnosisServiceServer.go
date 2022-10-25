// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	context "context"

	diagnosis "github.com/ca-risken/datasource-api/proto/diagnosis"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// DiagnosisServiceServer is an autogenerated mock type for the DiagnosisServiceServer type
type DiagnosisServiceServer struct {
	mock.Mock
}

// DeleteApplicationScan provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) DeleteApplicationScan(_a0 context.Context, _a1 *diagnosis.DeleteApplicationScanRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.DeleteApplicationScanRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.DeleteApplicationScanRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) DeleteApplicationScanBasicSetting(_a0 context.Context, _a1 *diagnosis.DeleteApplicationScanBasicSettingRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.DeleteApplicationScanBasicSettingRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.DeleteApplicationScanBasicSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDiagnosisDataSource provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) DeleteDiagnosisDataSource(_a0 context.Context, _a1 *diagnosis.DeleteDiagnosisDataSourceRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.DeleteDiagnosisDataSourceRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.DeleteDiagnosisDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePortscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) DeletePortscanSetting(_a0 context.Context, _a1 *diagnosis.DeletePortscanSettingRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.DeletePortscanSettingRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.DeletePortscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePortscanTarget provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) DeletePortscanTarget(_a0 context.Context, _a1 *diagnosis.DeletePortscanTargetRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.DeletePortscanTargetRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.DeletePortscanTargetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteWpscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) DeleteWpscanSetting(_a0 context.Context, _a1 *diagnosis.DeleteWpscanSettingRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.DeleteWpscanSettingRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.DeleteWpscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApplicationScan provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) GetApplicationScan(_a0 context.Context, _a1 *diagnosis.GetApplicationScanRequest) (*diagnosis.GetApplicationScanResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.GetApplicationScanResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.GetApplicationScanRequest) *diagnosis.GetApplicationScanResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.GetApplicationScanResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.GetApplicationScanRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) GetApplicationScanBasicSetting(_a0 context.Context, _a1 *diagnosis.GetApplicationScanBasicSettingRequest) (*diagnosis.GetApplicationScanBasicSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.GetApplicationScanBasicSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.GetApplicationScanBasicSettingRequest) *diagnosis.GetApplicationScanBasicSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.GetApplicationScanBasicSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.GetApplicationScanBasicSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDiagnosisDataSource provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) GetDiagnosisDataSource(_a0 context.Context, _a1 *diagnosis.GetDiagnosisDataSourceRequest) (*diagnosis.GetDiagnosisDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.GetDiagnosisDataSourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.GetDiagnosisDataSourceRequest) *diagnosis.GetDiagnosisDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.GetDiagnosisDataSourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.GetDiagnosisDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) GetPortscanSetting(_a0 context.Context, _a1 *diagnosis.GetPortscanSettingRequest) (*diagnosis.GetPortscanSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.GetPortscanSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.GetPortscanSettingRequest) *diagnosis.GetPortscanSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.GetPortscanSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.GetPortscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortscanTarget provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) GetPortscanTarget(_a0 context.Context, _a1 *diagnosis.GetPortscanTargetRequest) (*diagnosis.GetPortscanTargetResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.GetPortscanTargetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.GetPortscanTargetRequest) *diagnosis.GetPortscanTargetResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.GetPortscanTargetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.GetPortscanTargetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWpscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) GetWpscanSetting(_a0 context.Context, _a1 *diagnosis.GetWpscanSettingRequest) (*diagnosis.GetWpscanSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.GetWpscanSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.GetWpscanSettingRequest) *diagnosis.GetWpscanSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.GetWpscanSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.GetWpscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScan provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) InvokeScan(_a0 context.Context, _a1 *diagnosis.InvokeScanRequest) (*diagnosis.InvokeScanResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.InvokeScanResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.InvokeScanRequest) *diagnosis.InvokeScanResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.InvokeScanResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.InvokeScanRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeScanAll provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) InvokeScanAll(_a0 context.Context, _a1 *diagnosis.InvokeScanAllRequest) (*emptypb.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.InvokeScanAllRequest) *emptypb.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.InvokeScanAllRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListApplicationScan provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) ListApplicationScan(_a0 context.Context, _a1 *diagnosis.ListApplicationScanRequest) (*diagnosis.ListApplicationScanResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.ListApplicationScanResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.ListApplicationScanRequest) *diagnosis.ListApplicationScanResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.ListApplicationScanResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.ListApplicationScanRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) ListApplicationScanBasicSetting(_a0 context.Context, _a1 *diagnosis.ListApplicationScanBasicSettingRequest) (*diagnosis.ListApplicationScanBasicSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.ListApplicationScanBasicSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.ListApplicationScanBasicSettingRequest) *diagnosis.ListApplicationScanBasicSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.ListApplicationScanBasicSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.ListApplicationScanBasicSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDiagnosisDataSource provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) ListDiagnosisDataSource(_a0 context.Context, _a1 *diagnosis.ListDiagnosisDataSourceRequest) (*diagnosis.ListDiagnosisDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.ListDiagnosisDataSourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.ListDiagnosisDataSourceRequest) *diagnosis.ListDiagnosisDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.ListDiagnosisDataSourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.ListDiagnosisDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPortscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) ListPortscanSetting(_a0 context.Context, _a1 *diagnosis.ListPortscanSettingRequest) (*diagnosis.ListPortscanSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.ListPortscanSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.ListPortscanSettingRequest) *diagnosis.ListPortscanSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.ListPortscanSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.ListPortscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPortscanTarget provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) ListPortscanTarget(_a0 context.Context, _a1 *diagnosis.ListPortscanTargetRequest) (*diagnosis.ListPortscanTargetResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.ListPortscanTargetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.ListPortscanTargetRequest) *diagnosis.ListPortscanTargetResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.ListPortscanTargetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.ListPortscanTargetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListWpscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) ListWpscanSetting(_a0 context.Context, _a1 *diagnosis.ListWpscanSettingRequest) (*diagnosis.ListWpscanSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.ListWpscanSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.ListWpscanSettingRequest) *diagnosis.ListWpscanSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.ListWpscanSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.ListWpscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutApplicationScan provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) PutApplicationScan(_a0 context.Context, _a1 *diagnosis.PutApplicationScanRequest) (*diagnosis.PutApplicationScanResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.PutApplicationScanResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.PutApplicationScanRequest) *diagnosis.PutApplicationScanResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.PutApplicationScanResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.PutApplicationScanRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutApplicationScanBasicSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) PutApplicationScanBasicSetting(_a0 context.Context, _a1 *diagnosis.PutApplicationScanBasicSettingRequest) (*diagnosis.PutApplicationScanBasicSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.PutApplicationScanBasicSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.PutApplicationScanBasicSettingRequest) *diagnosis.PutApplicationScanBasicSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.PutApplicationScanBasicSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.PutApplicationScanBasicSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutDiagnosisDataSource provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) PutDiagnosisDataSource(_a0 context.Context, _a1 *diagnosis.PutDiagnosisDataSourceRequest) (*diagnosis.PutDiagnosisDataSourceResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.PutDiagnosisDataSourceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.PutDiagnosisDataSourceRequest) *diagnosis.PutDiagnosisDataSourceResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.PutDiagnosisDataSourceResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.PutDiagnosisDataSourceRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPortscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) PutPortscanSetting(_a0 context.Context, _a1 *diagnosis.PutPortscanSettingRequest) (*diagnosis.PutPortscanSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.PutPortscanSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.PutPortscanSettingRequest) *diagnosis.PutPortscanSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.PutPortscanSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.PutPortscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPortscanTarget provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) PutPortscanTarget(_a0 context.Context, _a1 *diagnosis.PutPortscanTargetRequest) (*diagnosis.PutPortscanTargetResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.PutPortscanTargetResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.PutPortscanTargetRequest) *diagnosis.PutPortscanTargetResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.PutPortscanTargetResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.PutPortscanTargetRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutWpscanSetting provides a mock function with given fields: _a0, _a1
func (_m *DiagnosisServiceServer) PutWpscanSetting(_a0 context.Context, _a1 *diagnosis.PutWpscanSettingRequest) (*diagnosis.PutWpscanSettingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *diagnosis.PutWpscanSettingResponse
	if rf, ok := ret.Get(0).(func(context.Context, *diagnosis.PutWpscanSettingRequest) *diagnosis.PutWpscanSettingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*diagnosis.PutWpscanSettingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *diagnosis.PutWpscanSettingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDiagnosisServiceServer creates a new instance of DiagnosisServiceServer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewDiagnosisServiceServer(t testing.TB) *DiagnosisServiceServer {
	mock := &DiagnosisServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}