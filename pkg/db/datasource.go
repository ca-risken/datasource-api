package db

import (
	"context"
	"fmt"
)

type DataSourceRepoInterface interface {
	CleanWithNoProject(ctx context.Context) error
}

var _ DataSourceRepoInterface = (*Client)(nil) // verify interface compliance

const (
	cleanTableWithNoProjectTemplate = "delete tbl from %s tbl where not exists(select * from project p where p.project_id = tbl.project_id)"
	cleanTableWithNoGithubTemplate  = "delete tbl from %s tbl where not exists(select * from code_github_setting github where github.code_github_setting_id = tbl.code_github_setting_id)"
)

func (c *Client) CleanWithNoProject(ctx context.Context) error {
	// AWS
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "aws")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "aws_rel_data_source")).Error; err != nil {
		return err
	}

	// Google
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "gcp")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "gcp_data_source")).Error; err != nil {
		return err
	}

	// Diagnosis
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "wpscan_setting")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "portscan_setting")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "portscan_target")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "application_scan")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "application_scan_basic_setting")).Error; err != nil {
		return err
	}

	// OSINT
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "osint")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "rel_osint_data_source")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "osint_detect_word")).Error; err != nil {
		return err
	}

	// Code
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "code_github_setting")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "code_gitleaks_setting")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoProjectTemplate, "code_dependency_setting")).Error; err != nil {
		return err
	}
	if err := c.MasterDB.WithContext(ctx).Exec(fmt.Sprintf(cleanTableWithNoGithubTemplate, "code_gitleaks_cache")).Error; err != nil {
		return err
	}

	return nil
}
