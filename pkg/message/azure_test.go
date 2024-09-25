package message

import (
	"reflect"
	"testing"
)

func TestValidateAzure(t *testing.T) {
	cases := []struct {
		name    string
		input   *AzureQueueMessage
		wantErr bool
	}{
		{
			name:  "OK (prowler)",
			input: &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, ProjectID: 1},
		},
		{
			name:    "NG Required(AzureID)",
			input:   &AzureQueueMessage{AzureID: 0, AzureDataSourceID: 2, ProjectID: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(AzureDataSourceID)",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 0, ProjectID: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectID)",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, ProjectID: 0},
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

func TestParseMessageAzure(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *AzureQueueMessage
		wantErr bool
	}{
		{
			name:  "OK",
			input: `{"azure_id":1, "azure_data_source_id":1, "project_id":1}`,
			want:  &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1},
		},
		{
			name:  "OK(scan_only)",
			input: `{"azure_id":1, "azure_data_source_id":1, "project_id":1, "scan_only":"true"}`,
			want:  &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, ScanOnly: true},
		},
		{
			name:    "NG Json parse error",
			input:   `{"parse...: error`,
			wantErr: true,
		},
		{
			name:    "NG Invalid message(required parammeter)",
			input:   `{}`,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseMessageAzure(c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error occured, wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpaeted response, want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
