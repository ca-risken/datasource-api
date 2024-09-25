package azure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/azure"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func convertAzureDataSource(data *model.AzureDataSource) *azure.AzureDataSource {
	if data == nil {
		return &azure.AzureDataSource{}
	}
	return &azure.AzureDataSource{
		AzureDataSourceId: data.AzureDataSourceID,
		Name:              data.Name,
		Description:       data.Description,
		MaxScore:          data.MaxScore,
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.UpdatedAt.Unix(),
	}
}

func (a *AzureService) ListAzureDataSource(ctx context.Context, req *azure.ListAzureDataSourceRequest) (*azure.ListAzureDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListAzureDataSource(ctx, req.AzureDataSourceId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &azure.ListAzureDataSourceResponse{}, nil
		}
		return nil, err
	}
	data := azure.ListAzureDataSourceResponse{}
	for _, d := range *list {
		data.AzureDataSource = append(data.AzureDataSource, convertAzureDataSource(&d))
	}
	return &data, nil
}

func convertAzure(data *model.Azure) *azure.Azure {
	if data == nil {
		return &azure.Azure{}
	}
	azure := azure.Azure{
		AzureId:          data.AzureID,
		Name:             data.Name,
		ProjectId:        data.ProjectID,
		SubscriptionId:   data.SubscriptionID,
		VerificationCode: data.VerificationCode,
		CreatedAt:        data.CreatedAt.Unix(),
		UpdatedAt:        data.UpdatedAt.Unix(),
	}
	return &azure
}

func (a *AzureService) ListAzure(ctx context.Context, req *azure.ListAzureRequest) (*azure.ListAzureResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListAzure(ctx, req.ProjectId, req.AzureId, req.SubscriptionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &azure.ListAzureResponse{}, nil
		}
		return nil, err
	}
	data := azure.ListAzureResponse{}
	for _, d := range *list {
		data.Azure = append(data.Azure, convertAzure(&d))
	}
	return &data, nil
}

func (a *AzureService) GetAzure(ctx context.Context, req *azure.GetAzureRequest) (*azure.GetAzureResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := a.repository.GetAzure(ctx, req.ProjectId, req.AzureId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &azure.GetAzureResponse{}, nil
		}
		return nil, err
	}
	return &azure.GetAzureResponse{Azure: convertAzure(data)}, nil
}

func (a *AzureService) PutAzure(ctx context.Context, req *azure.PutAzureRequest) (*azure.PutAzureResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registerd, err := a.repository.UpsertAzure(ctx, req.Azure)
	if err != nil {
		return nil, err
	}
	return &azure.PutAzureResponse{Azure: convertAzure(registerd)}, nil
}

func (a *AzureService) DeleteAzure(ctx context.Context, req *azure.DeleteAzureRequest) (*azure.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListRelAzureDataSource(ctx, req.ProjectId, req.AzureId)
	if err != nil {
		return nil, err
	}
	for _, ds := range *list {
		if err := a.repository.DeleteRelAzureDataSource(ctx, req.ProjectId, req.AzureId, ds.AzureDataSourceID); err != nil {
			return nil, err
		}
	}
	if err := a.repository.DeleteAzure(ctx, req.ProjectId, req.AzureId); err != nil {
		return nil, err
	}
	return &azure.Empty{}, nil
}

func (a *AzureService) ListRelAzureDataSource(ctx context.Context, req *azure.ListRelAzureDataSourceRequest) (*azure.ListRelAzureDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListRelAzureDataSource(ctx, req.ProjectId, req.AzureId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &azure.ListRelAzureDataSourceResponse{}, nil
		}
		return nil, err
	}
	data := azure.ListRelAzureDataSourceResponse{}
	for _, d := range *list {
		data.RelAzureDataSource = append(data.RelAzureDataSource, convertRelAzureDataSource(&d))
	}
	return &data, nil
}

func convertRelAzureDataSource(data *db.RelAzureDataSource) *azure.RelAzureDataSource {
	if data == nil {
		return &azure.RelAzureDataSource{}
	}
	azure := azure.RelAzureDataSource{
		AzureId:           data.AzureID,
		AzureDataSourceId: data.AzureDataSourceID,
		ProjectId:         data.ProjectID,
		Status:            getStatus(data.Status),
		StatusDetail:      data.StatusDetail,
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.UpdatedAt.Unix(),
		Name:              data.Name,           // azure_data_source.name
		MaxScore:          data.MaxScore,       // azure_data_source.max_score
		Description:       data.Description,    // azure_data_source.description
		SubscriptionId:    data.SubscriptionID, // azure.subscription_id
	}
	if !zero.IsZeroVal(data.ScanAt) {
		azure.ScanAt = data.ScanAt.Unix()
	}
	return &azure
}

func getStatus(s string) azure.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := azure.Status_value[statusKey]; !ok {
		return azure.Status_UNKNOWN
	}
	switch statusKey {
	case azure.Status_OK.String():
		return azure.Status_OK
	case azure.Status_CONFIGURED.String():
		return azure.Status_CONFIGURED
	case azure.Status_IN_PROGRESS.String():
		return azure.Status_IN_PROGRESS
	case azure.Status_ERROR.String():
		return azure.Status_ERROR
	default:
		return azure.Status_UNKNOWN
	}
}

func (a *AzureService) GetRelAzureDataSource(ctx context.Context, req *azure.GetRelAzureDataSourceRequest) (*azure.GetRelAzureDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := a.repository.GetRelAzureDataSource(ctx, req.ProjectId, req.AzureId, req.AzureDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &azure.GetRelAzureDataSourceResponse{}, nil
		}
		return nil, err
	}
	return &azure.GetRelAzureDataSourceResponse{RelAzureDataSource: convertRelAzureDataSource(data)}, nil
}

func (a *AzureService) AttachRelAzureDataSource(ctx context.Context, req *azure.AttachRelAzureDataSourceRequest) (*azure.AttachRelAzureDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	az, err := a.repository.GetAzure(ctx, req.ProjectId, req.RelAzureDataSource.AzureId)
	if err != nil {
		return nil, err
	}
	if ok, err := a.azureClient.VerifyCode(ctx, az.SubscriptionID, az.VerificationCode); !ok || err != nil {
		return nil, err
	}
	registered, err := a.repository.UpsertRelAzureDataSource(ctx, req.RelAzureDataSource)
	if err != nil {
		return nil, err
	}
	if !registered.ErrorNotifiedAt.IsZero() && registered.Status != azure.Status_ERROR.String() {
		if err := a.repository.UpdateAzureErrorNotifiedAt(ctx, gorm.Expr("NULL"), registered.AzureID, registered.AzureDataSourceID, registered.ProjectID); err != nil {
			return nil, err
		}
	}
	return &azure.AttachRelAzureDataSourceResponse{RelAzureDataSource: convertRelAzureDataSource(registered)}, nil
}

func (a *AzureService) DetachRelAzureDataSource(ctx context.Context, req *azure.DetachRelAzureDataSourceRequest) (*azure.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := a.repository.DeleteRelAzureDataSource(ctx, req.ProjectId, req.AzureId, req.AzureDataSourceId)
	if err != nil {
		return nil, err
	}
	return &azure.Empty{}, nil
}

const (
	prowlerDataSourceID uint32 = 1001
)

func (a *AzureService) InvokeScanAzure(ctx context.Context, req *azure.InvokeScanAzureRequest) (*azure.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	az, err := a.repository.GetAzure(ctx, req.ProjectId, req.AzureId)
	if err != nil {
		return nil, err
	}
	data, err := a.repository.GetRelAzureDataSource(ctx, req.ProjectId, req.AzureId, req.AzureDataSourceId)
	if err != nil {
		return nil, err
	}
	if ok, err := a.azureClient.VerifyCode(ctx, az.SubscriptionID, az.VerificationCode); !ok || err != nil {
		if _, upErr := a.repository.UpsertRelAzureDataSource(ctx, &azure.RelAzureDataSourceForUpsert{
			AzureId:           data.AzureID,
			AzureDataSourceId: data.AzureDataSourceID,
			ProjectId:         data.ProjectID,
			Status:            azure.Status_ERROR,
			StatusDetail:      err.Error(),
			ScanAt:            data.ScanAt.Unix(),
		}); upErr != nil {
			return nil, upErr
		}
		a.logger.Warnf(ctx, "Failed to verify code: azure_id=%d, subscription_id=%s, err=%v", data.AzureID, data.SubscriptionID, err)
		return &azure.Empty{}, nil
	}
	msg := &message.AzureQueueMessage{
		AzureID:           data.AzureID,
		ProjectID:         data.ProjectID,
		AzureDataSourceID: data.AzureDataSourceID,
		ScanOnly:          req.ScanOnly,
	}
	var resp *sqs.SendMessageOutput
	switch data.AzureDataSourceID {
	case prowlerDataSourceID:
		resp, err = a.sqs.Send(ctx, a.sqs.AzureProwlerQueueURL, msg)
	default:
		return nil, fmt.Errorf("unknown azureDataSourceID: %d", data.AzureDataSourceID)
	}
	if err != nil {
		return nil, err
	}
	if _, err = a.repository.UpsertRelAzureDataSource(ctx, &azure.RelAzureDataSourceForUpsert{
		AzureId:           data.AzureID,
		AzureDataSourceId: data.AzureDataSourceID,
		ProjectId:         data.ProjectID,
		Status:            azure.Status_IN_PROGRESS,
		StatusDetail:      fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:            data.ScanAt.Unix(),
	}); err != nil {
		return nil, err
	}
	a.logger.Infof(ctx, "Invoke scanned: subscription_id=%s, messageId=%v", data.SubscriptionID, resp.MessageId)
	return &azure.Empty{}, nil
}

func (a *AzureService) InvokeScanAll(ctx context.Context, req *azure.InvokeScanAllRequest) (*azure.Empty, error) {
	list, err := a.repository.ListRelAzureDataSourceByDataSourceID(ctx, req.AzureDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &azure.Empty{}, nil
		}
		return nil, err
	}
	for _, az := range *list {
		if resp, err := a.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: az.ProjectID}); err != nil {
			a.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			a.logger.Infof(ctx, "Skip deactive project, project_id=%d", az.ProjectID)
			continue
		}

		if _, err := a.InvokeScanAzure(ctx, &azure.InvokeScanAzureRequest{
			AzureId:           az.AzureID,
			ProjectId:         az.ProjectID,
			AzureDataSourceId: az.AzureDataSourceID,
			ScanOnly:          true,
		}); err != nil {
			// In Azure, an error may occur during InvokeScan due to user misconfiguration(e.a. invalid verification_code).
			// But to avoid having a single error stop the entire process, a notification log is output and other processes continue.
			a.logger.Notifyf(ctx, logging.ErrorLevel, "InvokeScanAzure error occured: azure_id=%d, subscription_id=%s, err=%+v", az.AzureID, az.SubscriptionID, err)
			continue
		}
	}
	return &azure.Empty{}, nil
}
