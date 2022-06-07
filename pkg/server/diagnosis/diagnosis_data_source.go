package diagnosis

import (
	"context"
	"errors"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
	"github.com/golang/protobuf/ptypes/empty"
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
