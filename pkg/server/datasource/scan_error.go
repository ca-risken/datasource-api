package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type ScanErrors struct {
	// TODO: OSINT, Diagnosis
	awsErrors    []*db.AWSScanError
	gcpErrors    []*db.GCPScanError
	githubErrors []*db.GitHubScanError
}

// getScanErrors returns the scan error as a map of scan error data keyed by the project ID
func (d *DataSourceService) getScanErrors(ctx context.Context) (map[uint32]*ScanErrors, error) {
	scanErrors := map[uint32]*ScanErrors{}
	// AWS
	awsList, err := d.dbClient.ListAWSScanErrorForNotify(ctx)
	if err != nil {
		return nil, err
	}
	for _, aws := range awsList {
		if _, ok := scanErrors[aws.ProjectID]; !ok {
			scanErrors[aws.ProjectID] = &ScanErrors{}
		}
		scanErrors[aws.ProjectID].awsErrors = append(scanErrors[aws.ProjectID].awsErrors, aws)
	}
	// GCP
	gcpList, err := d.dbClient.ListGCPScanErrorForNotify(ctx)
	if err != nil {
		return nil, err
	}
	for _, gcp := range gcpList {
		if _, ok := scanErrors[gcp.ProjectID]; !ok {
			scanErrors[gcp.ProjectID] = &ScanErrors{}
		}
		scanErrors[gcp.ProjectID].gcpErrors = append(scanErrors[gcp.ProjectID].gcpErrors, gcp)
	}
	// Code
	githubList, err := d.dbClient.ListCodeGitHubScanErrorForNotify(ctx)
	if err != nil {
		return nil, err
	}
	for _, github := range githubList {
		if _, ok := scanErrors[github.ProjectID]; !ok {
			scanErrors[github.ProjectID] = &ScanErrors{}
		}
		scanErrors[github.ProjectID].githubErrors = append(scanErrors[github.ProjectID].githubErrors, github)
	}

	return scanErrors, nil
}

func (d *DataSourceService) updateScanErrorNotifiedAt(ctx context.Context, projectID uint32, errs *ScanErrors) error {
	for _, aws := range errs.awsErrors {
		if err := d.dbClient.UpdateAWSErrorNotifiedAt(ctx, time.Now(), aws.AWSID, aws.AWSDataSourceID, projectID); err != nil {
			return err
		}
	}
	for _, gcp := range errs.gcpErrors {
		if err := d.dbClient.UpdateGCPErrorNotifiedAt(ctx, time.Now(), gcp.GCPID, gcp.GoogleDataSourceID, projectID); err != nil {
			return err
		}
	}
	for _, github := range errs.githubErrors {
		switch {
		case github.DataSource == message.GitleaksDataSource:
			if err := d.dbClient.UpdateCodeGitleaksErrorNotifiedAt(ctx, time.Now(), github.CodeGithubSettingID, projectID); err != nil {
				return err
			}
		case github.DataSource == message.DependencyDataSource:
			if err := d.dbClient.UpdateCodeDependencyErrorNotifiedAt(ctx, time.Now(), github.CodeGithubSettingID, projectID); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown data source: %s", github.DataSource)
		}
	}
	return nil
}
