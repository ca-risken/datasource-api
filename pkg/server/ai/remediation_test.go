package ai

import (
	"context"
	"errors"
	"testing"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	coreai "github.com/ca-risken/core/proto/ai"
	coreaimocks "github.com/ca-risken/core/proto/ai/mocks"
	"github.com/ca-risken/core/proto/finding"
	findingmocks "github.com/ca-risken/core/proto/finding/mocks"
	"github.com/ca-risken/datasource-api/pkg/db"
	dbmocks "github.com/ca-risken/datasource-api/pkg/db/mocks"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	aipb "github.com/ca-risken/datasource-api/proto/datasource_ai"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type mockDBClient struct {
	*dbmocks.AWSRepoInterface
}

type mockSQS struct {
	mock.Mock
}

func (m *mockSQS) Send(ctx context.Context, url string, msg interface{}) (*awssqs.SendMessageOutput, error) {
	args := m.Called(ctx, url, msg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*awssqs.SendMessageOutput), args.Error(1)
}

func TestGenerateRemediationProposal(t *testing.T) {
	targetFinding := &finding.GetFindingResponse{
		Finding: &finding.Finding{FindingId: 1001, ProjectId: 1, DataSource: "aws:cloudsploit"},
	}
	accountTags := &finding.ListFindingTagResponse{
		Tag: []*finding.FindingTag{
			{Tag: "aws"},
			{Tag: "cloudsploit"},
			{Tag: "123456789012"},
		},
	}
	awsData := &model.AWS{AWSID: 5, ProjectID: 1, AWSAccountID: "123456789012"}
	awsDataSources := &[]db.DataSource{
		{AWSDataSourceID: 1003, DataSource: "aws:cloudsploit", AWSID: 5, ProjectID: 1},
	}
	dsForMessage := &db.DataSource{
		AWSDataSourceID: 1003,
		DataSource:      "aws:cloudsploit",
		AWSID:           5,
		ProjectID:       1,
		AWSAccountID:    "123456789012",
		AssumeRoleArn:   "arn:aws:iam::123456789012:role/risken",
		ExternalID:      "ext-id",
	}
	createdProposal := &coreai.CreateRemediationProposalResponse{
		RemediationProposal: &coreai.RemediationProposal{RemediationProposalId: 2001},
	}

	cases := []struct {
		name      string
		input     *aipb.GenerateRemediationProposalRequest
		mockSetup func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name:  "OK",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(targetFinding, nil).Once()
				f.On("ListFindingTag", mock.Anything, mock.Anything).Return(accountTags, nil).Once()
				awsRepo.On("GetAWSByAccountID", mock.Anything, uint32(1), "123456789012").Return(awsData, nil).Once()
				awsRepo.On("ListAWSDataSource", mock.Anything, uint32(1), uint32(5), "aws:cloudsploit").Return(awsDataSources, nil).Once()
				awsRepo.On("GetAWSDataSourceForMessage", mock.Anything, uint32(5), uint32(1003), uint32(1)).Return(dsForMessage, nil).Once()
				a.On("CreateRemediationProposal", mock.Anything, mock.Anything).Return(createdProposal, nil).Once()
				s.On("Send", mock.Anything, "https://example.com/queue/aws-remediation-proposal", mock.MatchedBy(func(msg interface{}) bool {
					m, ok := msg.(*message.RemediationProposalQueueMessage)
					return ok &&
						m.RemediationProposalID == 2001 &&
						m.FindingID == 1001 &&
						m.ProjectID == 1 &&
						m.AssumeRoleArn == "arn:aws:iam::123456789012:role/risken" &&
						m.ExternalID == "ext-id"
				})).Return(&awssqs.SendMessageOutput{}, nil).Once()
			},
		},
		{
			name:  "NG validation error",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 0, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
			},
			wantErr: true,
		},
		{
			name:  "NG finding not found",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 9999},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(&finding.GetFindingResponse{}, nil).Once()
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name:  "NG unsupported data_source",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(&finding.GetFindingResponse{
					Finding: &finding.Finding{FindingId: 1001, ProjectId: 1, DataSource: "aws:guard-duty"},
				}, nil).Once()
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name:  "NG account_id tag not found",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(targetFinding, nil).Once()
				f.On("ListFindingTag", mock.Anything, mock.Anything).Return(&finding.ListFindingTagResponse{
					Tag: []*finding.FindingTag{{Tag: "aws"}},
				}, nil).Once()
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name:  "NG aws account not registered",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(targetFinding, nil).Once()
				f.On("ListFindingTag", mock.Anything, mock.Anything).Return(accountTags, nil).Once()
				awsRepo.On("GetAWSByAccountID", mock.Anything, uint32(1), "123456789012").Return(nil, gorm.ErrRecordNotFound).Once()
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name:  "NG aws data_source not attached",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(targetFinding, nil).Once()
				f.On("ListFindingTag", mock.Anything, mock.Anything).Return(accountTags, nil).Once()
				awsRepo.On("GetAWSByAccountID", mock.Anything, uint32(1), "123456789012").Return(awsData, nil).Once()
				awsRepo.On("ListAWSDataSource", mock.Anything, uint32(1), uint32(5), "aws:cloudsploit").Return(&[]db.DataSource{}, nil).Once()
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name:  "NG create remediation proposal error",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(targetFinding, nil).Once()
				f.On("ListFindingTag", mock.Anything, mock.Anything).Return(accountTags, nil).Once()
				awsRepo.On("GetAWSByAccountID", mock.Anything, uint32(1), "123456789012").Return(awsData, nil).Once()
				awsRepo.On("ListAWSDataSource", mock.Anything, uint32(1), uint32(5), "aws:cloudsploit").Return(awsDataSources, nil).Once()
				awsRepo.On("GetAWSDataSourceForMessage", mock.Anything, uint32(5), uint32(1003), uint32(1)).Return(dsForMessage, nil).Once()
				a.On("CreateRemediationProposal", mock.Anything, mock.Anything).Return(nil, errors.New("core error")).Once()
			},
			wantErr: true,
		},
		{
			name:  "NG core returned empty proposal",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(targetFinding, nil).Once()
				f.On("ListFindingTag", mock.Anything, mock.Anything).Return(accountTags, nil).Once()
				awsRepo.On("GetAWSByAccountID", mock.Anything, uint32(1), "123456789012").Return(awsData, nil).Once()
				awsRepo.On("ListAWSDataSource", mock.Anything, uint32(1), uint32(5), "aws:cloudsploit").Return(awsDataSources, nil).Once()
				awsRepo.On("GetAWSDataSourceForMessage", mock.Anything, uint32(5), uint32(1003), uint32(1)).Return(dsForMessage, nil).Once()
				a.On("CreateRemediationProposal", mock.Anything, mock.Anything).Return(&coreai.CreateRemediationProposalResponse{}, nil).Once()
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name:  "NG sqs send error",
			input: &aipb.GenerateRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			mockSetup: func(f *findingmocks.FindingServiceClient, a *coreaimocks.AIServiceClient, awsRepo *dbmocks.AWSRepoInterface, s *mockSQS) {
				f.On("GetFinding", mock.Anything, mock.Anything).Return(targetFinding, nil).Once()
				f.On("ListFindingTag", mock.Anything, mock.Anything).Return(accountTags, nil).Once()
				awsRepo.On("GetAWSByAccountID", mock.Anything, uint32(1), "123456789012").Return(awsData, nil).Once()
				awsRepo.On("ListAWSDataSource", mock.Anything, uint32(1), uint32(5), "aws:cloudsploit").Return(awsDataSources, nil).Once()
				awsRepo.On("GetAWSDataSourceForMessage", mock.Anything, uint32(5), uint32(1003), uint32(1)).Return(dsForMessage, nil).Once()
				a.On("CreateRemediationProposal", mock.Anything, mock.Anything).Return(createdProposal, nil).Once()
				s.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("sqs error")).Once()
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			findingMock := findingmocks.NewFindingServiceClient(t)
			coreAIMock := coreaimocks.NewAIServiceClient(t)
			awsRepoMock := dbmocks.NewAWSRepoInterface(t)
			sqsMock := &mockSQS{}
			svc := AIService{
				dbClient:                    &mockDBClient{AWSRepoInterface: awsRepoMock},
				findingClient:               findingMock,
				coreAIClient:                coreAIMock,
				sqs:                         sqsMock,
				remediationProposalQueueURL: "https://example.com/queue/aws-remediation-proposal",
				logger:                      logging.NewLogger(),
			}
			c.mockSetup(findingMock, coreAIMock, awsRepoMock, sqsMock)

			result, err := svc.GenerateRemediationProposal(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if c.wantCode != codes.OK {
				if st, ok := status.FromError(err); !ok || st.Code() != c.wantCode {
					t.Fatalf("unexpected error code: want=%v, got=%v (err=%+v)", c.wantCode, st.Code(), err)
				}
			}
			if !c.wantErr && result.RemediationProposalId == 0 {
				t.Fatal("remediation_proposal_id is empty")
			}
			sqsMock.AssertExpectations(t)
		})
	}
}
