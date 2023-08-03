package attackflow

import "testing"

func TestGetEventBridgeResourceType(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "event bus",
			input: "arn:aws:events:ap-northeast-1:123456789012:event-bus/resource-name",
			want:  "event-bus",
		},
		{
			name:  "api destination",
			input: "arn:aws:events:ap-northeast-1:123456789012:api-destination/resource-name/uuid",
			want:  "api-destination",
		},
		{
			name:  "invalid arn",
			input: "arn:aws:events:ap-northeast-1:123456789012:invalid",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getEventBridgeResourceType(c.input)
			if got != c.want {
				t.Errorf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
