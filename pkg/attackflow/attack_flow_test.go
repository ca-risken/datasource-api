package attackflow

import (
	"reflect"
	"testing"

	"github.com/ca-risken/datasource-api/proto/datasource"
)

func TestGetEdge(t *testing.T) {
	type args struct {
		source, target, edgeLabel string
	}
	cases := []struct {
		name  string
		input *args
		want  *datasource.ResourceRelationship
	}{
		{
			name: "OK",
			input: &args{
				source:    "s",
				target:    "t",
				edgeLabel: "ec2",
			},
			want: &datasource.ResourceRelationship{
				RelationId:         "ed-[s]-[t]",
				SourceResourceName: "s",
				TargetResourceName: "t",
				RelationLabel:      "ec2",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := GetEdge(c.input.source, c.input.target, c.input.edgeLabel)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestParseMetadata(t *testing.T) {
	type args struct {
		Test string `json:"test"`
	}
	cases := []struct {
		name  string
		input *args
		want  string
	}{
		{
			name: "OK",
			input: &args{
				Test: "test",
			},
			want: `{"test":"test"}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := ParseMetadata(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestSetNode(t *testing.T) {
	type args struct {
		isPublic          bool
		internetEdgeLabel string
		resource          *datasource.Resource
		resp              *datasource.AnalyzeAttackFlowResponse
	}
	cases := []struct {
		name  string
		input *args
		want  *datasource.AnalyzeAttackFlowResponse
	}{
		{
			name: "OK(not public)",
			input: &args{
				isPublic: false,
				resource: &datasource.Resource{
					ResourceName: "test",
				},
				resp: &datasource.AnalyzeAttackFlowResponse{},
			},
			want: &datasource.AnalyzeAttackFlowResponse{
				Nodes: []*datasource.Resource{
					{
						ResourceName: "test",
					},
				},
			},
		},
		{
			name: "OK(public)",
			input: &args{
				isPublic:          true,
				internetEdgeLabel: "label",
				resource: &datasource.Resource{
					ResourceName: "test",
				},
				resp: &datasource.AnalyzeAttackFlowResponse{},
			},
			want: &datasource.AnalyzeAttackFlowResponse{
				Nodes: []*datasource.Resource{
					{
						ResourceName: RESOURCE_INTERNET,
						ShortName:    RESOURCE_INTERNET,
						Layer:        LAYER_INTERNET,
						Region:       REGION_GLOBAL,
						Service:      "internet",
					},
					{
						ResourceName: "test",
					},
				},
				Edges: []*datasource.ResourceRelationship{
					{
						RelationId:         "ed-[Internet]-[test]",
						RelationLabel:      "label",
						SourceResourceName: RESOURCE_INTERNET,
						TargetResourceName: "test",
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SetNode(c.input.isPublic, c.input.internetEdgeLabel, c.input.resource, c.input.resp)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetExternalServiceNode(t *testing.T) {
	testCases := []struct {
		name             string
		target           string
		expectedResource *datasource.Resource
	}{
		{
			name:   "example.com",
			target: "example.com",
			expectedResource: &datasource.Resource{
				ResourceName: "example.com",
				ShortName:    "example.com",
				Layer:        LAYER_EXTERNAL_SERVICE,
				Region:       REGION_GLOBAL,
				Service:      "external-service",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource := GetExternalServiceNode(tc.target)

			if !reflect.DeepEqual(resource, tc.expectedResource) {
				t.Errorf("Expected resource %+v, but got %+v", tc.expectedResource, resource)
			}
		})
	}
}

func TestGetInternalServiceNode(t *testing.T) {
	testCases := []struct {
		name             string
		target           string
		region           string
		expectedResource *datasource.Resource
	}{
		{
			name:   "us-west-2",
			target: "internal-service",
			region: "us-west-2",
			expectedResource: &datasource.Resource{
				ResourceName: "internal-service",
				ShortName:    "internal-service",
				Layer:        LAYER_INTERNAL_SERVICE,
				Region:       "us-west-2",
				Service:      "internal-service",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource := GetInternalServiceNode(tc.target, tc.region)
			if !reflect.DeepEqual(resource, tc.expectedResource) {
				t.Errorf("Expected resource %+v, but got %+v", tc.expectedResource, resource)
			}
		})
	}
}

func TestGetCodeRepositoryNode(t *testing.T) {
	testCases := []struct {
		name             string
		repository       string
		service          string
		expectedResource *datasource.Resource
	}{
		{
			name:       "GitHub",
			repository: "github.com/example/repo",
			service:    "GitHub",
			expectedResource: &datasource.Resource{
				ResourceName: "github.com/example/repo",
				ShortName:    "github.com/example/repo",
				Layer:        LAYER_CODE_REPOSITORY,
				Region:       REGION_GLOBAL,
				Service:      "GitHub",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource := GetCodeRepositoryNode(tc.repository, tc.service)
			if !reflect.DeepEqual(resource, tc.expectedResource) {
				t.Errorf("Expected resource %+v, but got %+v", tc.expectedResource, resource)
			}
		})
	}
}
