package osint

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/aws"
	"github.com/ca-risken/datasource-api/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (o *OsintService) ListRelOsintDataSource(ctx context.Context, req *osint.ListRelOsintDataSourceRequest) (*osint.ListRelOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := o.repository.ListRelOsintDataSource(ctx, req.ProjectId, req.OsintId, req.OsintDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.ListRelOsintDataSourceResponse{}, nil
		}
		o.logger.Errorf(ctx, "Failed to List RelOsintDataSource. error: %v", err)
		return nil, err
	}
	data := osint.ListRelOsintDataSourceResponse{}
	for _, d := range *list {
		data.RelOsintDataSource = append(data.RelOsintDataSource, convertRelOsintDataSource(&d))
	}
	return &data, nil
}

func (o *OsintService) GetRelOsintDataSource(ctx context.Context, req *osint.GetRelOsintDataSourceRequest) (*osint.GetRelOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := o.repository.GetRelOsintDataSource(ctx, req.ProjectId, req.RelOsintDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.GetRelOsintDataSourceResponse{}, nil
		}
		o.logger.Errorf(ctx, "Failed to Get RelOsintDataSource. error: %v", err)
		return nil, err
	}

	return &osint.GetRelOsintDataSourceResponse{RelOsintDataSource: convertRelOsintDataSource(getData)}, nil
}

func (o *OsintService) PutRelOsintDataSource(ctx context.Context, req *osint.PutRelOsintDataSourceRequest) (*osint.PutRelOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.RelOsintDataSource{
		RelOsintDataSourceID: req.RelOsintDataSource.RelOsintDataSourceId,
		ProjectID:            req.ProjectId,
		OsintDataSourceID:    req.RelOsintDataSource.OsintDataSourceId,
		OsintID:              req.RelOsintDataSource.OsintId,
		Status:               req.RelOsintDataSource.Status.String(),
		StatusDetail:         req.RelOsintDataSource.StatusDetail,
		ScanAt:               time.Unix(req.RelOsintDataSource.ScanAt, 0),
	}

	registeredData, err := o.repository.UpsertRelOsintDataSource(ctx, data)
	if err != nil {
		o.logger.Errorf(ctx, "Failed to Put RelOsintDataSource. error: %v", err)
		return nil, err
	}
	if registeredData.ErrorNotifiedAt != nil &&
		!registeredData.ErrorNotifiedAt.IsZero() &&
		registeredData.Status != aws.Status_ERROR.String() {
		if err := o.repository.UpdateOsintErrorNotifiedAt(ctx, gorm.Expr("NULL"), registeredData.RelOsintDataSourceID, registeredData.ProjectID); err != nil {
			return nil, err
		}
	}
	return &osint.PutRelOsintDataSourceResponse{RelOsintDataSource: convertRelOsintDataSource(registeredData)}, nil
}

func (o *OsintService) DeleteRelOsintDataSource(ctx context.Context, req *osint.DeleteRelOsintDataSourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := o.deleteRelOsintDataSourceWithDetectWord(ctx, req.ProjectId, req.RelOsintDataSourceId); err != nil {
		o.logger.Errorf(ctx, "Failed to DeleteRelOsintDataSource. error: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (o *OsintService) deleteRelOsintDataSourceWithDetectWord(ctx context.Context, projectID, relOsintDataSourceID uint32) error {

	detectWords, err := o.repository.ListOsintDetectWord(ctx, projectID, relOsintDataSourceID)
	if err != nil {
		return err
	}

	for _, d := range *detectWords {
		if err := o.repository.DeleteOsintDetectWord(ctx, projectID, d.OsintDetectWordID); err != nil {
			return err
		}
	}

	if err := o.repository.DeleteRelOsintDataSource(ctx, projectID, relOsintDataSourceID); err != nil {
		return err
	}
	return nil
}

func convertRelOsintDataSource(data *model.RelOsintDataSource) *osint.RelOsintDataSource {
	if data == nil {
		return &osint.RelOsintDataSource{}
	}
	return &osint.RelOsintDataSource{
		RelOsintDataSourceId: data.RelOsintDataSourceID,
		OsintDataSourceId:    data.OsintDataSourceID,
		OsintId:              data.OsintID,
		ProjectId:            data.ProjectID,
		CreatedAt:            data.CreatedAt.Unix(),
		UpdatedAt:            data.CreatedAt.Unix(),
		Status:               getStatus(data.Status),
		StatusDetail:         data.StatusDetail,
		ScanAt:               data.ScanAt.Unix(),
	}
}

func (o *OsintService) InvokeScan(ctx context.Context, req *osint.InvokeScanRequest) (*osint.InvokeScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	relOsintDataSourceData, err := o.repository.GetRelOsintDataSource(ctx, req.ProjectId, req.RelOsintDataSourceId)
	if err != nil {
		return nil, err
	}
	osintDataSourceData, err := o.repository.GetOsintDataSource(ctx, relOsintDataSourceData.ProjectID, relOsintDataSourceData.OsintDataSourceID)
	if err != nil {
		return nil, err
	}
	osintData, err := o.repository.GetOsint(ctx, relOsintDataSourceData.ProjectID, relOsintDataSourceData.OsintID)
	if err != nil {
		return nil, err
	}
	detectWord, err := o.repository.ListOsintDetectWord(ctx, relOsintDataSourceData.ProjectID, relOsintDataSourceData.RelOsintDataSourceID)
	if err != nil {
		return nil, err
	}
	jsonDetectWord, err := json.Marshal(map[string][]model.OsintDetectWord{"DetectWord": *detectWord})
	if err != nil {
		return nil, err
	}
	msg := &message.OsintQueueMessage{
		DataSource:           osintDataSourceData.Name,
		RelOsintDataSourceID: req.RelOsintDataSourceId,
		OsintID:              relOsintDataSourceData.OsintID,
		OsintDataSourceID:    relOsintDataSourceData.OsintDataSourceID,
		ProjectID:            req.ProjectId,
		ResourceType:         osintData.ResourceType,
		ResourceName:         osintData.ResourceName,
		DetectWord:           string(jsonDetectWord),
		ScanOnly:             req.ScanOnly,
	}

	var resp *sqs.SendMessageOutput
	switch msg.DataSource {
	case message.SubdomainDataSource:
		resp, err = o.sqs.Send(ctx, o.sqs.OSINTSubdomainQueueURL, msg)
	case message.WebsiteDataSource:
		resp, err = o.sqs.Send(ctx, o.sqs.OSINTWebsiteQueueURL, msg)
	default:
		return nil, fmt.Errorf("unknown datasource, datasource=%s", msg.DataSource)
	}
	if err != nil {
		return nil, err
	}
	if _, err = o.repository.UpsertRelOsintDataSource(ctx, &model.RelOsintDataSource{
		RelOsintDataSourceID: relOsintDataSourceData.RelOsintDataSourceID,
		OsintID:              relOsintDataSourceData.OsintID,
		OsintDataSourceID:    relOsintDataSourceData.OsintDataSourceID,
		ProjectID:            relOsintDataSourceData.ProjectID,
		Status:               osint.Status_IN_PROGRESS.String(),
		StatusDetail:         fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:               relOsintDataSourceData.ScanAt,
	}); err != nil {
		o.logger.Errorf(ctx, "Failed to update scan status: %+v", err)
		return nil, err
	}
	o.logger.Infof(ctx, "Invoked scan. MessageId: %v", *resp.MessageId)
	return &osint.InvokeScanResponse{Message: "Invoke Scan."}, nil
}

func (o *OsintService) InvokeScanAll(ctx context.Context, req *osint.InvokeScanAllRequest) (*empty.Empty, error) {

	list, err := o.repository.ListAllRelOsintDataSource(ctx, req.OsintDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &empty.Empty{}, nil
		}
		o.logger.Errorf(ctx, "Failed to List AllRelOsintDataSource. error: %v", err)
		return nil, err
	}

	for _, relOsintDataSource := range *list {
		if resp, err := o.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: relOsintDataSource.ProjectID}); err != nil {
			o.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			o.logger.Infof(ctx, "Skip deactive project, project_id=%d", relOsintDataSource.ProjectID)
			continue
		}

		if _, err := o.InvokeScan(ctx, &osint.InvokeScanRequest{
			ProjectId:            relOsintDataSource.ProjectID,
			RelOsintDataSourceId: relOsintDataSource.RelOsintDataSourceID,
			ScanOnly:             true,
		}); err != nil {
			o.logger.Errorf(ctx, "InvokeScanAll error: project_id=%d, rel_osint_data_source_id=%d, err=%+v",
				relOsintDataSource.ProjectID, relOsintDataSource.RelOsintDataSourceID, err)
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

func getStatus(s string) osint.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := osint.Status_value[statusKey]; !ok {
		return osint.Status_UNKNOWN
	}
	switch statusKey {
	case osint.Status_OK.String():
		return osint.Status_OK
	case osint.Status_CONFIGURED.String():
		return osint.Status_CONFIGURED
	case osint.Status_IN_PROGRESS.String():
		return osint.Status_IN_PROGRESS
	case osint.Status_ERROR.String():
		return osint.Status_ERROR
	default:
		return osint.Status_UNKNOWN
	}
}
