package google

import (
	"context"
	"fmt"
	"os"

	"github.com/ca-risken/common/pkg/logging"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ResourceManagerServiceClient interface {
	verifyCode(ctx context.Context, gcpProjectID, verificationCode string) (bool, error)
}

type ResourceManagerClient struct {
	logger logging.Logger
	svc    *cloudresourcemanager.Service
}

func newResourceManagerClient(ctx context.Context, credentialPath string, logger logging.Logger) (ResourceManagerServiceClient, error) {
	svc, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(credentialPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create new Cloud Resource Manager service: err=%w", err)
	}

	// Remove credential file for Security
	if err := os.Remove(credentialPath); err != nil {
		return nil, fmt.Errorf("failed to remove file: path=%s, err=%w", credentialPath, err)
	}
	return &ResourceManagerClient{
		svc:    svc,
		logger: logger,
	}, nil
}

const (
	verificationLabelKey       = "risken"
	verificationErrMsgTemplate = "[Invalid code] Please check your GCP project label(key=%s), And then the registered verification_code must be the same value.(verification_code=%s)"
)

func (r *ResourceManagerClient) verifyCode(ctx context.Context, gcpProjectID, verificationCode string) (bool, error) {
	if verificationCode == "" {
		return true, nil
	}
	// https://cloud.google.com/resource-manager/reference/rest/v1/projects/get
	cspan, cctx := tracer.StartSpanFromContext(ctx, "GetProject")
	resp, err := r.svc.Projects.Get(gcpProjectID).Context(cctx).Do()
	cspan.Finish(tracer.WithError(err))
	if err != nil {
		r.logger.Warnf(ctx, "Failed to ResourceManager.Projects.Get API, err=%+v", err)
		return false, fmt.Errorf("Failed to ResourceManager.Projects.Get API, err=%+v", err)
	}
	r.logger.Debugf(ctx, "Got the project info: %+v", resp)
	if v, ok := resp.Labels[verificationLabelKey]; !ok || v != verificationCode {
		return false, fmt.Errorf(verificationErrMsgTemplate, verificationLabelKey, verificationCode)
	}
	return true, nil
}
