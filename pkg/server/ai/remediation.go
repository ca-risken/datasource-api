package ai

import (
	"context"
	"errors"
	"regexp"
	"time"

	coreai "github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	aipb "github.com/ca-risken/datasource-api/proto/datasource_ai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	remediationProposalStatusPending   = "PENDING"
	remediationProposalStatusSucceeded = "SUCCEEDED"

	aiRemediationCooldownPeriod = time.Hour
	listFindingTagLimit         = 200
)

var aiRemediationTargetDataSources = []string{
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
	if !isAIRemediationTarget(dataSource) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported data_source for AI remediation: %s", dataSource)
	}

	proposals, err := a.coreAIClient.ListRemediationProposal(ctx, &coreai.ListRemediationProposalRequest{
		ProjectId: req.ProjectId,
		FindingId: req.FindingId,
		Status:    []string{remediationProposalStatusPending, remediationProposalStatusSucceeded},
	})
	if err != nil {
		return nil, err
	}
	for _, p := range proposals.RemediationProposal {
		if time.Since(time.Unix(p.CreatedAt, 0)) < aiRemediationCooldownPeriod {
			return nil, status.Errorf(codes.FailedPrecondition,
				"remediation proposal is in cooldown period: finding_id=%d, remediation_proposal_id=%d", req.FindingId, p.RemediationProposalId)
		}
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

	putResp, err := a.coreAIClient.PutRemediationProposal(ctx, &coreai.PutRemediationProposalRequest{
		ProjectId: req.ProjectId,
		FindingId: req.FindingId,
		Status:    remediationProposalStatusPending,
	})
	if err != nil {
		return nil, err
	}
	if putResp.RemediationProposal == nil || putResp.RemediationProposal.RemediationProposalId == 0 {
		return nil, status.Error(codes.Internal, "failed to create remediation proposal")
	}
	remediationProposalID := putResp.RemediationProposal.RemediationProposalId

	msg := &message.AIRemediationQueueMessage{
		RemediationProposalID: remediationProposalID,
		FindingID:             req.FindingId,
		ProjectID:             req.ProjectId,
		DataSource:            ds.DataSource,
		AWSID:                 ds.AWSID,
		AccountID:             ds.AWSAccountID,
		AssumeRoleArn:         ds.AssumeRoleArn,
		ExternalID:            ds.ExternalID,
	}
	resp, err := a.sqs.Send(ctx, a.aiRemediationQueueURL, msg)
	if err != nil {
		return nil, err
	}
	a.logger.Infof(ctx, "Generated remediation proposal: remediation_proposal_id=%d, finding_id=%d, messageId=%v", remediationProposalID, req.FindingId, resp.MessageId)
	return &aipb.GenerateRemediationProposalResponse{RemediationProposalId: remediationProposalID}, nil
}

func isAIRemediationTarget(dataSource string) bool {
	for _, target := range aiRemediationTargetDataSources {
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
