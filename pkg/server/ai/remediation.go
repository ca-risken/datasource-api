package ai

import (
	"context"
	"errors"
	"regexp"

	coreai "github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	aipb "github.com/ca-risken/datasource-api/proto/datasource_ai"
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

func (a *AIService) GenerateRemediationProposal(ctx context.Context, req *aipb.GenerateRemediationProposalRequest) (*aipb.GenerateRemediationProposalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	findingResp, err := a.findingClient.GetFinding(ctx, &finding.GetFindingRequest{
		ProjectId: req.ProjectId,
		FindingId: req.FindingId,
	})
	if err != nil {
		return nil, err
	}
	if findingResp.Finding == nil {
		return nil, status.Errorf(codes.NotFound, "finding not found: finding_id=%d", req.FindingId)
	}

	dataSource := findingResp.Finding.DataSource
	if !isRemediationProposalTarget(dataSource) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported data_source for remediation proposal: %s", dataSource)
	}

	accountID, err := a.getAWSAccountIDFromFindingTag(ctx, req.ProjectId, req.FindingId)
	if err != nil {
		return nil, err
	}
	awsData, err := a.dbClient.GetAWSByAccountID(ctx, req.ProjectId, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "aws account is not registered: account_id=%s", accountID)
		}
		return nil, err
	}
	awsDataSources, err := a.dbClient.ListAWSDataSource(ctx, req.ProjectId, awsData.AWSID, dataSource)
	if err != nil {
		return nil, err
	}
	if awsDataSources == nil || len(*awsDataSources) == 0 {
		return nil, status.Errorf(codes.NotFound, "aws data_source not found: data_source=%s", dataSource)
	}
	ds, err := a.dbClient.GetAWSDataSourceForMessage(ctx, awsData.AWSID, (*awsDataSources)[0].AWSDataSourceID, req.ProjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "aws data_source is not attached: aws_id=%d, data_source=%s", awsData.AWSID, dataSource)
		}
		return nil, err
	}

	createResp, err := a.coreAIClient.CreateRemediationProposal(ctx, &coreai.CreateRemediationProposalRequest{
		ProjectId: req.ProjectId,
		FindingId: req.FindingId,
	})
	if err != nil {
		return nil, err
	}
	if createResp.RemediationProposal == nil || createResp.RemediationProposal.RemediationProposalId == 0 {
		return nil, status.Error(codes.Internal, "failed to create remediation proposal")
	}
	remediationProposalID := createResp.RemediationProposal.RemediationProposalId

	msg := &message.RemediationProposalQueueMessage{
		RemediationProposalID: remediationProposalID,
		FindingID:             req.FindingId,
		ProjectID:             req.ProjectId,
		AssumeRoleArn:         ds.AssumeRoleArn,
		ExternalID:            ds.ExternalID,
	}
	resp, err := a.sqs.Send(ctx, a.remediationProposalQueueURL, msg)
	if err != nil {
		return nil, err
	}
	a.logger.Infof(ctx, "Generated remediation proposal: remediation_proposal_id=%d, finding_id=%d, messageId=%v", remediationProposalID, req.FindingId, resp.MessageId)
	return &aipb.GenerateRemediationProposalResponse{RemediationProposalId: remediationProposalID}, nil
}

func isRemediationProposalTarget(dataSource string) bool {
	for _, target := range remediationProposalTargetDataSources {
		if dataSource == target {
			return true
		}
	}
	return false
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
	return "", status.Errorf(codes.NotFound, "aws account_id tag not found: finding_id=%d", findingID)
}
