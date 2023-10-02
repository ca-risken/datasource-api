package attackflow

import (
	"reflect"
	"testing"

	"github.com/ca-risken/datasource-api/proto/datasource"
)

func TestSetAndGetAttackFlowCache(t *testing.T) {
	type args struct {
		cloudID      string
		resourceName string
		data         *datasource.Resource
	}
	testCases := []struct {
		name  string
		input *args
		want  *datasource.Resource
	}{
		{
			name: "valid data",
			input: &args{
				cloudID:      "aws",
				resourceName: "resource-1",
				data: &datasource.Resource{
					ResourceName: "resource-1",
					ShortName:    "resource-1",
					Layer:        LAYER_COMPUTE,
					Region:       "us-west-2",
					Service:      "EC2",
				},
			},
			want: &datasource.Resource{
				ResourceName: "resource-1",
				ShortName:    "resource-1",
				Layer:        LAYER_COMPUTE,
				Region:       "us-west-2",
				Service:      "EC2",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetAttackFlowCache(tc.input.cloudID, tc.input.resourceName, tc.input.data)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			cachedData, err := GetAttackFlowCache(tc.input.cloudID, tc.input.resourceName)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(cachedData, tc.input.data) {
				t.Errorf("Expected data %+v, but got %+v", tc.input.data, cachedData)
			}
		})
	}
}
