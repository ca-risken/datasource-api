package diagnosis

import (
	"context"
	"errors"
	"time"

	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (d *DiagnosisService) ListWpscanSetting(ctx context.Context, req *diagnosis.ListWpscanSettingRequest) (*diagnosis.ListWpscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := d.repository.ListWpscanSetting(ctx, req.ProjectId, req.DiagnosisDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListWpscanSettingResponse{}, nil
		}
		d.logger.Errorf(ctx, "Failed to List WpscanSettinng, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListWpscanSettingResponse{}
	for _, d := range *list {
		data.WpscanSetting = append(data.WpscanSetting, convertWpscanSetting(&d))
	}
	return &data, nil
}

func (d *DiagnosisService) GetWpscanSetting(ctx context.Context, req *diagnosis.GetWpscanSettingRequest) (*diagnosis.GetWpscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := d.repository.GetWpscanSetting(ctx, req.ProjectId, req.WpscanSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		d.logger.Errorf(ctx, "Failed to Get WpscanSettinng, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetWpscanSettingResponse{WpscanSetting: convertWpscanSetting(getData)}, nil
}

func (d *DiagnosisService) PutWpscanSetting(ctx context.Context, req *diagnosis.PutWpscanSettingRequest) (*diagnosis.PutWpscanSettingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := d.repository.GetWpscanSetting(ctx, req.ProjectId, req.WpscanSetting.WpscanSettingId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		d.logger.Errorf(ctx, "Failed to Get WpscanSetting, error: %v", err)
		return nil, err
	}

	var wpscanSettingID uint32
	if !noRecord {
		wpscanSettingID = savedData.WpscanSettingID
	}
	data := &model.WpscanSetting{
		WpscanSettingID:       wpscanSettingID,
		ProjectID:             req.ProjectId,
		DiagnosisDataSourceID: req.WpscanSetting.DiagnosisDataSourceId,
		TargetURL:             req.WpscanSetting.TargetUrl,
		Options:               req.WpscanSetting.Options,
		Status:                req.WpscanSetting.Status.String(),
		StatusDetail:          req.WpscanSetting.StatusDetail,
		ScanAt:                time.Unix(req.WpscanSetting.ScanAt, 0),
	}

	registerdData, err := d.repository.UpsertWpscanSetting(ctx, data)
	if err != nil {
		d.logger.Errorf(ctx, "Failed to Put WpscanSetting, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutWpscanSettingResponse{WpscanSetting: convertWpscanSetting(registerdData)}, nil
}

func (d *DiagnosisService) DeleteWpscanSetting(ctx context.Context, req *diagnosis.DeleteWpscanSettingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := d.repository.DeleteWpscanSetting(ctx, req.ProjectId, req.WpscanSettingId); err != nil {
		d.logger.Errorf(ctx, "Failed to Delete WpscanSettinng, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertWpscanSetting(data *model.WpscanSetting) *diagnosis.WpscanSetting {
	if data == nil {
		return &diagnosis.WpscanSetting{}
	}
	return &diagnosis.WpscanSetting{
		WpscanSettingId:       data.WpscanSettingID,
		DiagnosisDataSourceId: data.DiagnosisDataSourceID,
		ProjectId:             data.ProjectID,
		TargetUrl:             data.TargetURL,
		Options:               data.Options,
		CreatedAt:             data.CreatedAt.Unix(),
		UpdatedAt:             data.CreatedAt.Unix(),
		Status:                getStatus(data.Status),
		StatusDetail:          data.StatusDetail,
		ScanAt:                data.ScanAt.Unix(),
	}
}

func makeWpscanMessage(ProjectID, SettingID uint32, targetURL, options string) (*message.WpscanQueueMessage, error) {
	msg := &message.WpscanQueueMessage{
		DataSource:      message.DataSourceNameWPScan,
		WpscanSettingID: SettingID,
		ProjectID:       ProjectID,
		TargetURL:       targetURL,
		Options:         options,
	}
	return msg, nil
}
