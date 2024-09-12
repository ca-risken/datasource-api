package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type ScanErrors struct {
	awsErrors       []*db.AWSScanError
	gcpErrors       []*db.GCPScanError
	githubErrors    []*db.GitHubScanError
	diagnosisErrors []*db.DiagnosisScanError
	osintErrors     []*db.OsintScanError
	azureErrors     []*db.AzureScanError
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
	// Diagnosis
	diagnosisList, err := d.dbClient.ListDiagnosisScanErrorForNotify(ctx)
	if err != nil {
		return nil, err
	}
	for _, diagnosis := range diagnosisList {
		if _, ok := scanErrors[diagnosis.ProjectID]; !ok {
			scanErrors[diagnosis.ProjectID] = &ScanErrors{}
		}
		scanErrors[diagnosis.ProjectID].diagnosisErrors = append(scanErrors[diagnosis.ProjectID].diagnosisErrors, diagnosis)
	}
	// OSINT
	osintList, err := d.dbClient.ListOsintScanErrorForNotify(ctx)
	if err != nil {
		return nil, err
	}
	for _, osint := range osintList {
		if _, ok := scanErrors[osint.ProjectID]; !ok {
			scanErrors[osint.ProjectID] = &ScanErrors{}
		}
		scanErrors[osint.ProjectID].osintErrors = append(scanErrors[osint.ProjectID].osintErrors, osint)
	}
	// Azure
	azureList, err := d.dbClient.ListAzureScanErrorForNotify(ctx)
	if err != nil {
		return nil, err
	}
	for _, azure := range azureList {
		if _, ok := scanErrors[azure.ProjectID]; !ok {
			scanErrors[azure.ProjectID] = &ScanErrors{}
		}
		scanErrors[azure.ProjectID].azureErrors = append(scanErrors[azure.ProjectID].azureErrors, azure)
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
		case github.DataSource == message.CodeScanDataSource:
			if err := d.dbClient.UpdateCodeCodeScanErrorNotifiedAt(ctx, time.Now(), github.CodeGithubSettingID, projectID); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown data source: %s", github.DataSource)
		}
	}
	for _, diagnosis := range errs.diagnosisErrors {
		switch {
		case diagnosis.DataSource == message.DataSourceNameWPScan:
			if err := d.dbClient.UpdateDiagnosisWpscanErrorNotifiedAt(ctx, time.Now(), diagnosis.ScanID, projectID); err != nil {
				return err
			}
		case diagnosis.DataSource == message.DataSourceNamePortScan:
			if err := d.dbClient.UpdateDiagnosisPortscanErrorNotifiedAt(ctx, time.Now(), diagnosis.ScanID, projectID); err != nil {
				return err
			}
		case diagnosis.DataSource == message.DataSourceNameApplicationScan:
			if err := d.dbClient.UpdateDiagnosisAppScanErrorNotifiedAt(ctx, time.Now(), diagnosis.ScanID, projectID); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown data source: %s", diagnosis.DataSource)
		}
	}
	for _, o := range errs.osintErrors {
		if err := d.dbClient.UpdateOsintErrorNotifiedAt(ctx, time.Now(), o.RelOsintDataSourceID, projectID); err != nil {
			return err
		}
	}
	for _, az := range errs.azureErrors {
		if err := d.dbClient.UpdateAzureErrorNotifiedAt(ctx, time.Now(), az.AzureID, az.AzureDataSourceID, projectID); err != nil {
			return err
		}
	}

	return nil
}
