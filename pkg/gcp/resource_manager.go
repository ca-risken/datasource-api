package gcp

import (
	"context"
	"fmt"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	verificationLabelKey       = "risken"
	verificationErrMsgTemplate = "[Invalid code] Please check your GCP project label(key=%s), And then the registered verification_code must be the same value.(verification_code=%s)"
)

func (g *GcpClient) VerifyCode(ctx context.Context, gcpProjectID, verificationCode string) (bool, error) {
	if verificationCode == "" {
		return true, nil
	}
	// https://cloud.google.com/resource-manager/reference/rest/v1/projects/get
	cspan, cctx := tracer.StartSpanFromContext(ctx, "GetProject")
	resp, err := g.crm.Projects.Get(gcpProjectID).Context(cctx).Do()
	cspan.Finish(tracer.WithError(err))
	if err != nil {
		g.logger.Warnf(ctx, "Failed to ResourceManager.Projects.Get API, err=%+v", err)
		return false, fmt.Errorf("Failed to ResourceManager.Projects.Get API, err=%+v", err)
	}
	if v, ok := resp.Labels[verificationLabelKey]; !ok || v != verificationCode {
		return false, fmt.Errorf(verificationErrMsgTemplate, verificationLabelKey, verificationCode)
	}
	return true, nil
}
