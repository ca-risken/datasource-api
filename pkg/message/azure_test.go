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
			input: &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, DataSource: "azure:prowler", ProjectID: 1, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: ""},
		},
		{
			name:    "NG Required(AzureID)",
			input:   &AzureQueueMessage{AzureID: 0, AzureDataSourceID: 2, DataSource: "azure:prowler", ProjectID: 1, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: ""},
			wantErr: true,
		},
		{
			name:    "NG Required(AzureDataSourceID)",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 0, DataSource: "azure:prowler", ProjectID: 1, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: ""},
			wantErr: true,
		},
		{
			name:    "NG Required(DataSource)",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, DataSource: "", ProjectID: 1, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: ""},
			wantErr: true,
		},
		{
			name:    "NG Unknown(DataSource)",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, DataSource: "azure:unknown", ProjectID: 1, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: ""},
			wantErr: true,
		},
		{
			name:    "NG Required(ProjectID)",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, DataSource: "azure:prowler", ProjectID: 0, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: ""},
			wantErr: true,
		},
		{
			name:    "NG Required(SubscriptionID)",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, DataSource: "azure:prowler", ProjectID: 0, SubscriptionID: "", VerificationCode: ""},
			wantErr: true,
		},
		{
			name:    "NG Invalid Length(SubscriptionID) short",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, DataSource: "azure:prowler", ProjectID: 0, SubscriptionID: "12345678901234567890123456789012345", VerificationCode: ""},
			wantErr: true,
		},
		{
			name:    "NG Invalid Length(SubscriptionID) long",
			input:   &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 2, DataSource: "azure:prowler", ProjectID: 0, SubscriptionID: "1234567890123456789012345678901234567", VerificationCode: ""},
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
			input: `{"azure_id":1, "azure_data_source_id":1, "data_source":"azure:prowler", "project_id":1, "subscription_id":"123456789012345678901234567890123456",  "verification_code":""}`,
			want:  &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 1, DataSource: "azure:prowler", ProjectID: 1, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: ""},
		},
		{
			name:  "OK(scan_only)",
			input: `{"azure_id":1, "azure_data_source_id":1, "data_source":"azure:prowler", "project_id":1, "subscription_id":"123456789012345678901234567890123456",  "verification_code":"", "scan_only":"true"}`,
			want:  &AzureQueueMessage{AzureID: 1, AzureDataSourceID: 1, DataSource: "azure:prowler", ProjectID: 1, SubscriptionID: "123456789012345678901234567890123456", VerificationCode: "", ScanOnly: true},
		},
		{
			name:    "NG Json parse erroro",
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
