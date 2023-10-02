package aws

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/attackflow"
	"github.com/ca-risken/datasource-api/proto/datasource"
)

func TestGetCpuMemLabel(t *testing.T) {
	type args struct {
		cpu string
		mem string
	}
	cases := []struct {
		name  string
		input args
		want  string
	}{
		{
			name:  "OK 1",
			input: args{cpu: "250", mem: "500"},
			want:  "CPU: 0.25vCPU, MEM: 0.50GB",
		},
		{
			name:  "OK 2",
			input: args{cpu: "250000", mem: "500000"},
			want:  "CPU: 250.00vCPU, MEM: 500.00GB",
		},
		{
			name:  "Unknown CPU & MEM",
			input: args{cpu: "hoge", mem: "fuga"},
			want:  "CPU: hoge, MEM: fuga",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			analyzer := appRunnerAnalyzer{logger: logging.NewLogger()}
			got := analyzer.getCpuMemLabel(context.TODO(), c.input.cpu, c.input.mem)
			if got != c.want {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetSourceCodeNode(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *datasource.Resource
	}{
		{
			name:  "GitHub",
			input: "https://github.com/ca-risken/datasource-api",
			want: &datasource.Resource{
				ResourceName: "ca-risken/datasource-api",
				ShortName:    "ca-risken/datasource-api",
				Layer:        attackflow.LAYER_CODE_REPOSITORY,
				Region:       attackflow.REGION_GLOBAL,
				Service:      "github",
			},
		},
		{
			name:  "Other",
			input: "https://gitlab.com/owner/repo",
			want: &datasource.Resource{
				ResourceName: "https://gitlab.com/owner/repo",
				ShortName:    "https://gitlab.com/owner/repo",
				Layer:        attackflow.LAYER_CODE_REPOSITORY,
				Region:       attackflow.REGION_GLOBAL,
				Service:      "code-repository",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getSourceCodeNode(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetPublicEcrNode(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *datasource.Resource
		wantErr bool
	}{
		{
			name:  "OK",
			input: "public.ecr.aws/risken/risken-datasource-api:v0.8.0",
			want: &datasource.Resource{
				ResourceName: "arn:aws:ecr-public::risken:repository/risken-datasource-api",
				ShortName:    "risken-datasource-api",
				CloudType:    "aws",
				CloudId:      "risken",
				Layer:        attackflow.LAYER_DATASTORE,
				Region:       attackflow.REGION_GLOBAL,
				Service:      "ecr-public",
			},
		},
		{
			name:    "NG Invalid name",
			input:   "public.ecr.aws/risken",
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := getPublicEcrNode(c.input)
			if err != nil && !c.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && c.wantErr {
				t.Errorf("Expected error, but got nil")
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetPrivateEcrNode(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *datasource.Resource
		wantErr bool
	}{
		{
			name:  "OK",
			input: "123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/attack-flow-test:latest",
			want: &datasource.Resource{
				ResourceName: "arn:aws:ecr:ap-northeast-1:123456789012:repository/attack-flow-test",
				ShortName:    "attack-flow-test",
				CloudType:    "aws",
				CloudId:      "123456789012",
				Layer:        attackflow.LAYER_DATASTORE,
				Region:       "ap-northeast-1",
				Service:      "ecr",
			},
		},
		{
			name:    "NG Invalid name",
			input:   "123456789012dkrecrap-northeast-1amazonawscom/attack-flow-test:latest",
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := getPrivateEcrNode(c.input)
			if err != nil && !c.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && c.wantErr {
				t.Errorf("Expected error, but got nil")
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetAppRunnerAttackFlowCache(t *testing.T) {
	type args struct {
		cloudID      string
		resourceName string
		resource     *datasource.Resource
	}
	type output struct {
		resource *datasource.Resource
		meta     *appRunnerMetadata
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
				resourceName: "arn:aws:apprunner:ap-northeast-1:123456789012:service/service-name/xxx",
				resource: &datasource.Resource{
					CloudId:      "123456789012",
					ResourceName: "arn:aws:apprunner:ap-northeast-1:123456789012:service/service-name/xxx",
					MetaData:     `{"is_public":true}`,
				},
			},
			want: output{
				resource: &datasource.Resource{
					CloudId:      "123456789012",
					ResourceName: "arn:aws:apprunner:ap-northeast-1:123456789012:service/service-name/xxx",
					MetaData:     `{"is_public":true}`,
				},
				meta: &appRunnerMetadata{
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
			gotResource, gotMeta, err := getAppRunnerAttackFlowCache(c.input.cloudID, c.input.resourceName)
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
