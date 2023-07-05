package attackflow

import (
	"testing"
)

func TestExtractApiID(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "OK(RestAPI)",
			input: "arn:aws:apigateway:ap-northeast-1::/restapis/xxx",
			want:  "xxx",
		},
		{
			name:  "OK(HTTP API or WebSocket)",
			input: "arn:aws:apigateway:ap-northeast-1::/apis/xxx",
			want:  "xxx",
		},
		{
			name:    "Invalid arn",
			input:   "arnawsapigatewayap-northeast-1/restapis/xxx",
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := extractApiID(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%s, got=%s", c.want, got)
			}
			if !c.wantErr && err != nil {
				t.Errorf("Unexpected error: err=%+v", err)
			}
			if c.wantErr && err == nil {
				t.Errorf("No error: wantErr=%t", c.wantErr)
			}
		})
	}
}

func TestExtractArnFromMethodIntegration(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OK(Lambda)",
			input: "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:my-function/invocations",
			want:  "arn:aws:lambda:us-east-1:123456789012:function:my-function",
		},
		{
			name:  "OK(S3 action)",
			input: "arn:aws:apigateway:us-west-2:s3:action/GetObject&Bucket=bucket-name&Key=index.html",
			want:  "arn:aws:s3:::bucket-name",
		},
		{
			name:  "OK(S3 path)",
			input: "arn:aws:apigateway:us-west-2:s3:path/bucket-name/index.html",
			want:  "arn:aws:s3:::bucket-name",
		},
		{
			name:  "Invalid arn",
			input: "arn:aws:apigateway:us-east-1:hoge-service:path/hoge/fuga",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := extractArnFromMethodIntegration(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%s, got=%s", c.want, got)
			}
		})
	}
}
