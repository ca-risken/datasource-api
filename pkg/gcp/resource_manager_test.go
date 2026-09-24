package gcp

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

func TestRetryGetProject(t *testing.T) {
	tests := []struct {
		name         string
		failures     int
		err          error
		wantAttempts int
		wantError    bool
	}{
		{name: "retry 429 then succeed", failures: 2, err: &googleapi.Error{Code: 429}, wantAttempts: 3},
		{name: "stop after ten retries", failures: 11, err: &googleapi.Error{Code: 429}, wantAttempts: 11, wantError: true},
		{name: "do not retry other errors", failures: 1, err: &googleapi.Error{Code: 403}, wantAttempts: 1, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempts := 0
			project, err := retryGetProject(context.Background(), 0, func() (*cloudresourcemanager.Project, error) {
				attempts++
				if attempts <= tt.failures {
					return nil, tt.err
				}
				return &cloudresourcemanager.Project{ProjectId: "test-project"}, nil
			}, nil)
			if attempts != tt.wantAttempts {
				t.Fatalf("attempts = %d, want %d", attempts, tt.wantAttempts)
			}
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError %t", err, tt.wantError)
			}
			if !tt.wantError && (project == nil || project.ProjectId != "test-project") {
				t.Fatalf("project = %+v", project)
			}
		})
	}
}

func TestRetryGetProjectCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0
	_, err := retryGetProject(ctx, time.Hour, func() (*cloudresourcemanager.Project, error) {
		attempts++
		cancel()
		return nil, &googleapi.Error{Code: 429}
	}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
