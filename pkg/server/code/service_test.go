package code

import (
	"strings"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/github"
	"github.com/ca-risken/datasource-api/pkg/queue"
)

func TestNewCodeService(t *testing.T) {
	cases := []struct {
		name      string
		dataKey   string
		appAuth   *github.AppAuthConfig
		wantError string
	}{
		{
			name:    "OK without github app auth",
			dataKey: "12345678901234567890123456789012",
		},
		{
			name:      "NG invalid data key",
			dataKey:   "short",
			wantError: "failed to create cipher",
		},
		{
			name:      "NG invalid github app auth",
			dataKey:   "12345678901234567890123456789012",
			appAuth:   &github.AppAuthConfig{AppID: "12345"},
			wantError: "failed to create github client",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := NewCodeService(c.dataKey, c.appAuth, nil, &queue.Client{}, nil, 0, logging.NewLogger())
			if c.wantError != "" {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				if !strings.Contains(err.Error(), c.wantError) {
					t.Fatalf("Expected error to contain %q, got %v", c.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("Expected service but got nil")
			}
		})
	}
}
