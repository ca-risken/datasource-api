package diagnosis

import (
	"testing"
)

//DiagnosisDataSource DataSource

func TestValidate_ListDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListDiagnosisDataSourceRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetDiagnosisDataSourceRequest{ProjectId: 1, DiagnosisDataSourceId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetDiagnosisDataSourceRequest{DiagnosisDataSourceId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(diagnosis_data_source_id)",
			input:   &GetDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutDiagnosisDataSourceRequest{ProjectId: 1, DiagnosisDataSource: &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 10.0}},
			wantErr: false,
		},
		{
			name:    "NG Required(DiagnosisDataSource)",
			input:   &PutDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutDiagnosisDataSourceRequest{DiagnosisDataSource: &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 10.0}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteDiagnosisDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteDiagnosisDataSourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteDiagnosisDataSourceRequest{ProjectId: 1, DiagnosisDataSourceId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteDiagnosisDataSourceRequest{DiagnosisDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &DeleteDiagnosisDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

//WpscanSetting DataSource

func TestValidate_ListWpscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListWpscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListWpscanSettingRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListWpscanSettingRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetWpscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetWpscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetWpscanSettingRequest{ProjectId: 1, WpscanSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetWpscanSettingRequest{WpscanSettingId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(wpscan_setting_id)",
			input:   &GetWpscanSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutWpscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutWpscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutWpscanSettingRequest{ProjectId: 1001, WpscanSetting: &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, TargetUrl: "hoge_target"}},
			wantErr: false,
		},
		{
			name:    "NG Required(WpscanSetting)",
			input:   &PutWpscanSettingRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != wpscan_setting.project_id)",
			input:   &PutWpscanSettingRequest{ProjectId: 1001, WpscanSetting: &WpscanSettingForUpsert{ProjectId: 1002, DiagnosisDataSourceId: 1, TargetUrl: "hoge_target"}},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutWpscanSettingRequest{WpscanSetting: &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, TargetUrl: "hoge_target"}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteWpscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteWpscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteWpscanSettingRequest{ProjectId: 1, WpscanSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteWpscanSettingRequest{WpscanSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &DeleteWpscanSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

//PortscanSetting DataSource

func TestValidate_ListPortscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListPortscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListPortscanSettingRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListPortscanSettingRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetPortscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetPortscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetPortscanSettingRequest{ProjectId: 1, PortscanSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetPortscanSettingRequest{PortscanSettingId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(portscan_setting_id)",
			input:   &GetPortscanSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutPortscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutPortscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutPortscanSettingRequest{ProjectId: 1001, PortscanSetting: &PortscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_target"}},
			wantErr: false,
		},
		{
			name:    "NG Required(PortscanSetting)",
			input:   &PutPortscanSettingRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != portscan_setting.project_id)",
			input:   &PutPortscanSettingRequest{ProjectId: 1001, PortscanSetting: &PortscanSettingForUpsert{ProjectId: 1002, DiagnosisDataSourceId: 1, Name: "hoge_target"}},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutPortscanSettingRequest{PortscanSetting: &PortscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_target"}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeletePortscanSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeletePortscanSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeletePortscanSettingRequest{ProjectId: 1, PortscanSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeletePortscanSettingRequest{PortscanSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &DeletePortscanSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

//PortscanTarget DataSource

func TestValidate_ListPortscanTargetRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListPortscanTargetRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListPortscanTargetRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListPortscanTargetRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetPortscanTargetRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetPortscanTargetRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetPortscanTargetRequest{ProjectId: 1, PortscanTargetId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetPortscanTargetRequest{PortscanTargetId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(portscan_target_id)",
			input:   &GetPortscanTargetRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutPortscanTargetRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutPortscanTargetRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutPortscanTargetRequest{ProjectId: 1001, PortscanTarget: &PortscanTargetForUpsert{ProjectId: 1001, PortscanSettingId: 1, PortscanTargetId: 1, Target: "hoge_target"}},
			wantErr: false,
		},
		{
			name:    "NG Required(PortscanTarget)",
			input:   &PutPortscanTargetRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != portscan_setting.project_id)",
			input:   &PutPortscanTargetRequest{ProjectId: 1001, PortscanTarget: &PortscanTargetForUpsert{ProjectId: 1002, PortscanSettingId: 1, PortscanTargetId: 1, Target: "hoge_target"}},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutPortscanTargetRequest{PortscanTarget: &PortscanTargetForUpsert{ProjectId: 1001, PortscanSettingId: 1, PortscanTargetId: 1, Target: "hoge_target"}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeletePortscanTargetRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeletePortscanTargetRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeletePortscanTargetRequest{ProjectId: 1, PortscanTargetId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeletePortscanTargetRequest{PortscanTargetId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &DeletePortscanTargetRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

//ApplicationScan DataSource

func TestValidate_ListApplicationScanRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListApplicationScanRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListApplicationScanRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListApplicationScanRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetApplicationScanRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetApplicationScanRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetApplicationScanRequest{ProjectId: 1, ApplicationScanId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetApplicationScanRequest{ApplicationScanId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(portscan_setting_id)",
			input:   &GetApplicationScanRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutApplicationScanRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutApplicationScanRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutApplicationScanRequest{ProjectId: 1001, ApplicationScan: &ApplicationScanForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_target", ScanAt: 0}},
			wantErr: false,
		},
		{
			name:    "NG Required(ApplicationScan)",
			input:   &PutApplicationScanRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != portscan_setting.project_id)",
			input:   &PutApplicationScanRequest{ProjectId: 1001, ApplicationScan: &ApplicationScanForUpsert{ProjectId: 1002, DiagnosisDataSourceId: 1, Name: "hoge_target", ScanAt: 0}},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutApplicationScanRequest{ApplicationScan: &ApplicationScanForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_target", ScanAt: 0}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteApplicationScanRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteApplicationScanRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteApplicationScanRequest{ProjectId: 1, ApplicationScanId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteApplicationScanRequest{ApplicationScanId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(application_scan_id)",
			input:   &DeleteApplicationScanRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

//ApplicationScanBasicSetting DataSource

func TestValidate_ListApplicationScanBasicSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListApplicationScanBasicSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListApplicationScanBasicSettingRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &ListApplicationScanBasicSettingRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetApplicationScanBasicSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetApplicationScanBasicSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetApplicationScanBasicSettingRequest{ProjectId: 1, ApplicationScanId: 2},
			wantErr: false,
		},
		{
			name:    "NG required(project_id)",
			input:   &GetApplicationScanBasicSettingRequest{ApplicationScanId: 2},
			wantErr: true,
		},
		{
			name:    "NG required(portscan_target_id)",
			input:   &GetApplicationScanBasicSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutApplicationScanBasicSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutApplicationScanBasicSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutApplicationScanBasicSettingRequest{ProjectId: 1001, ApplicationScanBasicSetting: &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, ApplicationScanId: 1, ApplicationScanBasicSettingId: 1, Target: "hoge_target"}},
			wantErr: false,
		},
		{
			name:    "NG Required(ApplicationScanBasicSetting)",
			input:   &PutApplicationScanBasicSettingRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != portscan_setting.project_id)",
			input:   &PutApplicationScanBasicSettingRequest{ProjectId: 1001, ApplicationScanBasicSetting: &ApplicationScanBasicSettingForUpsert{ProjectId: 1002, ApplicationScanId: 1, ApplicationScanBasicSettingId: 1, Target: "hoge_target"}},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectId)",
			input:   &PutApplicationScanBasicSettingRequest{ApplicationScanBasicSetting: &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, ApplicationScanId: 1, ApplicationScanBasicSettingId: 1, Target: "hoge_target"}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteApplicationScanBasicSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteApplicationScanBasicSettingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteApplicationScanBasicSettingRequest{ProjectId: 1, ApplicationScanBasicSettingId: 2},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteApplicationScanBasicSettingRequest{ApplicationScanBasicSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(application_scan_basic_setting_id)",
			input:   &DeleteApplicationScanBasicSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_InvokeScanRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *InvokeScanRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &InvokeScanRequest{ProjectId: 1, SettingId: 2, DiagnosisDataSourceId: 3},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &InvokeScanRequest{SettingId: 2, DiagnosisDataSourceId: 3},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &InvokeScanRequest{ProjectId: 1, DiagnosisDataSourceId: 3},
			wantErr: true,
		},
		{
			name:    "NG Required(setting_id)",
			input:   &InvokeScanRequest{ProjectId: 1, SettingId: 2},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DiagnosisDataSourceForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *DiagnosisDataSourceForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 100},
			wantErr: false,
		},
		{
			name:    "NG Required(name)",
			input:   &DiagnosisDataSourceForUpsert{Description: "description", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Length(description)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Required(description)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", MaxScore: 100},
			wantErr: true,
		},
		{
			name:    "NG Num Over(Max Score)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: 100000},
			wantErr: true,
		},
		{
			name:    "NG Num Under(Max Score)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description", MaxScore: -1.0},
			wantErr: true,
		},
		{
			name:    "NG Required(Max Score)",
			input:   &DiagnosisDataSourceForUpsert{Name: "name", Description: "description"},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_WpscanSettingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *WpscanSettingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, TargetUrl: "hoge_url"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &WpscanSettingForUpsert{DiagnosisDataSourceId: 1, TargetUrl: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &WpscanSettingForUpsert{ProjectId: 1001, TargetUrl: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "Too long(target_url)",
			input:   &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, TargetUrl: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"},
			wantErr: true,
		},
		{
			name:    "NG Required(target_url)",
			input:   &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Too small scan_at",
			input:   &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, TargetUrl: "hoge_url", ScanAt: -1},
			wantErr: true,
		},
		{
			name:    "NG Too large scan_at",
			input:   &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, TargetUrl: "hoge_url", ScanAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG not json options",
			input:   &WpscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, TargetUrl: "hoge_url", ScanAt: 253402268400, Options: "hogehoge"},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PortscanSettingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *PortscanSettingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PortscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_url"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &PortscanSettingForUpsert{DiagnosisDataSourceId: 1, Name: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &PortscanSettingForUpsert{ProjectId: 1001, Name: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "Too long(name)",
			input:   &PortscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &PortscanSettingForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PortscanTargetForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *PortscanTargetForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PortscanTargetForUpsert{ProjectId: 1001, PortscanSettingId: 1, Target: "hoge_url"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &PortscanTargetForUpsert{PortscanSettingId: 1, Target: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "NG Required(portscan_setting_id)",
			input:   &PortscanTargetForUpsert{ProjectId: 1001, Target: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "Too long(target)",
			input:   &PortscanTargetForUpsert{ProjectId: 1001, PortscanSettingId: 1, Target: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"},
			wantErr: true,
		},
		{
			name:    "NG Required(target)",
			input:   &PortscanTargetForUpsert{ProjectId: 1001, PortscanSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Too small scan_at",
			input:   &PortscanTargetForUpsert{ProjectId: 1001, PortscanSettingId: 1, Target: "hoge_url", ScanAt: -1},
			wantErr: true,
		},
		{
			name:    "NG Too large scan_at",
			input:   &PortscanTargetForUpsert{ProjectId: 1001, PortscanSettingId: 1, Target: "hoge_url", ScanAt: 253402268400},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_ApplicationScanForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *ApplicationScanForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ApplicationScanForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "hoge_url"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ApplicationScanForUpsert{DiagnosisDataSourceId: 1, Name: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "NG Required(diagnosis_data_source_id)",
			input:   &ApplicationScanForUpsert{ProjectId: 1001, Name: "hoge_url"},
			wantErr: true,
		},
		{
			name:    "Too long(name)",
			input:   &ApplicationScanForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, Name: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &ApplicationScanForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Too small scan_at",
			input:   &ApplicationScanForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, ScanAt: -1},
			wantErr: true,
		},
		{
			name:    "NG Too large max_depth",
			input:   &ApplicationScanForUpsert{ProjectId: 1001, DiagnosisDataSourceId: 1, ScanAt: 253402268400},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_ApplicationScanBasicSettingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *ApplicationScanBasicSettingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, ApplicationScanId: 1, Target: "hoge_url", MaxDepth: 10, MaxChildren: 10},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ApplicationScanBasicSettingForUpsert{ApplicationScanId: 1, Target: "hoge_url", MaxDepth: 10, MaxChildren: 10},
			wantErr: true,
		},
		{
			name:    "NG Required(portscan_setting_id)",
			input:   &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, Target: "hoge_url", MaxDepth: 10, MaxChildren: 10},
			wantErr: true,
		},
		{
			name:    "Too long(target)",
			input:   &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, ApplicationScanId: 1, Target: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", MaxDepth: 10, MaxChildren: 10},
			wantErr: true,
		},
		{
			name:    "NG Required(target)",
			input:   &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, ApplicationScanId: 1, MaxDepth: 10, MaxChildren: 10},
			wantErr: true,
		},
		{
			name:    "NG Too large max_depth",
			input:   &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, ApplicationScanId: 1, Target: "hoge_url", MaxDepth: 101, MaxChildren: 10},
			wantErr: true,
		},
		{
			name:    "NG Too large max_children",
			input:   &ApplicationScanBasicSettingForUpsert{ProjectId: 1001, ApplicationScanId: 1, Target: "hoge_url", MaxDepth: 10, MaxChildren: 101},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
