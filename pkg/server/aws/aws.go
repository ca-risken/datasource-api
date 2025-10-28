package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/aws"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (a *AWSService) ListAWS(ctx context.Context, req *aws.ListAWSRequest) (*aws.ListAWSResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListAWS(ctx, req.ProjectId, req.AwsId, req.AwsAccountId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &aws.ListAWSResponse{}, nil
		}
		return nil, err
	}
	data := aws.ListAWSResponse{}
	for _, d := range *list {
		data.Aws = append(data.Aws, convertAWS(&d))
	}
	return &data, nil
}

func (a *AWSService) PutAWS(ctx context.Context, req *aws.PutAWSRequest) (*aws.PutAWSResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := a.repository.GetAWSByAccountID(ctx, req.ProjectId, req.Aws.AwsAccountId)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var awsID uint32
	if !noRecord {
		awsID = savedData.AWSID
	}
	data := &model.AWS{
		AWSID:        awsID,
		Name:         req.Aws.Name,
		ProjectID:    req.Aws.ProjectId,
		AWSAccountID: req.Aws.AwsAccountId,
	}

	// aws upsert
	registerdData, err := a.repository.UpsertAWS(ctx, data)
	if err != nil {
		return nil, err
	}
	return &aws.PutAWSResponse{Aws: convertAWS(registerdData)}, nil
}

func (a *AWSService) DeleteAWS(ctx context.Context, req *aws.DeleteAWSRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListAWSRelDataSource(ctx, req.ProjectId, req.AwsId)
	if err != nil {
		return nil, err
	}
	for _, ds := range *list {
		if err := a.repository.DeleteAWSRelDataSource(ctx, req.ProjectId, req.AwsId, ds.AWSDataSourceID); err != nil {
			return nil, err
		}
	}
	if err := a.repository.DeleteAWS(ctx, req.ProjectId, req.AwsId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func getStatus(s string) aws.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := aws.Status_value[statusKey]; !ok {
		return aws.Status_UNKNOWN
	}
	switch statusKey {
	case aws.Status_OK.String():
		return aws.Status_OK
	case aws.Status_CONFIGURED.String():
		return aws.Status_CONFIGURED
	case aws.Status_IN_PROGRESS.String():
		return aws.Status_IN_PROGRESS
	case aws.Status_ERROR.String():
		return aws.Status_ERROR
	default:
		return aws.Status_UNKNOWN
	}
}

func (a *AWSService) ListDataSource(ctx context.Context, req *aws.ListDataSourceRequest) (*aws.ListDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := a.repository.ListAWSDataSource(ctx, req.ProjectId, req.AwsId, req.DataSource)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &aws.ListDataSourceResponse{DataSource: []*aws.DataSource{}}, nil
		}
		return nil, err
	}
	ds := []*aws.DataSource{}
	for _, d := range *list {
		var scanAt int64
		if !d.ScanAt.IsZero() {
			scanAt = d.ScanAt.Unix()
		}
		ds = append(ds, &aws.DataSource{
			AwsDataSourceId: d.AWSDataSourceID,
			DataSource:      d.DataSource,
			MaxScore:        d.MaxScore,
			AwsId:           d.AWSID,
			ProjectId:       d.ProjectID,
			AssumeRoleArn:   d.AssumeRoleArn,
			ExternalId:      d.ExternalID,
			SpecificVersion: d.SpecificVersion,
			Status:          getStatus(d.Status),
			StatusDetail:    d.StatusDetail,
			ScanAt:          scanAt,
		})
	}
	return &aws.ListDataSourceResponse{DataSource: ds}, nil
}

func (a *AWSService) AttachDataSource(ctx context.Context, req *aws.AttachDataSourceRequest) (*aws.AttachDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registered, err := a.repository.UpsertAWSRelDataSource(ctx, req.AttachDataSource)
	if err != nil {
		return nil, err
	}
	if !registered.ErrorNotifiedAt.IsZero() && registered.Status != aws.Status_ERROR.String() {
		if err := a.repository.UpdateAWSErrorNotifiedAt(ctx, gorm.Expr("NULL"), registered.AWSID, registered.AWSDataSourceID, registered.ProjectID); err != nil {
			return nil, err
		}
	}
	var scanAt int64
	if !registered.ScanAt.IsZero() {
		scanAt = registered.ScanAt.Unix()
	}
	return &aws.AttachDataSourceResponse{DataSource: &aws.AWSRelDataSource{
		AwsId:           registered.AWSID,
		AwsDataSourceId: registered.AWSDataSourceID,
		ProjectId:       registered.ProjectID,
		AssumeRoleArn:   registered.AssumeRoleArn,
		ExternalId:      registered.ExternalID,
		SpecificVersion: registered.SpecificVersion,
		Status:          getStatus(registered.Status),
		StatusDetail:    registered.StatusDetail,
		ScanAt:          scanAt,
		CreatedAt:       registered.CreatedAt.Unix(),
		UpdatedAt:       registered.UpdatedAt.Unix(),
	}}, nil
}

func (a *AWSService) DetachDataSource(ctx context.Context, req *aws.DetachDataSourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := a.repository.DeleteAWSRelDataSource(ctx, req.ProjectId, req.AwsId, req.AwsDataSourceId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (a *AWSService) InvokeScan(ctx context.Context, req *aws.InvokeScanRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	ds, err := a.repository.GetAWSDataSourceForMessage(ctx, req.AwsId, req.AwsDataSourceId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	msg := &message.AWSQueueMessage{
		AWSID:           ds.AWSID,
		AWSDataSourceID: ds.AWSDataSourceID,
		DataSource:      ds.DataSource,
		AccountID:       ds.AWSAccountID,
		ProjectID:       ds.ProjectID,
		AssumeRoleArn:   ds.AssumeRoleArn,
		ExternalID:      ds.ExternalID,
		ScanOnly:        req.ScanOnly,
		SpecificVersion: ds.SpecificVersion,
		FullScan:        req.FullScan,
	}
	var resp *sqs.SendMessageOutput
	switch msg.DataSource {
	case message.AWSAccessAnalyzerDataSource:
		resp, err = a.sqs.Send(ctx, a.sqs.AWSAccessAnalyzerQueueURL, msg)
	case message.AWSAdminCheckerDataSource:
		resp, err = a.sqs.Send(ctx, a.sqs.AWSAdminCheckerQueueURL, msg)
	case message.AWSCloudSploitDataSource:
		if ds.SpecificVersion == "" {
			resp, err = a.sqs.Send(ctx, a.sqs.AWSCloudSploitQueueURL, msg)
		} else {
			resp, err = a.sqs.Send(ctx, a.sqs.AWSCloudSploitOldQueueURL, msg)
		}
	case message.AWSGuardDutyDataSource:
		resp, err = a.sqs.Send(ctx, a.sqs.AWSGuardDutyQueueURL, msg)
	case message.AWSPortscanDataSource:
		resp, err = a.sqs.Send(ctx, a.sqs.AWSPortscanQueueURL, msg)
	default:
		return nil, fmt.Errorf("unknown datasource, datasource=%s", msg.DataSource)
	}
	if err != nil {
		return nil, err
	}
	data, err := a.repository.GetAWSRelDataSourceByID(ctx, msg.AWSID, msg.AWSDataSourceID, msg.ProjectID)
	if err != nil {
		return nil, err
	}
	if _, err = a.repository.UpsertAWSRelDataSource(ctx, &aws.DataSourceForAttach{
		AwsId:           data.AWSID,
		AwsDataSourceId: data.AWSDataSourceID,
		ProjectId:       data.ProjectID,
		AssumeRoleArn:   data.AssumeRoleArn,
		ExternalId:      data.ExternalID,
		SpecificVersion: data.SpecificVersion,
		Status:          aws.Status_IN_PROGRESS,
		StatusDetail:    fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:          data.ScanAt.Unix(),
	}); err != nil {
		return nil, err
	}
	a.logger.Infof(ctx, "Invoke scanned, messageId: %v", resp.MessageId)
	return &empty.Empty{}, nil
}

func convertAWS(data *model.AWS) *aws.AWS {
	if data == nil {
		return &aws.AWS{}
	}
	return &aws.AWS{
		AwsId:        data.AWSID,
		Name:         data.Name,
		ProjectId:    data.ProjectID,
		AwsAccountId: data.AWSAccountID,
		CreatedAt:    data.CreatedAt.Unix(),
		UpdatedAt:    data.CreatedAt.Unix(),
	}
}

func (a *AWSService) InvokeScanAll(ctx context.Context, req *aws.InvokeScanAllRequest) (*empty.Empty, error) {
	a.logger.Infof(ctx, "Start InvokeScanAll, param: %+v", req)
	list, err := a.repository.ListDataSourceByAWSDataSourceID(ctx, req.AwsDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &empty.Empty{}, nil
		}
		return nil, err
	}
	for _, dataSource := range *list {
		if resp, err := a.projectClient.IsActive(ctx, &project.IsActiveRequest{ProjectId: dataSource.ProjectID}); err != nil {
			a.logger.Errorf(ctx, "Failed to project.IsActive API, err=%+v", err)
			return nil, err
		} else if !resp.Active {
			a.logger.Infof(ctx, "Skip deactive project, project_id=%d", dataSource.ProjectID)
			continue
		}
		if _, err := a.InvokeScan(ctx, &aws.InvokeScanRequest{
			ProjectId:       dataSource.ProjectID,
			AwsId:           dataSource.AWSID,
			AwsDataSourceId: dataSource.AWSDataSourceID,
			ScanOnly:        true,
			FullScan:        req.FullScan,
		}); err != nil {
			a.logger.Errorf(ctx, "AWS InvokeScan error: project_id=%d, aws_id=%d, aws_datasource_id=%d, err=%+v",
				dataSource.ProjectID, dataSource.AWSID, dataSource.AWSDataSourceID, err)
			return nil, err
		}
	}
	a.logger.Info(ctx, "End InvokeScanAll")
	return &empty.Empty{}, nil
}
