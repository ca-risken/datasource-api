package gcp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	verificationLabelKey       = "risken"
	verificationErrMsgTemplate = "[Invalid code] Please check your GCP project label(key=%s), And then the registered verification_code must be the same value.(verification_code=%s)"
	projectGetRetryInterval    = 65 * time.Second
	projectGetMaxRetries       = 10
)

func (g *GcpClient) VerifyCode(ctx context.Context, gcpProjectID, verificationCode string) (bool, error) {
	if verificationCode == "" {
		return true, nil
	}
	// https://cloud.google.com/resource-manager/reference/rest/v1/projects/get
	cspan, cctx := tracer.StartSpanFromContext(ctx, "GetProject")
	resp, err := retryGetProject(cctx, projectGetRetryInterval, func() (*cloudresourcemanager.Project, error) {
		return g.crm.Projects.Get(gcpProjectID).Context(cctx).Do()
	}, g.newRetryLogger(ctx, "ResourceManager.Projects.Get"))
	cspan.Finish(tracer.WithError(err))
	if err != nil {
		g.logger.Warnf(ctx, "Failed to ResourceManager.Projects.Get API, err=%+v", err)
		return false, fmt.Errorf("failed to ResourceManager.Projects.Get API, err=%+v", err)
	}
	if v, ok := resp.Labels[verificationLabelKey]; !ok || v != verificationCode {
		return false, fmt.Errorf(verificationErrMsgTemplate, verificationLabelKey, verificationCode)
	}
	return true, nil
}

func retryGetProject(ctx context.Context, interval time.Duration, get func() (*cloudresourcemanager.Project, error), notify backoff.Notify) (*cloudresourcemanager.Project, error) {
	operation := func() (*cloudresourcemanager.Project, error) {
		project, err := get()
		if err != nil {
			var apiErr *googleapi.Error
			if !errors.As(err, &apiErr) || apiErr.Code != 429 {
				return nil, backoff.Permanent(err)
			}
		}
		return project, err
	}
	retryer := backoff.WithContext(backoff.WithMaxRetries(backoff.NewConstantBackOff(interval), projectGetMaxRetries), ctx)
	return backoff.RetryNotifyWithData(operation, retryer, notify)
}
