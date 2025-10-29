package diagnosis

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
)

// DiagnosisDataSource

// Validate ListDiagnosisDataSourceRequest
func (r *ListDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetDiagnosisDataSourceRequest
func (r *GetDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
	)
}

// Validate PutDiagnosisDataSourceRequest
func (r *PutDiagnosisDataSourceRequest) Validate() error {
	if r.DiagnosisDataSource == nil {
		return errors.New("required DiagnosisDataSource")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	); err != nil {
		return err
	}
	return r.DiagnosisDataSource.Validate()
}

// Validate DeleteDataSourceRequest
func (r *DeleteDiagnosisDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
	)
}

// WpscanSetting

// Validate ListWpscanSettingRequest
func (r *ListWpscanSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetWpscanSettingRequest
func (r *GetWpscanSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.WpscanSettingId, validation.Required),
	)
}

// Validate PutWpscanSettingRequest
func (r *PutWpscanSettingRequest) Validate() error {
	if r.WpscanSetting == nil {
		return errors.New("required WpscanSetting")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.WpscanSetting.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.WpscanSetting.Validate()
}

// Validate DeleteWpscanSettingRequest
func (r *DeleteWpscanSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.WpscanSettingId, validation.Required),
	)
}

// PortscanSetting

// Validate ListPortscanSettingRequest
func (r *ListPortscanSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetPortscanSettingRequest
func (r *GetPortscanSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.PortscanSettingId, validation.Required),
	)
}

// Validate PutPortscanSettingRequest
func (r *PutPortscanSettingRequest) Validate() error {
	if r.PortscanSetting == nil {
		return errors.New("required PortscanSetting")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.PortscanSetting.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.PortscanSetting.Validate()
}

// Validate DeletePortscanSettingRequest
func (r *DeletePortscanSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.PortscanSettingId, validation.Required),
	)
}

// PortscanTarget

// Validate ListPortscanTargetRequest
func (r *ListPortscanTargetRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetPortscanTargetRequest
func (r *GetPortscanTargetRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.PortscanTargetId, validation.Required),
	)
}

// Validate PutPortscanTargetRequest
func (r *PutPortscanTargetRequest) Validate() error {
	if r.PortscanTarget == nil {
		return errors.New("required PortscanTarget")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.PortscanTarget.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.PortscanTarget.Validate()
}

// Validate DeletePortscanTargetRequest
func (r *DeletePortscanTargetRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.PortscanTargetId, validation.Required),
	)
}

// ApplicationScan

// Validate ListApplicationScanRequest
func (r *ListApplicationScanRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetApplicationScanRequest
func (r *GetApplicationScanRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.ApplicationScanId, validation.Required),
	)
}

// Validate PutApplicationScanRequest
func (r *PutApplicationScanRequest) Validate() error {
	if r.ApplicationScan == nil {
		return errors.New("required ApplicationScan")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.ApplicationScan.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.ApplicationScan.Validate()
}

// Validate DeleteApplicationScanRequest
func (r *DeleteApplicationScanRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.ApplicationScanId, validation.Required),
	)
}

// ApplicationScanBasicSetting

// Validate ListApplicationScanBasicSettingRequest
func (r *ListApplicationScanBasicSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetApplicationScanBasicSettingRequest
func (r *GetApplicationScanBasicSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.ApplicationScanId, validation.Required),
	)
}

// Validate PutApplicationScanBasicSettingRequest
func (r *PutApplicationScanBasicSettingRequest) Validate() error {
	if r.ApplicationScanBasicSetting == nil {
		return errors.New("required ApplicationScanBasicSetting")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.ApplicationScanBasicSetting.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.ApplicationScanBasicSetting.Validate()
}

// Validate DeleteApplicationScanBasicSettingRequest
func (r *DeleteApplicationScanBasicSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.ApplicationScanBasicSettingId, validation.Required),
	)
}

// Validate InvokeScanRequest
func (r *InvokeScanRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.SettingId, validation.Required),
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
	)
}

/**
 * Entity
**/

// Validate DiagnosisDataSourceForUpsert
func (d *DiagnosisDataSourceForUpsert) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Name, validation.Required, validation.Length(0, 50)),
		validation.Field(&d.Description, validation.Required, validation.Length(0, 200)),
		validation.Field(&d.MaxScore, validation.Required, validation.Min(0.0), validation.Max(99999.0)),
	)
}

// Validate WpscanSettingForUpsert
func (r *WpscanSettingForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.TargetUrl, validation.Required, validation.Length(0, 200)),
		validation.Field(&r.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.StatusDetail, validation.Length(0, 255)),
		validation.Field(&r.Options, is.JSON),
	)
}

// Validate PortscanSettingForUpsert
func (r *PortscanSettingForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(0, 200)),
	)
}

// Validate PortscanTargetForUpsert
func (r *PortscanTargetForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.PortscanSettingId, validation.Required),
		validation.Field(&r.Target, validation.Required, validation.Length(0, 300)),
		validation.Field(&r.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.StatusDetail, validation.Length(0, 255)),
	)
}

// Validate ApplicationScanForUpsert
func (r *ApplicationScanForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.DiagnosisDataSourceId, validation.Required),
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(0, 200)),
		validation.Field(&r.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.StatusDetail, validation.Length(0, 255)),
	)
}

// Validate ApplicationScanBasicSettingForUpsert
func (r *ApplicationScanBasicSettingForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.ApplicationScanId, validation.Required),
		validation.Field(&r.Target, validation.Required, validation.Length(0, 255)),
		validation.Field(&r.MaxDepth, validation.Min(uint32(0)), validation.Max(uint32(100))),
		validation.Field(&r.MaxChildren, validation.Min(uint32(0)), validation.Max(uint32(100))),
	)
}
