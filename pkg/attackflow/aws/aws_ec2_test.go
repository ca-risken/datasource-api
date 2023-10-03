package aws

import (
	"reflect"
	"testing"

	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

func TestGetEC2AttackFlowCache(t *testing.T) {
	type args struct {
		cloudID      string
		resourceName string
		resource     *datasource.Resource
	}
	type output struct {
		resource *datasource.Resource
		meta     *ec2Metadata
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
				resourceName: "arn:aws:ec2:us-east-1:123456789012:instance/i-xxxxxxx",
				resource: &datasource.Resource{
					CloudId:      "123456789012",
					ResourceName: "arn:aws:ec2:us-east-1:123456789012:instance/i-xxxxxxx",
					MetaData:     `{"is_public":true}`,
				},
			},
			want: output{
				resource: &datasource.Resource{
					CloudId:      "123456789012",
					ResourceName: "arn:aws:ec2:us-east-1:123456789012:instance/i-xxxxxxx",
					MetaData:     `{"is_public":true}`,
				},
				meta: &ec2Metadata{
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
			gotResource, gotMeta, err := getEC2AttackFlowCache(c.input.cloudID, c.input.resourceName)
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
