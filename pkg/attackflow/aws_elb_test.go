package attackflow

import "testing"

func TestIsV2LoadBalancer(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "classic load balancer",
			input: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/lb-name",
			want:  false,
		},
		{
			name:  "application load balancer",
			input: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/app/lb-name/xxxxxxxxx",
			want:  true,
		},
		{
			name:  "invalid arn",
			input: "arnawselasticloadbalancingap-northeast-1123456789012:loadbalancer/app/lb-name/xxxxxxxxx",
			want:  false,
		},
		{
			name:  "unknown lb format",
			input: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/xxx/yyy",
			want:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := isV2LoadBalancer(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetV2LoadBalancerName(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "application load balancer",
			input: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/app/lb-name/xxxxxxxxx",
			want:  "lb-name",
		},
		{
			name:  "invalid arn",
			input: "arnawselasticloadbalancingap-northeast-1123456789012:loadbalancer/app/lb-name/xxxxxxxxx",
			want:  "",
		},
		{
			name:  "unknown lb format",
			input: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/xxx",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getV2LoadBalancerName(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
