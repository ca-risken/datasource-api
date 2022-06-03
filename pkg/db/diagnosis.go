package db

import (
	"context"

	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/vikyd/zero"
)

type DiagnosisRepoInterface interface {
	ListDiagnosisDataSource(context.Context, uint32, string) (*[]model.DiagnosisDataSource, error)
	GetDiagnosisDataSource(context.Context, uint32, uint32) (*model.DiagnosisDataSource, error)
	UpsertDiagnosisDataSource(context.Context, *model.DiagnosisDataSource) (*model.DiagnosisDataSource, error)
	DeleteDiagnosisDataSource(context.Context, uint32, uint32) error
	ListWpscanSetting(context.Context, uint32, uint32) (*[]model.WpscanSetting, error)
	GetWpscanSetting(context.Context, uint32, uint32) (*model.WpscanSetting, error)
	UpsertWpscanSetting(context.Context, *model.WpscanSetting) (*model.WpscanSetting, error)
	DeleteWpscanSetting(context.Context, uint32, uint32) error
	ListPortscanSetting(context.Context, uint32, uint32) (*[]model.PortscanSetting, error)
	GetPortscanSetting(context.Context, uint32, uint32) (*model.PortscanSetting, error)
	UpsertPortscanSetting(context.Context, *model.PortscanSetting) (*model.PortscanSetting, error)
	DeletePortscanSetting(context.Context, uint32, uint32) error
	ListPortscanTarget(context.Context, uint32, uint32) (*[]model.PortscanTarget, error)
	GetPortscanTarget(context.Context, uint32, uint32) (*model.PortscanTarget, error)
	UpsertPortscanTarget(context.Context, *model.PortscanTarget) (*model.PortscanTarget, error)
	DeletePortscanTarget(context.Context, uint32, uint32) error
	DeletePortscanTargetByPortscanSettingID(context.Context, uint32, uint32) error
	ListApplicationScan(context.Context, uint32, uint32) (*[]model.ApplicationScan, error)
	GetApplicationScan(context.Context, uint32, uint32) (*model.ApplicationScan, error)
	UpsertApplicationScan(context.Context, *model.ApplicationScan) (*model.ApplicationScan, error)
	DeleteApplicationScan(context.Context, uint32, uint32) error
	ListApplicationScanBasicSetting(context.Context, uint32, uint32) (*[]model.ApplicationScanBasicSetting, error)
	GetApplicationScanBasicSetting(context.Context, uint32, uint32) (*model.ApplicationScanBasicSetting, error)
	UpsertApplicationScanBasicSetting(context.Context, *model.ApplicationScanBasicSetting) (*model.ApplicationScanBasicSetting, error)
	DeleteApplicationScanBasicSetting(context.Context, uint32, uint32) error

	//for InvokeScan
	ListAllWpscanSetting(context.Context) (*[]model.WpscanSetting, error)
}

var _ DiagnosisRepoInterface = (*Client)(nil) // verify interface compliance

func (c *Client) ListDiagnosisDataSource(ctx context.Context, projectID uint32, name string) (*[]model.DiagnosisDataSource, error) {
	var data []model.DiagnosisDataSource
	paramName := "%" + name + "%"
	if err := c.SlaveDB.WithContext(ctx).Where("name like ?", paramName).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetDiagnosisDataSource(ctx context.Context, projectID uint32, diagnosisDataSourceID uint32) (*model.DiagnosisDataSource, error) {
	var data model.DiagnosisDataSource
	if err := c.SlaveDB.WithContext(ctx).Where("diagnosis_data_source_id = ?", diagnosisDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertDiagnosisDataSource(ctx context.Context, input *model.DiagnosisDataSource) (*model.DiagnosisDataSource, error) {
	var data model.DiagnosisDataSource
	if err := c.MasterDB.WithContext(ctx).Where("diagnosis_data_source_id = ?", input.DiagnosisDataSourceID).Assign(model.DiagnosisDataSource{Name: input.Name, Description: input.Description, MaxScore: input.MaxScore}).FirstOrCreate(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) DeleteDiagnosisDataSource(ctx context.Context, projectID uint32, diagnosisDataSourceID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("diagnosis_data_source_id =  ?", diagnosisDataSourceID).Delete(model.DiagnosisDataSource{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListWpscanSetting(ctx context.Context, projectID, diagnoosisDataSourceID uint32) (*[]model.WpscanSetting, error) {
	query := `select * from wpscan_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnoosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnoosisDataSourceID)
	}
	var data []model.WpscanSetting
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetWpscanSetting(ctx context.Context, projectID uint32, wpscanSettingID uint32) (*model.WpscanSetting, error) {
	var data model.WpscanSetting
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND wpscan_setting_id = ?", projectID, wpscanSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertWpscanSetting(ctx context.Context, data *model.WpscanSetting) (*model.WpscanSetting, error) {
	var savedData model.WpscanSetting
	update := wpscanSettingToMap(data)
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND wpscan_setting_id = ?", data.ProjectID, data.WpscanSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeleteWpscanSetting(ctx context.Context, projectID uint32, wpscanSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND wpscan_setting_id = ?", projectID, wpscanSettingID).Delete(model.WpscanSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAllWpscanSetting(ctx context.Context) (*[]model.WpscanSetting, error) {
	var data []model.WpscanSetting
	if err := c.SlaveDB.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func wpscanSettingToMap(wpscanSetting *model.WpscanSetting) map[string]interface{} {
	settingMap := map[string]interface{}{
		"wpscan_setting_id":        wpscanSetting.WpscanSettingID,
		"diagnosis_data_source_id": wpscanSetting.DiagnosisDataSourceID,
		"project_id":               wpscanSetting.ProjectID,
		"target_url":               wpscanSetting.TargetURL,
		"status":                   wpscanSetting.Status,
		"options":                  wpscanSetting.Options,
		"status_detail":            wpscanSetting.StatusDetail,
	}
	if !zero.IsZeroVal(wpscanSetting.ScanAt) {
		settingMap["scan_at"] = wpscanSetting.ScanAt
	}
	return settingMap
}

func (c *Client) ListPortscanSetting(ctx context.Context, projectID, diagnoosisDataSourceID uint32) (*[]model.PortscanSetting, error) {
	query := `select * from portscan_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnoosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnoosisDataSourceID)
	}
	var data []model.PortscanSetting
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetPortscanSetting(ctx context.Context, projectID uint32, portscanSettingID uint32) (*model.PortscanSetting, error) {
	var data model.PortscanSetting
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", projectID, portscanSettingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertPortscanSetting(ctx context.Context, data *model.PortscanSetting) (*model.PortscanSetting, error) {
	var savedData model.PortscanSetting
	update := portscanSettingToMap(data)
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", data.ProjectID, data.PortscanSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeletePortscanSetting(ctx context.Context, projectID uint32, portscanSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", projectID, portscanSettingID).Delete(model.PortscanSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListPortscanTarget(ctx context.Context, projectID, portscanSettingID uint32) (*[]model.PortscanTarget, error) {
	query := `select * from portscan_target where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(portscanSettingID) {
		query += " and portscan_setting_id = ?"
		params = append(params, portscanSettingID)
	}
	var data []model.PortscanTarget
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetPortscanTarget(ctx context.Context, projectID uint32, portscanTargetID uint32) (*model.PortscanTarget, error) {
	var data model.PortscanTarget
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND portscan_target_id = ?", projectID, portscanTargetID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertPortscanTarget(ctx context.Context, data *model.PortscanTarget) (*model.PortscanTarget, error) {
	var savedData model.PortscanTarget
	update := portscanTargetToMap(data)
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_target_id = ?", data.ProjectID, data.PortscanTargetID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeletePortscanTarget(ctx context.Context, projectID uint32, portscanTargetID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_target_id = ?", projectID, portscanTargetID).Delete(model.PortscanTarget{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) DeletePortscanTargetByPortscanSettingID(ctx context.Context, projectID uint32, portscanSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND portscan_setting_id = ?", projectID, portscanSettingID).Delete(model.PortscanTarget{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAllPortscanSetting(ctx context.Context) (*[]model.PortscanSetting, error) {
	var data []model.PortscanSetting
	if err := c.SlaveDB.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func portscanSettingToMap(portscanSetting *model.PortscanSetting) map[string]interface{} {
	return map[string]interface{}{
		"portscan_setting_id":      portscanSetting.PortscanSettingID,
		"diagnosis_data_source_id": portscanSetting.DiagnosisDataSourceID,
		"project_id":               portscanSetting.ProjectID,
		"name":                     portscanSetting.Name,
	}
}

func portscanTargetToMap(portscanTarget *model.PortscanTarget) map[string]interface{} {
	settingMap := map[string]interface{}{
		"portscan_Target_id":  portscanTarget.PortscanTargetID,
		"portscan_setting_id": portscanTarget.PortscanSettingID,
		"project_id":          portscanTarget.ProjectID,
		"target":              portscanTarget.Target,
		"status":              portscanTarget.Status,
		"status_detail":       portscanTarget.StatusDetail,
	}
	if !zero.IsZeroVal(portscanTarget.ScanAt) {
		settingMap["scan_at"] = portscanTarget.ScanAt
	}
	return settingMap
}

func (c *Client) ListApplicationScan(ctx context.Context, projectID, diagnoosisDataSourceID uint32) (*[]model.ApplicationScan, error) {
	query := `select * from application_scan where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(diagnoosisDataSourceID) {
		query += " and diagnosis_data_source_id = ?"
		params = append(params, diagnoosisDataSourceID)
	}
	var data []model.ApplicationScan
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetApplicationScan(ctx context.Context, projectID uint32, applicationScanID uint32) (*model.ApplicationScan, error) {
	var data model.ApplicationScan
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", projectID, applicationScanID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertApplicationScan(ctx context.Context, data *model.ApplicationScan) (*model.ApplicationScan, error) {
	var savedData model.ApplicationScan
	update := applicationScanToMap(data)
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", data.ProjectID, data.ApplicationScanID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeleteApplicationScan(ctx context.Context, projectID uint32, applicationScanID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", projectID, applicationScanID).Delete(model.ApplicationScan{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListApplicationScanBasicSetting(ctx context.Context, projectID, applicationScanID uint32) (*[]model.ApplicationScanBasicSetting, error) {
	query := `select * from application_scan_basic_setting where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(applicationScanID) {
		query += " and application_scan_id = ?"
		params = append(params, applicationScanID)
	}
	var data []model.ApplicationScanBasicSetting
	if err := c.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetApplicationScanBasicSetting(ctx context.Context, projectID uint32, applicationScanID uint32) (*model.ApplicationScanBasicSetting, error) {
	var data model.ApplicationScanBasicSetting
	if err := c.SlaveDB.WithContext(ctx).Where("project_id = ? AND application_scan_id = ?", projectID, applicationScanID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertApplicationScanBasicSetting(ctx context.Context, data *model.ApplicationScanBasicSetting) (*model.ApplicationScanBasicSetting, error) {
	var savedData model.ApplicationScanBasicSetting
	update := applicationScanBasicSettingToMap(data)
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_basic_setting_id = ?", data.ProjectID, data.ApplicationScanBasicSettingID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (c *Client) DeleteApplicationScanBasicSetting(ctx context.Context, projectID uint32, applicationScanBasicSettingID uint32) error {
	if err := c.MasterDB.WithContext(ctx).Where("project_id = ? AND application_scan_basic_setting_id = ?", projectID, applicationScanBasicSettingID).Delete(model.ApplicationScanBasicSetting{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ListAllApplicationScan(ctx context.Context) (*[]model.ApplicationScan, error) {
	var data []model.ApplicationScan
	if err := c.SlaveDB.WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func applicationScanToMap(applicationScan *model.ApplicationScan) map[string]interface{} {
	settingMap := map[string]interface{}{
		"application_scan_id":      applicationScan.ApplicationScanID,
		"diagnosis_data_source_id": applicationScan.DiagnosisDataSourceID,
		"project_id":               applicationScan.ProjectID,
		"name":                     applicationScan.Name,
		"scan_type":                applicationScan.ScanType,
		"status":                   applicationScan.Status,
		"status_detail":            applicationScan.StatusDetail,
	}
	if !zero.IsZeroVal(applicationScan.ScanAt) {
		settingMap["scan_at"] = applicationScan.ScanAt
	}
	return settingMap
}

func applicationScanBasicSettingToMap(applicationScanBasicSetting *model.ApplicationScanBasicSetting) map[string]interface{} {
	settingMap := map[string]interface{}{
		"application_scan_basic_setting_id": applicationScanBasicSetting.ApplicationScanBasicSettingID,
		"application_scan_id":               applicationScanBasicSetting.ApplicationScanID,
		"project_id":                        applicationScanBasicSetting.ProjectID,
		"target":                            applicationScanBasicSetting.Target,
		"max_depth":                         applicationScanBasicSetting.MaxDepth,
		"max_children":                      applicationScanBasicSetting.MaxChildren,
	}
	return settingMap
}
