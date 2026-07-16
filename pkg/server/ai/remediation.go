package ai

import (
	"context"
	"errors"
	"regexp"
	"strings"

	coreai "github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/message"
	remediationpb "github.com/ca-risken/datasource-api/proto/remediation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	listFindingTagLimit = 200
)

var remediationProposalTargetDataSources = []string{
	message.AWSAccessAnalyzerDataSource,
	message.AWSAdminCheckerDataSource,
	message.AWSCloudSploitDataSource,
	message.AWSPortscanDataSource,
}

var awsAccountIDPattern = regexp.MustCompile(`^[0-9]{12}$`)

func (a *AIService) GenerateRemediationProposal(ctx context.Context, req *remediationpb.GenerateRemediationProposalRequest) (*remediationpb.GenerateRemediationProposalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	targetFinding, err := a.getRemediationProposalTargetFinding(ctx, req.ProjectId, req.FindingId)
	if err != nil {
		return nil, err
	}
	ds, err := a.getAWSDataSourceForRemediationProposal(ctx, req.ProjectId, req.FindingId, targetFinding.DataSource)
	if err != nil {
		return nil, err
	}
	remediationProposalID, err := a.createRemediationProposal(ctx, req.ProjectId, req.FindingId)
	if err != nil {
		return nil, err
	}
	if err := a.sendRemediationProposalMessage(ctx, req.ProjectId, req.FindingId, remediationProposalID, ds); err != nil {
		return nil, err
	}
	return &remediationpb.GenerateRemediationProposalResponse{RemediationProposalId: remediationProposalID}, nil
}

func (a *AIService) getRemediationProposalTargetFinding(ctx context.Context, projectID uint32, findingID uint64) (*finding.Finding, error) {
	findingResp, err := a.findingClient.GetFinding(ctx, &finding.GetFindingRequest{
		ProjectId: projectID,
		FindingId: findingID,
	})
	if err != nil {
		return nil, err
	}
	if findingResp.Finding == nil {
		return nil, status.Errorf(codes.NotFound, "finding not found: finding_id=%d", findingID)
	}
	if !isRemediationProposalTarget(findingResp.Finding.DataSource) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported data_source for remediation proposal: %s", findingResp.Finding.DataSource)
	}
	return findingResp.Finding, nil
}

func (a *AIService) getAWSDataSourceForRemediationProposal(ctx context.Context, projectID uint32, findingID uint64, dataSource string) (*db.DataSource, error) {
	accountID, err := a.getAWSAccountIDFromFindingTag(ctx, projectID, findingID)
	if err != nil {
		return nil, err
	}
	awsData, err := a.dbClient.GetAWSByAccountID(ctx, projectID, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.logger.Warnf(ctx, "AWS account is not registered for remediation proposal: project_id=%d, finding_id=%d, account_id=%s", projectID, findingID, accountID)
			return nil, status.Error(codes.NotFound, "aws remediation target is not configured")
		}
		return nil, err
	}
	awsDataSources, err := a.dbClient.ListAWSDataSource(ctx, projectID, awsData.AWSID, dataSource)
	if err != nil {
		return nil, err
	}
	if awsDataSources == nil || len(*awsDataSources) == 0 {
		a.logger.Warnf(ctx, "AWS data_source is not found for remediation proposal: project_id=%d, finding_id=%d, aws_id=%d, data_source=%s", projectID, findingID, awsData.AWSID, dataSource)
		return nil, status.Error(codes.NotFound, "aws remediation target is not configured")
	}
	ds, err := a.dbClient.GetAWSDataSourceForMessage(ctx, awsData.AWSID, (*awsDataSources)[0].AWSDataSourceID, projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.logger.Warnf(ctx, "AWS data_source is not attached for remediation proposal: project_id=%d, finding_id=%d, aws_id=%d, aws_data_source_id=%d, data_source=%s", projectID, findingID, awsData.AWSID, (*awsDataSources)[0].AWSDataSourceID, dataSource)
			return nil, status.Error(codes.NotFound, "aws remediation target is not configured")
		}
		return nil, err
	}
	if !isAWSAccountIDInAssumeRoleArn(accountID, ds.AssumeRoleArn) {
		a.logger.Warnf(ctx, "AWS account_id does not match assume_role_arn for remediation proposal: project_id=%d, finding_id=%d, aws_id=%d, aws_data_source_id=%d, account_id=%s, assume_role_arn=%s", projectID, findingID, ds.AWSID, ds.AWSDataSourceID, accountID, ds.AssumeRoleArn)
		return nil, status.Error(codes.NotFound, "aws remediation target is not configured")
	}
	return ds, nil
}

func (a *AIService) createRemediationProposal(ctx context.Context, projectID uint32, findingID uint64) (uint32, error) {
	createResp, err := a.coreAIClient.CreateRemediationProposal(ctx, &coreai.CreateRemediationProposalRequest{
		ProjectId: projectID,
		FindingId: findingID,
	})
	if err != nil {
		return 0, err
	}
	if createResp.RemediationProposal == nil || createResp.RemediationProposal.RemediationProposalId == 0 {
		return 0, status.Error(codes.Internal, "failed to create remediation proposal")
	}
	return createResp.RemediationProposal.RemediationProposalId, nil
}

func (a *AIService) sendRemediationProposalMessage(ctx context.Context, projectID uint32, findingID uint64, remediationProposalID uint32, ds *db.DataSource) error {
	msg := &message.RemediationProposalQueueMessage{
		RemediationProposalID: remediationProposalID,
		FindingID:             findingID,
		ProjectID:             projectID,
		AssumeRoleArn:         ds.AssumeRoleArn,
		ExternalID:            ds.ExternalID,
	}
	resp, err := a.sqs.Send(ctx, a.remediationProposalQueueURL, msg)
	if err != nil {
		return err
	}
	a.logger.Infof(ctx, "Generated remediation proposal: remediation_proposal_id=%d, finding_id=%d, messageId=%v", remediationProposalID, findingID, resp.MessageId)
	return nil
}

func isRemediationProposalTarget(dataSource string) bool {
	for _, target := range remediationProposalTargetDataSources {
		if dataSource == target {
			return true
		}
	}
	return false
}

func isAWSAccountIDInAssumeRoleArn(accountID, assumeRoleArn string) bool {
	parts := strings.SplitN(assumeRoleArn, ":", 6)
	if len(parts) != 6 {
		return false
	}
	return parts[0] == "arn" && parts[2] == "iam" && parts[4] == accountID
}

func (a *AIService) getAWSAccountIDFromFindingTag(ctx context.Context, projectID uint32, findingID uint64) (string, error) {
	tags, err := a.findingClient.ListFindingTag(ctx, &finding.ListFindingTagRequest{
		ProjectId: projectID,
		FindingId: findingID,
		Limit:     listFindingTagLimit,
	})
	if err != nil {
		return "", err
	}
	for _, t := range tags.Tag {
		if awsAccountIDPattern.MatchString(t.Tag) {
			return t.Tag, nil
		}
	}
	a.logger.Warnf(ctx, "AWS account_id tag is not found for remediation proposal: project_id=%d, finding_id=%d", projectID, findingID)
	return "", status.Error(codes.NotFound, "aws remediation target is not configured")
}
