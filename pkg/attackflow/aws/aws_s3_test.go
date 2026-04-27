package aws

import (
	"reflect"
	"testing"

	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

func TestGetS3ARNFromDomain(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "s3",
			input: "bucket-name.s3.ap-northeast-1.amazonaws.com",
			want:  "arn:aws:s3:::bucket-name",
		},
		{
			name:  "s3(include dot)",
			input: "bucket-name.co.jp.s3.ap-northeast-1.amazonaws.com",
			want:  "arn:aws:s3:::bucket-name.co.jp",
		},
		{
			name:  "Unknown",
			input: "xxxxxx.example.com",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getS3ARNFromDomain(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetS3AttackFlowCache(t *testing.T) {
	type args struct {
		cloudID      string
		resourceName string
		resource     *datasource.Resource
	}
	type output struct {
		resource *datasource.Resource
		meta     *S3Metadata
	}
	cases := []struct {
		name    string
		input   args
		want    output
		wantErr bool
	}{
		{
			name: "OK",
			input: args{
				cloudID:      "123456789012",
				resourceName: "arn:aws:s3:::bucket-name",
				resource: &datasource.Resource{
					CloudId:      "123456789012",
					ResourceName: "arn:aws:s3:::bucket-name",
					MetaData:     `{"is_public":true}`,
				},
			},
			want: output{
				resource: &datasource.Resource{
					CloudId:      "123456789012",
					ResourceName: "arn:aws:s3:::bucket-name",
					MetaData:     `{"is_public":true}`,
				},
				meta: &S3Metadata{
					IsPublic: true,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// set cache
			_ = attackflow.SetAttackFlowCache(c.input.cloudID, c.input.resourceName, c.input.resource)
			// get cache
			gotResource, gotMeta, err := getS3AttackFlowCache(c.input.cloudID, c.input.resourceName)
			if err != nil && !c.wantErr {
				t.Errorf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Errorf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(gotResource, c.want.resource) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want.resource, gotResource)
			}
			if !reflect.DeepEqual(gotMeta, c.want.meta) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want.meta, gotMeta)
			}
		})
	}
}
