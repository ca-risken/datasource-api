package google

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
	"github.com/ca-risken/datasource-api/proto/google"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func convertGoogleDataSource(data *model.GoogleDataSource) *google.GoogleDataSource {
	if data == nil {
		return &google.GoogleDataSource{}
	}
	return &google.GoogleDataSource{
		GoogleDataSourceId: data.GoogleDataSourceID,
		Name:               data.Name,
		Description:        data.Description,
		MaxScore:           data.MaxScore,
		CreatedAt:          data.CreatedAt.Unix(),
		UpdatedAt:          data.UpdatedAt.Unix(),
	}
}

func (g *GoogleService) ListGoogleDataSource(ctx context.Context, req *google.ListGoogleDataSourceRequest) (*google.ListGoogleDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := g.repository.ListGoogleDataSource(ctx, req.GoogleDataSourceId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &google.ListGoogleDataSourceResponse{}, nil
		}
		return nil, err
	}
	data := google.ListGoogleDataSourceResponse{}
	for _, d := range *list {
		data.GoogleDataSource = append(data.GoogleDataSource, convertGoogleDataSource(&d))
	}
	return &data, nil
}

func convertGCP(data *model.GCP) *google.GCP {
	if data == nil {
		return &google.GCP{}
	}
	gcp := google.GCP{
		GcpId:            data.GCPID,
		Name:             data.Name,
		ProjectId:        data.ProjectID,
		GcpProjectId:     data.GCPProjectID,
		VerificationCode: data.VerificationCode,
		CreatedAt:        data.CreatedAt.Unix(),
		UpdatedAt:        data.UpdatedAt.Unix(),
	}
	return &gcp
}

func (g *GoogleService) ListGCP(ctx context.Context, req *google.ListGCPRequest) (*google.ListGCPResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := g.repository.ListGCP(ctx, req.ProjectId, req.GcpId, req.GcpProjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &google.ListGCPResponse{}, nil
		}
		return nil, err
	}
	data := google.ListGCPResponse{}
	for _, d := range *list {
		data.Gcp = append(data.Gcp, convertGCP(&d))
	}
	return &data, nil
}

func (g *GoogleService) GetGCP(ctx context.Context, req *google.GetGCPRequest) (*google.GetGCPResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := g.repository.GetGCP(ctx, req.ProjectId, req.GcpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &google.GetGCPResponse{}, nil
		}
		return nil, err
	}
	return &google.GetGCPResponse{Gcp: convertGCP(data)}, nil
}

func (g *GoogleService) PutGCP(ctx context.Context, req *google.PutGCPRequest) (*google.PutGCPResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registerd, err := g.repository.UpsertGCP(ctx, req.Gcp)
	if err != nil {
		return nil, err
	}
	return &google.PutGCPResponse{Gcp: convertGCP(registerd)}, nil
}

func (g *GoogleService) DeleteGCP(ctx context.Context, req *google.DeleteGCPRequest) (*google.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := g.repository.ListGCPDataSource(ctx, req.ProjectId, req.GcpId)
	if err != nil {
		return nil, err
	}
	for _, ds := range *list {
		if err := g.repository.DeleteGCPDataSource(ctx, req.ProjectId, req.GcpId, ds.GoogleDataSourceID); err != nil {
			return nil, err
		}
	}
	if err := g.repository.DeleteGCP(ctx, req.ProjectId, req.GcpId); err != nil {
		return nil, err
	}
	return &google.Empty{}, nil
}

func (g *GoogleService) ListGCPDataSource(ctx context.Context, req *google.ListGCPDataSourceRequest) (*google.ListGCPDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := g.repository.ListGCPDataSource(ctx, req.ProjectId, req.GcpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &google.ListGCPDataSourceResponse{}, nil
		}
		return nil, err
	}
	data := google.ListGCPDataSourceResponse{}
	for _, d := range *list {
		data.GcpDataSource = append(data.GcpDataSource, convertGCPDataSource(&d))
	}
	return &data, nil
}

func convertGCPDataSource(data *db.GCPDataSource) *google.GCPDataSource {
	if data == nil {
		return &google.GCPDataSource{}
	}
	gcp := google.GCPDataSource{
		GcpId:              data.GCPID,
		GoogleDataSourceId: data.GoogleDataSourceID,
		SpecificVersion:    data.SpecificVersion,
		ProjectId:          data.ProjectID,
		Status:             getStatus(data.Status),
		StatusDetail:       data.StatusDetail,
		CreatedAt:          data.CreatedAt.Unix(),
		UpdatedAt:          data.UpdatedAt.Unix(),
		Name:               data.Name,         // google_data_source.name
		MaxScore:           data.MaxScore,     // google_data_source.max_score
		Description:        data.Description,  // google_data_source.description
		GcpProjectId:       data.GCPProjectID, // gcp.gcp_project_id
	}
	if !zero.IsZeroVal(data.ScanAt) {
		gcp.ScanAt = data.ScanAt.Unix()
	}
	return &gcp
}

func getStatus(s string) google.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := google.Status_value[statusKey]; !ok {
		return google.Status_UNKNOWN
	}
	switch statusKey {
	case google.Status_OK.String():
		return google.Status_OK
	case google.Status_CONFIGURED.String():
		return google.Status_CONFIGURED
	case google.Status_IN_PROGRESS.String():
		return google.Status_IN_PROGRESS
	case google.Status_ERROR.String():
		return google.Status_ERROR
	default:
		return google.Status_UNKNOWN
	}
}

func (g *GoogleService) GetGCPDataSource(ctx context.Context, req *google.GetGCPDataSourceRequest) (*google.GetGCPDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := g.repository.GetGCPDataSource(ctx, req.ProjectId, req.GcpId, req.GoogleDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &google.GetGCPDataSourceResponse{}, nil
		}
		return nil, err
	}
	return &google.GetGCPDataSourceResponse{GcpDataSource: convertGCPDataSource(data)}, nil
}

func (g *GoogleService) AttachGCPDataSource(ctx context.Context, req *google.AttachGCPDataSourceRequest) (*google.AttachGCPDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	gcp, err := g.repository.GetGCP(ctx, req.ProjectId, req.GcpDataSource.GcpId)
	if err != nil {
		return nil, err
	}
	if ok, err := g.gcpClient.VerifyCode(ctx, gcp.GCPProjectID, gcp.VerificationCode); !ok || err != nil {
		return nil, err
	}
	registered, err := g.repository.UpsertGCPDataSource(ctx, req.GcpDataSource)
	if err != nil {
		return nil, err
	}
	if !registered.ErrorNotifiedAt.IsZero() && registered.Status != google.Status_ERROR.String() {
		if err := g.repository.UpdateGCPErrorNotifiedAt(ctx, gorm.Expr("NULL"), registered.GCPID, registered.GoogleDataSourceID, registered.ProjectID); err != nil {
			return nil, err
		}
	}
	return &google.AttachGCPDataSourceResponse{GcpDataSource: convertGCPDataSource(registered)}, nil
}

func (g *GoogleService) DetachGCPDataSource(ctx context.Context, req *google.DetachGCPDataSourceRequest) (*google.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	err := g.repository.DeleteGCPDataSource(ctx, req.ProjectId, req.GcpId, req.GoogleDataSourceId)
	if err != nil {
		return nil, err
	}
	return &google.Empty{}, nil
}

const (
	cloudAssetDataSourceID  uint32 = 1001
	cloudSploitDataSourceID uint32 = 1002
	sccDataSourceID         uint32 = 1003
	portscanDataSourceID    uint32 = 1004
)

func (g *GoogleService) InvokeScanGCP(ctx context.Context, req *google.InvokeScanGCPRequest) (*google.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	gcp, err := g.repository.GetGCP(ctx, req.ProjectId, req.GcpId)
	if err != nil {
		return nil, err
	}
	data, err := g.repository.GetGCPDataSource(ctx, req.ProjectId, req.GcpId, req.GoogleDataSourceId)
	if err != nil {
		return nil, err
	}
	if ok, err := g.gcpClient.VerifyCode(ctx, gcp.GCPProjectID, gcp.VerificationCode); !ok || err != nil {
		if _, upErr := g.repository.UpsertGCPDataSource(ctx, &google.GCPDataSourceForUpsert{
			GcpId:              data.GCPID,
			GoogleDataSourceId: data.GoogleDataSourceID,
			ProjectId:          data.ProjectID,
			Status:             google.Status_ERROR,
			StatusDetail:       err.Error(),
			ScanAt:             data.ScanAt.Unix(),
			SpecificVersion:    data.SpecificVersion,
		}); upErr != nil {
			return nil, upErr
		}
		g.logger.Warnf(ctx, "Failed to verify code: gcp_id=%d, gcp_project_id=%s, err=%v", data.GCPID, data.GCPProjectID, err)
		return &google.Empty{}, nil
	}
	msg := &message.GCPQueueMessage{
		GCPID:              data.GCPID,
		ProjectID:          data.ProjectID,
		GoogleDataSourceID: data.GoogleDataSourceID,
		ScanOnly:           req.ScanOnly,
	}
	var resp *sqs.SendMessageOutput
	switch data.GoogleDataSourceID {
	case cloudAssetDataSourceID:
		resp, err = g.sqs.Send(ctx, g.sqs.GoogleAssetQueueURL, msg)
	case cloudSploitDataSourceID:
		if data.SpecificVersion == "" {
			resp, err = g.sqs.Send(ctx, g.sqs.GoogleCloudSploitQueueURL, msg)
		} else {
			resp, err = g.sqs.Send(ctx, g.sqs.GoogleCloudSploitOldQueueURL, msg)
		}
	case sccDataSourceID:
		resp, err = g.sqs.Send(ctx, g.sqs.GoogleSCCQueueURL, msg)
	case portscanDataSourceID:
		resp, err = g.sqs.Send(ctx, g.sqs.GooglePortscanQueueURL, msg)
	default:
		return nil, fmt.Errorf("Unknown googleDataSourceID: %d", data.GoogleDataSourceID)
	}
	if err != nil {
		return nil, err
	}
	if _, err = g.repository.UpsertGCPDataSource(ctx, &google.GCPDataSourceForUpsert{
		GcpId:              data.GCPID,
		GoogleDataSourceId: data.GoogleDataSourceID,
		ProjectId:          data.ProjectID,
		Status:             google.Status_IN_PROGRESS,
		StatusDetail:       fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:             data.ScanAt.Unix(),
		SpecificVersion:    data.SpecificVersion,
	}); err != nil {
		return nil, err
	}
	g.logger.Infof(ctx, "Invoke scanned: gcp_project_id=%s, messageId=%v", data.GCPProjectID, resp.MessageId)
	return &google.Empty{}, nil
}

func (g *GoogleService) InvokeScanAll(ctx context.Context, req *google.InvokeScanAllRequest) (*google.Empty, error) {
	list, err := g.repository.ListGCPDataSourceByDataSourceID(ctx, req.GoogleDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &google.Empty{}, nil
		}
		return nil, err
	}
	for _, gcp := range *list {
		if resp, err := g.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: gcp.ProjectID}); err != nil {
			g.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			g.logger.Infof(ctx, "Skip deactive project, project_id=%d", gcp.ProjectID)
			continue
		}

		if _, err := g.InvokeScanGCP(ctx, &google.InvokeScanGCPRequest{
			GcpId:              gcp.GCPID,
			ProjectId:          gcp.ProjectID,
			GoogleDataSourceId: gcp.GoogleDataSourceID,
			ScanOnly:           true,
		}); err != nil {
			// In GCP, an error may occur during InvokeScan due to user misconfiguration(e.g. invalid verification_code).
			// But to avoid having a single error stop the entire process, a notification log is output and other processes continue.
			g.logger.Notifyf(ctx, logging.ErrorLevel, "InvokeScanGCP error occured: gcp_id=%d, gcp_project_id=%s, err=%+v", gcp.GCPID, gcp.GCPProjectID, err)
			continue
		}
	}
	return &google.Empty{}, nil
}
