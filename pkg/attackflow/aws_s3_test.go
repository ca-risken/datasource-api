package attackflow

import (
	"testing"
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
			name:  "s3(inclued dot)",
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
