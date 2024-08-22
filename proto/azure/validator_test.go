package azure

import (
	"testing"
	"time"
)

const (
	stringLength65           = "12345678901234567890123456789012345678901234567890123456789012345"
	stringLength129          = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789"
	stringLength256          = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456"
	unixtime19691231T235959  = -1
	unixtime100000101T000000 = 253402268400
)

func TestValidate_ListAzureDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListAzureDataSourceRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListAzureDataSourceRequest{AzureDataSourceId: 1, Name: "name"},
		},
		{
			name:  "OK empty",
			input: &ListAzureDataSourceRequest{},
		},
		{
			name:    "NG length(name)",
			input:   &ListAzureDataSourceRequest{AzureDataSourceId: 1, Name: stringLength65},
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

func TestValidate_ListAzureRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListAzureRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListAzureRequest{ProjectId: 1, AzureId: 1, SubscriptionId: "pj"},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAzureRequest{ProjectId: 1, AzureId: 1, SubscriptionId: stringLength129},
			wantErr: true,
		},
		{
			name:    "NG Length(gcp_project_id)",
			input:   &ListAzureRequest{AzureId: 1, SubscriptionId: "pj"},
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

func TestValidate_GetAzureRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetAzureRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &GetAzureRequest{ProjectId: 1, AzureId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetAzureRequest{AzureId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(gcp_id)",
			input:   &GetAzureRequest{ProjectId: 1},
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

func TestValidate_PutAzureRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutAzureRequest
		wantErr bool
	}{
		{
			name: "OK",
			input: &PutAzureRequest{ProjectId: 1, Azure: &AzureForUpsert{
				Name: "name", ProjectId: 1, SubscriptionId: "1", VerificationCode: "12345678",
			}},
		},
		{
			name:    "NG No Azure param",
			input:   &PutAzureRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "NG Invalid project_id",
			input: &PutAzureRequest{ProjectId: 999, Azure: &AzureForUpsert{
				Name: "name", ProjectId: 1, SubscriptionId: "1",
			}},
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

func TestValidate_DeleteAzureRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteAzureRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DeleteAzureRequest{ProjectId: 1, AzureId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteAzureRequest{AzureId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(gcp_id)",
			input:   &DeleteAzureRequest{ProjectId: 1},
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

func TestValidate_ListRelAzureDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListRelAzureDataSourceRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListRelAzureDataSourceRequest{AzureId: 1},
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

func TestValidate_GetRelAzureDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetRelAzureDataSourceRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &GetRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetRelAzureDataSourceRequest{AzureId: 1, AzureDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(gcp_id)",
			input:   &GetRelAzureDataSourceRequest{ProjectId: 1, AzureDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(google_data_source_id)",
			input:   &GetRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1},
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

func TestValidate_AttachRelAzureDataSourceRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *AttachRelAzureDataSourceRequest
		wantErr bool
	}{
		{
			name: "OK",
			input: &AttachRelAzureDataSourceRequest{ProjectId: 1, RelAzureDataSource: &RelAzureDataSourceForUpsert{
				AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			}},
		},
		{
			name:    "NG No AzureDataSource param",
			input:   &AttachRelAzureDataSourceRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "NG Invalid project_id",
			input: &AttachRelAzureDataSourceRequest{ProjectId: 999, RelAzureDataSource: &RelAzureDataSourceForUpsert{
				AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			}},
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

func TestValidate_DetachRelAzureDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DetachRelAzureDataSourceRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DetachRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DetachRelAzureDataSourceRequest{AzureId: 1, AzureDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(gcp_id)",
			input:   &DetachRelAzureDataSourceRequest{ProjectId: 1, AzureDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(google_data_source_id)",
			input:   &DetachRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1},
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

func TestValidate_InvokeScanAzureRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *InvokeScanAzureRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &InvokeScanAzureRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &InvokeScanAzureRequest{AzureId: 1, AzureDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(gcp_id)",
			input:   &InvokeScanAzureRequest{ProjectId: 1, AzureDataSourceId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(google_data_source_id)",
			input:   &InvokeScanAzureRequest{ProjectId: 1, AzureId: 1},
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

func TestValidate_AzureForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *AzureForUpsert
		wantErr bool
	}{
		{
			name: "OK",
			input: &AzureForUpsert{
				AzureId: 1, Name: "name", ProjectId: 1, SubscriptionId: "my-pj", VerificationCode: "12345678",
			},
		},
		{
			name: "OK minimize",
			input: &AzureForUpsert{
				Name: "name", ProjectId: 1, SubscriptionId: "my-pj", VerificationCode: "12345678",
			},
		},
		{
			name: "NG Required(name)",
			input: &AzureForUpsert{
				AzureId: 1, ProjectId: 1, SubscriptionId: "my-pj", VerificationCode: "12345678",
			},
			wantErr: true,
		},
		{
			name: "NG Length(name)",
			input: &AzureForUpsert{
				AzureId: 1, Name: stringLength65, ProjectId: 1, SubscriptionId: "my-pj", VerificationCode: "12345678",
			},
			wantErr: true,
		},
		{
			name: "NG Required(project_id)",
			input: &AzureForUpsert{
				AzureId: 1, Name: "name", SubscriptionId: "my-pj", VerificationCode: "12345678",
			},
			wantErr: true,
		},
		{
			name: "NG Required(gcp_project_id)",
			input: &AzureForUpsert{
				AzureId: 1, Name: "name", ProjectId: 1, VerificationCode: "12345678",
			},
			wantErr: true,
		},
		{
			name: "NG Length(gcp_project_id)",
			input: &AzureForUpsert{
				AzureId: 1, Name: "name", ProjectId: 1, SubscriptionId: stringLength129, VerificationCode: "12345678",
			},
			wantErr: true,
		},
		{
			name: "NG Required(verification_code)",
			input: &AzureForUpsert{
				AzureId: 1, Name: "name", ProjectId: 1, SubscriptionId: "my-pj",
			},
			wantErr: true,
		},
		{
			name: "NG MinLength(verification_code)",
			input: &AzureForUpsert{
				AzureId: 1, Name: "name", ProjectId: 1, SubscriptionId: stringLength129, VerificationCode: "1234567",
			},
			wantErr: true,
		},
		{
			name: "NG MaxLength(verification_code)",
			input: &AzureForUpsert{
				AzureId: 1, Name: "name", ProjectId: 1, SubscriptionId: stringLength129, VerificationCode: stringLength256,
			},
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

func TestValidate_RelAzureDataSourceForUpsert(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *RelAzureDataSourceForUpsert
		wantErr bool
	}{
		{
			name: "OK",
			input: &RelAzureDataSourceForUpsert{
				AzureDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
		},
		{
			name: "OK minimize",
			input: &RelAzureDataSourceForUpsert{
				AzureDataSourceId: 1, ProjectId: 1,
			},
		},
		{
			name: "NG Required(google_data_source_id)",
			input: &RelAzureDataSourceForUpsert{
				ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Required(project_id)",
			input: &RelAzureDataSourceForUpsert{
				AzureDataSourceId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Min(scan_at)",
			input: &RelAzureDataSourceForUpsert{
				AzureDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: unixtime19691231T235959,
			},
			wantErr: true,
		},
		{
			name: "NG Max(scan_at)",
			input: &RelAzureDataSourceForUpsert{
				AzureDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: unixtime100000101T000000,
			},
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
