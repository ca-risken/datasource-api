package db

import (
	"context"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/vikyd/zero"
)

type OSINTRepoInterface interface {
	ListOsint(context.Context, uint32) (*[]model.Osint, error)
	GetOsint(context.Context, uint32, uint32) (*model.Osint, error)
	UpsertOsint(context.Context, *model.Osint) (*model.Osint, error)
	DeleteOsint(context.Context, uint32, uint32) error
	ListOsintDataSource(context.Context, uint32, string) (*[]model.OsintDataSource, error)
	GetOsintDataSource(context.Context, uint32, uint32) (*model.OsintDataSource, error)
	UpsertOsintDataSource(context.Context, *model.OsintDataSource) (*model.OsintDataSource, error)
	DeleteOsintDataSource(context.Context, uint32, uint32) error
	ListRelOsintDataSource(context.Context, uint32, uint32, uint32) (*[]model.RelOsintDataSource, error)
	GetRelOsintDataSource(context.Context, uint32, uint32) (*model.RelOsintDataSource, error)
	UpsertRelOsintDataSource(context.Context, *model.RelOsintDataSource) (*model.RelOsintDataSource, error)
	DeleteRelOsintDataSource(context.Context, uint32, uint32) error
	ListOsintDetectWord(context.Context, uint32, uint32) (*[]model.OsintDetectWord, error)
	GetOsintDetectWord(context.Context, uint32, uint32) (*model.OsintDetectWord, error)
	UpsertOsintDetectWord(context.Context, *model.OsintDetectWord) (*model.OsintDetectWord, error)
	DeleteOsintDetectWord(context.Context, uint32, uint32) error
	// For Invoke
	ListAllRelOsintDataSource(context.Context, uint32) (*[]model.RelOsintDataSource, error)
}

var _ OSINTRepoInterface = (*Client)(nil) // verify interface compliance

func (c *Client) ListOsint(ctx context.Context, projectID uint32) (*[]model.Osint, error) {
	query := `select * from osint where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	var data []model.Osint
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetOsint(ctx context.Context, projectID uint32, osintID uint32) (*model.Osint, error) {
	var data model.Osint
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND osint_id = ?", projectID, osintID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertOsint(ctx context.Context, data *model.Osint) (*model.Osint, error) {
	var savedData model.Osint
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND osint_id = ?", data.ProjectID, data.OsintID).Assign(data).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeleteOsint(ctx context.Context, projectID uint32, osintID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND osint_id = ?", projectID, osintID).Delete(&model.Osint{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListOsintDataSource(ctx context.Context, projectID uint32, name string) (*[]model.OsintDataSource, error) {
	var data []model.OsintDataSource
	paramName := "%" + name + "%"
	if err := c.SlaveDB.WithContext(ctx).Where("name like ?", paramName).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetOsintDataSource(ctx context.Context, projectID uint32, osintDataSourceID uint32) (*model.OsintDataSource, error) {
	var data model.OsintDataSource
	if err := c.SlaveDB.WithContext(ctx).Where("osint_data_source_id = ?", osintDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertOsintDataSource(ctx context.Context, input *model.OsintDataSource) (*model.OsintDataSource, error) {
	var data model.OsintDataSource
	if err := c.MasterDB.WithContext(ctx).Where("osint_data_source_id = ?", input.OsintDataSourceID).Assign(&model.OsintDataSource{Name: input.Name, Description: input.Description, MaxScore: input.MaxScore}).FirstOrCreate(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) DeleteOsintDataSource(ctx context.Context, projectID uint32, osintDataSourceID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("osint_data_source_id =  ?", osintDataSourceID).Delete(&model.OsintDataSource{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListRelOsintDataSource(ctx context.Context, projectID, osintID, osintDataSourceID uint32) (*[]model.RelOsintDataSource, error) {
	query := `select * from rel_osint_data_source where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(osintID) {
		query += " and osint_id = ?"
		params = append(params, osintID)
	}
	if !zero.IsZeroVal(osintDataSourceID) {
		query += " and osint_data_source_id = ?"
		params = append(params, osintDataSourceID)
	}
	var data []model.RelOsintDataSource
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetRelOsintDataSource(ctx context.Context, projectID uint32, relOsintDataSourceID uint32) (*model.RelOsintDataSource, error) {
	var data model.RelOsintDataSource
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND rel_osint_data_source_id = ?", projectID, relOsintDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertRelOsintDataSource(ctx context.Context, data *model.RelOsintDataSource) (*model.RelOsintDataSource, error) {
	var savedData model.RelOsintDataSource
	update := relOsintDataSourceToMap(data)
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND rel_osint_data_source_id = ?", data.ProjectID, data.RelOsintDataSourceID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeleteRelOsintDataSource(ctx context.Context, projectID uint32, relOsintDataSourceID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND rel_osint_data_source_id = ?", projectID, relOsintDataSourceID).Delete(&model.RelOsintDataSource{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAllRelOsintDataSource(ctx context.Context, osintDataSourceID uint32) (*[]model.RelOsintDataSource, error) {
	query := `select * from rel_osint_data_source`
	var params []interface{}
	if !zero.IsZeroVal(osintDataSourceID) {
		query += " where osint_data_source_id = ?"
		params = append(params, osintDataSourceID)
	}
	var data []model.RelOsintDataSource
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func relOsintDataSourceToMap(relOsintDataSource *model.RelOsintDataSource) map[string]interface{} {
	return map[string]interface{}{
		"rel_osint_data_source_id": relOsintDataSource.RelOsintDataSourceID,
		"osint_id":                 relOsintDataSource.OsintID,
		"osint_data_source_id":     relOsintDataSource.OsintDataSourceID,
		"project_id":               relOsintDataSource.ProjectID,
		"status":                   relOsintDataSource.Status,
		"status_detail":            relOsintDataSource.StatusDetail,
		"scan_at":                  relOsintDataSource.ScanAt,
	}
}

func (c *Client) ListOsintDetectWord(ctx context.Context, projectID, relOsintDataSourceID uint32) (*[]model.OsintDetectWord, error) {
	query := `select * from osint_detect_word where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(relOsintDataSourceID) {
		query += " and rel_osint_data_source_id = ?"
		params = append(params, relOsintDataSourceID)
	}
	var data []model.OsintDetectWord
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetOsintDetectWord(ctx context.Context, projectID uint32, osintDetectWordID uint32) (*model.OsintDetectWord, error) {
	var data model.OsintDetectWord
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND osint_detect_word_id = ?", projectID, osintDetectWordID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertOsintDetectWord(ctx context.Context, data *model.OsintDetectWord) (*model.OsintDetectWord, error) {
	var savedData model.OsintDetectWord
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND osint_detect_word_id = ?", data.ProjectID, data.OsintDetectWordID).Assign(data).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeleteOsintDetectWord(ctx context.Context, projectID uint32, osintDetectWordID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND osint_detect_word_id = ?", projectID, osintDetectWordID).Delete(&model.OsintDetectWord{}).Error; err != nil {
		return err
	}
	return nil
}
