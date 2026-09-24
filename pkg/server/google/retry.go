package google

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/api/googleapi"
)

const (
	projectGetRetryInterval = 65 * time.Second
	projectGetMaxRetries    = 10
)

func retryVerifyCode(ctx context.Context, interval time.Duration, verify func() (bool, error), notify backoff.Notify) (bool, error) {
	operation := func() (bool, error) {
		ok, err := verify()
		if err != nil {
			var apiErr *googleapi.Error
			if !errors.As(err, &apiErr) || apiErr.Code != 429 {
				return false, backoff.Permanent(err)
			}
		}
		return ok, err
	}
	retryer := backoff.WithContext(backoff.WithMaxRetries(backoff.NewConstantBackOff(interval), projectGetMaxRetries), ctx)
	return backoff.RetryNotifyWithData(operation, retryer, notify)
}
