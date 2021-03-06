package osint

import (
	"context"
	"errors"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (o *OsintService) ListOsint(ctx context.Context, req *osint.ListOsintRequest) (*osint.ListOsintResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := o.repository.ListOsint(ctx, req.ProjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.ListOsintResponse{}, nil
		}
		o.logger.Errorf(ctx, "Failed to List Osint. error: %v", err)
		return nil, err
	}
	data := osint.ListOsintResponse{}
	for _, d := range *list {
		data.Osint = append(data.Osint, convertOsint(&d))
	}
	return &data, nil
}

func (o *OsintService) GetOsint(ctx context.Context, req *osint.GetOsintRequest) (*osint.GetOsintResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := o.repository.GetOsint(ctx, req.ProjectId, req.OsintId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.GetOsintResponse{}, nil
		}
		o.logger.Errorf(ctx, "Failed to Get Osint. error: %v", err)
		return nil, err
	}

	return &osint.GetOsintResponse{Osint: convertOsint(getData)}, nil
}

func (o *OsintService) PutOsint(ctx context.Context, req *osint.PutOsintRequest) (*osint.PutOsintResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.Osint{
		OsintID:      req.Osint.OsintId,
		ResourceType: req.Osint.ResourceType,
		ResourceName: req.Osint.ResourceName,
		ProjectID:    req.Osint.ProjectId,
	}

	registerdData, err := o.repository.UpsertOsint(ctx, data)
	if err != nil {
		o.logger.Errorf(ctx, "Failed to Put Osint. error: %v", err)
		return nil, err
	}
	return &osint.PutOsintResponse{Osint: convertOsint(registerdData)}, nil
}

func (o *OsintService) DeleteOsint(ctx context.Context, req *osint.DeleteOsintRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	relOsintDataSources, err := o.repository.ListRelOsintDataSource(ctx, req.ProjectId, req.OsintId, 0)
	if err != nil {
		o.logger.Errorf(ctx, "Failed to List RelOsintDataSource when delete osint. error: %v", err)
		return nil, err
	}

	for _, relOsintDataSource := range *relOsintDataSources {
		if err := o.deleteRelOsintDataSourceWithDetectWord(ctx, relOsintDataSource.ProjectID, relOsintDataSource.RelOsintDataSourceID); err != nil {
			o.logger.Errorf(ctx, "Failed to DeleteRelOsintDataSource. error: %v", err)
			return nil, err
		}
	}

	if err := o.repository.DeleteOsint(ctx, req.ProjectId, req.OsintId); err != nil {
		o.logger.Errorf(ctx, "Failed to DeleteOsint. error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertOsint(data *model.Osint) *osint.Osint {
	if data == nil {
		return &osint.Osint{}
	}
	return &osint.Osint{
		OsintId:      data.OsintID,
		ResourceType: data.ResourceType,
		ResourceName: data.ResourceName,
		ProjectId:    data.ProjectID,
		CreatedAt:    data.CreatedAt.Unix(),
		UpdatedAt:    data.CreatedAt.Unix(),
	}
}
