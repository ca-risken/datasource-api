package attackflow

import (
	"reflect"
	"testing"

	"github.com/ca-risken/datasource-api/proto/datasource"
)

func TestGetAWSInfoFromARN(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *datasource.Resource
	}{
		{
			name:  "OK(iam)",
			input: "arn:aws:iam::123456789012:user/MyUser",
			want: &datasource.Resource{
				ResourceName: "arn:aws:iam::123456789012:user/MyUser",
				ShortName:    "MyUser",
				CloudType:    "aws",
				CloudId:      "123456789012",
				Service:      "iam",
				Region:       "global",
				Layer:        LAYER_LATERAL_MOVEMENT,
			},
		},
		{
			name:  "OK(ec2)",
			input: "arn:aws:ec2:us-east-1:123456789012:instance/i-xxxxxxx",
			want: &datasource.Resource{
				ResourceName: "arn:aws:ec2:us-east-1:123456789012:instance/i-xxxxxxx",
				ShortName:    "i-xxxxxxx",
				CloudType:    "aws",
				CloudId:      "123456789012",
				Service:      "ec2",
				Region:       "us-east-1",
				Layer:        "",
			},
		},
		{
			name:  "OK(s3)",
			input: "arn:aws:s3:::bucket_name",
			want: &datasource.Resource{
				ResourceName: "arn:aws:s3:::bucket_name",
				ShortName:    "bucket_name",
				CloudType:    "aws",
				CloudId:      "",
				Service:      "s3",
				Region:       "global",
				Layer:        LAYER_DATASTORE,
			},
		},
		{
			name:  "OK(lambda)",
			input: "arn:aws:lambda:ap-northeast-1:123456789012:function:attack-flow-test",
			want: &datasource.Resource{
				ResourceName: "arn:aws:lambda:ap-northeast-1:123456789012:function:attack-flow-test",
				ShortName:    "attack-flow-test",
				CloudType:    "aws",
				CloudId:      "123456789012",
				Service:      "lambda",
				Region:       "ap-northeast-1",
				Layer:        LAYER_COMPUTE,
			},
		},
		{
			name:  "Blank",
			input: "",
			want:  nil,
		},
		{
			name:  "Invalid arn",
			input: "arnaws:iam123456789012:user/MyUser",
			want:  nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getAWSInfoFromARN(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestIsSupportedAWSService(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "OK",
			input: "cloudfront",
			want:  true,
		},
		{
			name:  "Blank",
			input: "",
			want:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := isSupportedAWSService(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestFindAWSServiceFromDomain(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "cloudfront",
			input: "distribution-id.cloudfront.net",
			want:  "cloudfront",
		},
		{
			name:  "s3",
			input: "bucket-name.s3.ap-northeast-1.amazonaws.com",
			want:  "s3",
		},
		{
			name:  "Unknown",
			input: "xxxxxx.example.com",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := findAWSServiceFromDomain(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
