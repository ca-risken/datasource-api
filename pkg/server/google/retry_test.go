package google

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/api/googleapi"
)

func TestRetryVerifyCode(t *testing.T) {
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
			ok, err := retryVerifyCode(context.Background(), 0, func() (bool, error) {
				attempts++
				if attempts <= tt.failures {
					return false, tt.err
				}
				return true, nil
			}, nil)
			if attempts != tt.wantAttempts {
				t.Fatalf("attempts = %d, want %d", attempts, tt.wantAttempts)
			}
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError %t", err, tt.wantError)
			}
			if !tt.wantError && !ok {
				t.Fatal("verification failed unexpectedly")
			}
		})
	}
}

func TestRetryVerifyCodeCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0
	_, err := retryVerifyCode(ctx, time.Hour, func() (bool, error) {
		attempts++
		cancel()
		return false, &googleapi.Error{Code: 429}
	}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
