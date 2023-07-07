package attackflow

import (
	"context"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
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
