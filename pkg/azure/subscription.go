package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	scope                      = "/subscriptions/%s"
	verificationLabelKey       = "risken"
	verificationErrMsgTemplate = "[Invalid code] Please check your Azure subscription label(key=%s), And then the registered verification_code must be the same value.(verification_code=%s)"
)

func (a *AzureClient) VerifyCode(ctx context.Context, subscriptionID, verificationCode string) (bool, error) {
	if verificationCode == "" {
		return true, nil
	}
	// https://learn.microsoft.com/ja-jp/rest/api/resources/tags/get-at-scope
	cspan, cctx := tracer.StartSpanFromContext(ctx, "GetTags")
	tagsClient, err := armresources.NewTagsClient(subscriptionID, a.cred, nil)
	if err != nil {
		return false, err
	}
	resp, err := tagsClient.GetAtScope(cctx, fmt.Sprintf(scope, subscriptionID), nil)
	if err != nil {
		return false, err
	}
	cspan.Finish(tracer.WithError(err))
	if resp.TagsResource.Properties == nil {
		return false, fmt.Errorf("tagsResource.Properties is nil")
	}
	for k, v := range resp.TagsResource.Properties.Tags {
		if k == verificationLabelKey && *v == verificationCode {
			return true, nil
		}
	}
	cspan.Finish(tracer.WithError(err))

	return false, nil
}
