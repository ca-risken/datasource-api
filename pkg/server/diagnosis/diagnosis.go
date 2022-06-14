package diagnosis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (d *DiagnosisService) ListDiagnosisDataSource(ctx context.Context, req *diagnosis.ListDiagnosisDataSourceRequest) (*diagnosis.ListDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := d.repository.ListDiagnosisDataSource(ctx, req.ProjectId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &diagnosis.ListDiagnosisDataSourceResponse{}, nil
		}
		d.logger.Errorf(ctx, "Failed to List DiagnosisDataSource, error: %v", err)
		return nil, err
	}
	data := diagnosis.ListDiagnosisDataSourceResponse{}
	for _, d := range *list {
		data.DiagnosisDataSource = append(data.DiagnosisDataSource, convertDiagnosisDataSource(&d))
	}
	return &data, nil
}

func (d *DiagnosisService) GetDiagnosisDataSource(ctx context.Context, req *diagnosis.GetDiagnosisDataSourceRequest) (*diagnosis.GetDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := d.repository.GetDiagnosisDataSource(ctx, req.ProjectId, req.DiagnosisDataSourceId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		d.logger.Errorf(ctx, "Failed to Get DiagnosisDataSource, error: %v", err)
		return nil, err
	}

	return &diagnosis.GetDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(getData)}, nil
}

func (d *DiagnosisService) PutDiagnosisDataSource(ctx context.Context, req *diagnosis.PutDiagnosisDataSourceRequest) (*diagnosis.PutDiagnosisDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := d.repository.GetDiagnosisDataSource(ctx, req.ProjectId, req.DiagnosisDataSource.DiagnosisDataSourceId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		d.logger.Errorf(ctx, "Failed to Get DiagnosisDataSource, error: %v", err)
		return nil, err
	}

	var diagnosisDataSourceID uint32
	if !noRecord {
		diagnosisDataSourceID = savedData.DiagnosisDataSourceID
	}
	data := &model.DiagnosisDataSource{
		DiagnosisDataSourceID: diagnosisDataSourceID,
		Name:                  req.DiagnosisDataSource.Name,
		Description:           req.DiagnosisDataSource.Description,
		MaxScore:              req.DiagnosisDataSource.MaxScore,
	}

	registerdData, err := d.repository.UpsertDiagnosisDataSource(ctx, data)
	if err != nil {
		d.logger.Errorf(ctx, "Failed to Put DiagnosisDataSource, error: %v", err)
		return nil, err
	}
	return &diagnosis.PutDiagnosisDataSourceResponse{DiagnosisDataSource: convertDiagnosisDataSource(registerdData)}, nil
}

func (d *DiagnosisService) DeleteDiagnosisDataSource(ctx context.Context, req *diagnosis.DeleteDiagnosisDataSourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := d.repository.DeleteDiagnosisDataSource(ctx, req.ProjectId, req.DiagnosisDataSourceId); err != nil {
		d.logger.Errorf(ctx, "Failed to Delete DiagnosisDataSource, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertDiagnosisDataSource(data *model.DiagnosisDataSource) *diagnosis.DiagnosisDataSource {
	if data == nil {
		return &diagnosis.DiagnosisDataSource{}
	}
	return &diagnosis.DiagnosisDataSource{
		DiagnosisDataSourceId: data.DiagnosisDataSourceID,
		Name:                  data.Name,
		Description:           data.Description,
		MaxScore:              data.MaxScore,
		CreatedAt:             data.CreatedAt.Unix(),
		UpdatedAt:             data.CreatedAt.Unix(),
	}
}

func (d *DiagnosisService) InvokeScan(ctx context.Context, req *diagnosis.InvokeScanRequest) (*diagnosis.InvokeScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	dataSource, err := d.repository.GetDiagnosisDataSource(ctx, 0, req.DiagnosisDataSourceId)
	if err != nil {
		return nil, err
	}
	var resp *sqs.SendMessageOutput
	switch dataSource.Name {
	case message.DataSourceNameWPScan:
		data, err := d.repository.GetWpscanSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			return nil, err
		}
		options := data.Options
		if zero.IsZeroVal(options) {
			options = "{}"
		}
		msg, err := makeWpscanMessage(req.ProjectId, req.SettingId, data.TargetURL, options)
		if err != nil {
			d.logger.Errorf(ctx, "Error occured when making WPScan message, error: %v", err)
			return nil, err
		}
		msg.ScanOnly = req.ScanOnly
		resp, err = d.sqs.Send(ctx, d.sqs.DiagnosisWpscanQueueURL, msg)
		if err != nil {
			d.logger.Errorf(ctx, "Error occured when sending WPScan message, error: %v", err)
			return nil, err
		}
		var scanAt time.Time
		if !zero.IsZeroVal(data.ScanAt) {
			scanAt = data.ScanAt
		}
		if _, err = d.repository.UpsertWpscanSetting(ctx, &model.WpscanSetting{
			WpscanSettingID:       data.WpscanSettingID,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			TargetURL:             data.TargetURL,
			Options:               options,
			Status:                diagnosis.Status_IN_PROGRESS.String(),
			StatusDetail:          fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
			ScanAt:                scanAt,
		}); err != nil {
			d.logger.Errorf(ctx, "Error occured when upsert WPScanSetting, error: %v", err)
			return nil, err
		}
	case message.DataSourceNamePortScan:
		data, err := d.repository.GetPortscanSetting(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			d.logger.Errorf(ctx, "Error occured when getting PortscanSetting, error: %v", err)
			return nil, err
		}
		portscanTargets, err := d.repository.ListPortscanTarget(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			d.logger.Errorf(ctx, "Error occured when getting PortscanTargets, error: %v", err)
			return nil, err
		}
		for _, target := range *portscanTargets {
			msg, err := makePortscanMessage(data.ProjectID, data.PortscanSettingID, target.PortscanTargetID, target.Target)
			if err != nil {
				d.logger.Errorf(ctx, "Error occured when making Portscan message, error: %v", err)
				continue
			}
			msg.ScanOnly = req.ScanOnly
			resp, err = d.sqs.Send(ctx, d.sqs.DiagnosisPortscanQueueURL, msg)
			if err != nil {
				d.logger.Errorf(ctx, "Error occured when sending Portscan message, error: %v", err)
				continue
			}
			var scanAt time.Time
			if !zero.IsZeroVal(target.ScanAt) {
				scanAt = target.ScanAt
			}
			if _, err = d.repository.UpsertPortscanTarget(ctx, &model.PortscanTarget{
				PortscanTargetID:  target.PortscanTargetID,
				PortscanSettingID: target.PortscanSettingID,
				ProjectID:         target.ProjectID,
				Target:            target.Target,
				Status:            diagnosis.Status_IN_PROGRESS.String(),
				StatusDetail:      fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
				ScanAt:            scanAt,
			}); err != nil {
				d.logger.Errorf(ctx, "Error occured when upsert Portscan target, error: %v", err)
				return nil, err
			}
		}
		if _, err = d.repository.UpsertPortscanSetting(ctx, &model.PortscanSetting{
			PortscanSettingID:     data.PortscanSettingID,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			Name:                  data.Name,
		}); err != nil {
			return nil, err
		}
	case message.DataSourceNameApplicationScan:
		data, err := d.repository.GetApplicationScan(ctx, req.ProjectId, req.SettingId)
		if err != nil {
			d.logger.Errorf(ctx, "Error occured when getting PortscanSetting, error: %v", err)
			return nil, err
		}
		msg, err := makeApplicationScanMessage(req.ProjectId, req.SettingId, data.Name, data.ScanType)
		if err != nil {
			return nil, err
		}
		msg.ScanOnly = req.ScanOnly
		resp, err = d.sqs.Send(ctx, d.sqs.DiagnosisApplicationScanQueueURL, msg)
		if err != nil {
			return nil, err
		}
		var scanAt time.Time
		if !zero.IsZeroVal(data.ScanAt) {
			scanAt = data.ScanAt
		}
		if _, err = d.repository.UpsertApplicationScan(ctx, &model.ApplicationScan{
			ApplicationScanID:     data.ApplicationScanID,
			DiagnosisDataSourceID: data.DiagnosisDataSourceID,
			ProjectID:             data.ProjectID,
			Name:                  data.Name,
			ScanType:              data.ScanType,
			Status:                diagnosis.Status_IN_PROGRESS.String(),
			StatusDetail:          fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
			ScanAt:                scanAt,
		}); err != nil {
			d.logger.Errorf(ctx, "Error occured when upsert Application scan, error: %v", err)
			return nil, err
		}
	default:
		return nil, nil
	}

	d.logger.Infof(ctx, "Invoke scanned, MessageID: %v", *resp.MessageId)
	return &diagnosis.InvokeScanResponse{Message: "Start Diagnosis."}, nil
}

func (d *DiagnosisService) InvokeScanAll(ctx context.Context, req *diagnosis.InvokeScanAllRequest) (*empty.Empty, error) {
	if !zero.IsZeroVal(req.DiagnosisDataSourceId) {
		dataSource, err := d.repository.GetDiagnosisDataSource(ctx, 0, req.DiagnosisDataSourceId)
		if err != nil {
			return nil, err
		}
		if dataSource.Name != message.DataSourceNameWPScan {
			return &empty.Empty{}, nil
		}
	}

	listWpscanSetting, err := d.repository.ListAllWpscanSetting(ctx)
	if err != nil {
		d.logger.Errorf(ctx, "Failed to List All WPScanSetting., error: %v", err)
		return nil, err
	}
	for _, WpscanSetting := range *listWpscanSetting {
		if resp, err := d.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: WpscanSetting.ProjectID}); err != nil {
			d.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			d.logger.Infof(ctx, "Skip deactive project, project_id=%d", WpscanSetting.ProjectID)
			continue
		}

		if _, err := d.InvokeScan(ctx, &diagnosis.InvokeScanRequest{
			ProjectId:             WpscanSetting.ProjectID,
			SettingId:             WpscanSetting.WpscanSettingID,
			DiagnosisDataSourceId: WpscanSetting.DiagnosisDataSourceID,
			ScanOnly:              true,
		}); err != nil {
			d.logger.Errorf(ctx, "InvokeScanAll error, error: %v", err)
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

func getStatus(s string) diagnosis.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := diagnosis.Status_value[statusKey]; !ok {
		return diagnosis.Status_UNKNOWN
	}
	switch statusKey {
	case diagnosis.Status_OK.String():
		return diagnosis.Status_OK
	case diagnosis.Status_CONFIGURED.String():
		return diagnosis.Status_CONFIGURED
	case diagnosis.Status_IN_PROGRESS.String():
		return diagnosis.Status_IN_PROGRESS
	case diagnosis.Status_ERROR.String():
		return diagnosis.Status_ERROR
	default:
		return diagnosis.Status_UNKNOWN
	}
}
