package datasource

import (
	"context"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/golang/protobuf/ptypes/empty"
)

type DataSourceService struct {
	repository db.DataSourceRepoInterface
	logger     logging.Logger
}

func NewDataSourceService(repo db.DataSourceRepoInterface, l logging.Logger) *DataSourceService {
	return &DataSourceService{
		repository: repo,
		logger:     l,
	}
}

func (d *DataSourceService) CleanDataSource(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	if err := d.repository.CleanWithNoProject(ctx); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
