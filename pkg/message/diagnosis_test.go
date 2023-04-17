package message

import (
	"reflect"
	"testing"
)

func TestValidateWPScan(t *testing.T) {
	cases := []struct {
		name    string
		input   *WpscanQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: &WpscanQueueMessage{WpscanSettingID: 1, ProjectID: 1, TargetURL: "http://example.com"},
		},
		{
			name:  "OK(scan_only)",
			input: &WpscanQueueMessage{WpscanSettingID: 1, ProjectID: 1, TargetURL: "http://example.com", ScanOnly: true},
		},
		{
			name:    "NG Required(wpscan_setting_id)",
			input:   &WpscanQueueMessage{ProjectID: 1, TargetURL: "http://example.com"},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &WpscanQueueMessage{WpscanSettingID: 1, TargetURL: "http://example.com"},
			wantErr: true,
		},
		{
			name:    "NG Required(target_url)",
			input:   &WpscanQueueMessage{WpscanSettingID: 1, ProjectID: 1},
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

func TestParseWPScanMessage(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *WpscanQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: `{"wpscan_setting_id":1, "project_id":1, "target_url":"http://example.com"}`,
			want:  &WpscanQueueMessage{WpscanSettingID: 1, ProjectID: 1, TargetURL: "http://example.com"},
		},
		{
			name:  "OK(scan_only)",
			input: `{"wpscan_setting_id":1, "project_id":1, "target_url":"http://example.com", "scan_only":"true"}`,
			want:  &WpscanQueueMessage{WpscanSettingID: 1, ProjectID: 1, TargetURL: "http://example.com", ScanOnly: true},
		},
		{
			name:    "NG Json parse error",
			input:   `{"parse...: error`,
			wantErr: true,
		},
		{
			name:    "NG Invalid mmessage(required parammeter)",
			input:   `{}`,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseWpscanMessage(c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error occured, wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected response, want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestValidatePortscan(t *testing.T) {
	cases := []struct {
		name    string
		input   *PortscanQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: &PortscanQueueMessage{PortscanSettingID: 1, PortscanTargetID: 1, ProjectID: 1, Target: "example.com"},
		},
		{
			name:  "OK(scan_only)",
			input: &PortscanQueueMessage{PortscanSettingID: 1, PortscanTargetID: 1, ProjectID: 1, Target: "example.com", ScanOnly: true},
		},
		{
			name:    "NG Required(portscan_setting_id)",
			input:   &PortscanQueueMessage{PortscanTargetID: 1, ProjectID: 1, Target: "example.com"},
			wantErr: true,
		},
		{
			name:    "NG Required(portscan_target_id)",
			input:   &PortscanQueueMessage{PortscanSettingID: 1, ProjectID: 1, Target: "example.com"},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &PortscanQueueMessage{PortscanSettingID: 1, PortscanTargetID: 1, Target: "example.com"},
			wantErr: true,
		},
		{
			name:    "NG Required(target)",
			input:   &PortscanQueueMessage{PortscanSettingID: 1, PortscanTargetID: 1, ProjectID: 1},
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

func TestParseportscanMessage(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *PortscanQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: `{"portscan_setting_id":1, "portscan_target_id":1, "project_id":1, "target":"example.com"}`,
			want:  &PortscanQueueMessage{PortscanSettingID: 1, PortscanTargetID: 1, ProjectID: 1, Target: "example.com"},
		},
		{
			name:  "OK(scan_only)",
			input: `{"portscan_setting_id":1, "portscan_target_id":1, "project_id":1, "target":"example.com", "scan_only":"true"}`,
			want:  &PortscanQueueMessage{PortscanSettingID: 1, PortscanTargetID: 1, ProjectID: 1, Target: "example.com", ScanOnly: true},
		},
		{
			name:    "NG Json parse error",
			input:   `{"parse...: error`,
			wantErr: true,
		},
		{
			name:    "NG Invalid mmessage(required parammeter)",
			input:   `{}`,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParsePortscanMessage(c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error occured, wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected response, want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestValidateApplicationScan(t *testing.T) {
	cases := []struct {
		name    string
		input   *ApplicationScanQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ApplicationScanQueueMessage{ApplicationScanID: 1, ProjectID: 1, Name: "test", ApplicationScanType: "BASIC"},
		},
		{
			name:  "OK(scan_only)",
			input: &ApplicationScanQueueMessage{ApplicationScanID: 1, ProjectID: 1, Name: "test", ApplicationScanType: "BASIC", ScanOnly: true},
		},
		{
			name:    "NG Required(application_scan_id)",
			input:   &ApplicationScanQueueMessage{ProjectID: 1, ApplicationScanType: "BASIC", Name: "test"},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ApplicationScanQueueMessage{ApplicationScanID: 1, Name: "test", ApplicationScanType: "BASIC"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &ApplicationScanQueueMessage{ApplicationScanID: 1, ProjectID: 1, ApplicationScanType: "BASIC"},
			wantErr: true,
		},
		{
			name:    "NG Required(target)",
			input:   &ApplicationScanQueueMessage{ApplicationScanID: 1, ProjectID: 1, Name: "test"},
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

func TestParseApplicationScanMessage(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *ApplicationScanQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: `{"application_scan_id":1,"project_id":1, "name":"test", "application_scan_type": "BASIC"}`,
			want:  &ApplicationScanQueueMessage{ApplicationScanID: 1, ProjectID: 1, Name: "test", ApplicationScanType: "BASIC"},
		},
		{
			name:  "OK(scan_only)",
			input: `{"application_scan_id":1,"project_id":1, "name":"test", "application_scan_type": "BASIC", "scan_only":"true"}`,
			want:  &ApplicationScanQueueMessage{ApplicationScanID: 1, ProjectID: 1, Name: "test", ApplicationScanType: "BASIC", ScanOnly: true},
		},
		{
			name:    "NG Json parse error",
			input:   `{"parse...: error`,
			wantErr: true,
		},
		{
			name:    "NG Invalid mmessage(required parammeter)",
			input:   `{}`,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseApplicationScanMessage(c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error occured, wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected response, want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
