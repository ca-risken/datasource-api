package code

import (
	"context"
	"crypto/aes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	projectmock "github.com/ca-risken/core/proto/project/mocks"
	"github.com/ca-risken/datasource-api/pkg/db/mocks"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/pkg/test"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/golang/protobuf/ptypes/empty"
	ghub "github.com/google/go-github/v44/github"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	length65string = "12345678901234567890123456789012345678901234567890123456789012345"
)

func TestListDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *code.ListDataSourceRequest
		want         *code.ListDataSourceResponse
		mockResponse *[]model.CodeDataSource
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &code.ListDataSourceRequest{CodeDataSourceId: 1},
			want: &code.ListDataSourceResponse{CodeDataSource: []*code.CodeDataSource{
				{Name: "one", Description: "desc", MaxScore: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{CodeDataSourceId: 2, Name: "two", Description: "desc", MaxScore: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &[]model.CodeDataSource{
				{Name: "one", Description: "desc", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
				{CodeDataSourceID: 2, Name: "two", Description: "desc", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty",
			input:     &code.ListDataSourceRequest{Name: "not exists name"},
			want:      &code.ListDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.ListDataSourceRequest{Name: length65string},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.ListDataSourceRequest{CodeDataSourceId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListCodeDataSource", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.ListDataSource(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestListGitHubSetting(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                   string
		input                  *code.ListGitHubSettingRequest
		want                   *code.ListGitHubSettingResponse
		mockResponse           *[]model.CodeGitHubSetting
		mockGitleaksResponse   *[]model.CodeGitleaksSetting
		mockDependencyResponse *[]model.CodeDependencySetting
		mockCodeScanResponse   *[]model.CodeCodeScanSetting
		mockError              error
		mockGitleaksError      error
		mockDependencyError    error
		mockCodeScanError      error
		wantErr                bool
	}{
		{
			name:  "OK",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			want: &code.ListGitHubSettingResponse{GithubSetting: []*code.GitHubSetting{
				{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					CodeScanSetting:   &code.CodeScanSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
				{GithubSettingId: 2, Name: "two", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 2, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo2", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					DependencySetting: &code.DependencySetting{GithubSettingId: 2, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					CodeScanSetting:   &code.CodeScanSetting{GithubSettingId: 2, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			}},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, Name: "two", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo2", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockCodeScanResponse: &[]model.CodeCodeScanSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:         "OK empty",
			input:        &code.ListGitHubSettingRequest{ProjectId: 1},
			want:         &code.ListGitHubSettingResponse{},
			mockResponse: &[]model.CodeGitHubSetting{},
		},
		{
			name:  "OK gitleaks setting empty",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			want: &code.ListGitHubSettingResponse{GithubSetting: []*code.GitHubSetting{
				{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					CodeScanSetting:   &code.CodeScanSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				}}},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now}},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockCodeScanResponse: &[]model.CodeCodeScanSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:  "OK dependency setting empty",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			want: &code.ListGitHubSettingResponse{GithubSetting: []*code.GitHubSetting{
				{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					GitleaksSetting: &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					CodeScanSetting: &code.CodeScanSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				}}},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now}},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockDependencyResponse: &[]model.CodeDependencySetting{},
			mockCodeScanResponse: &[]model.CodeCodeScanSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:  "OK code scan setting empty",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			want: &code.ListGitHubSettingResponse{GithubSetting: []*code.GitHubSetting{
				{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				}}},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now}},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockCodeScanResponse: &[]model.CodeCodeScanSetting{},
		},
		{
			name:    "NG invalid param",
			input:   &code.ListGitHubSettingRequest{},
			wantErr: true,
		},
		{
			name:      "Invalid DB error when listGitHubSetting",
			input:     &code.ListGitHubSettingRequest{ProjectId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
		{
			name:  "Invalid DB error when getGitleaksSetting",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			},
			mockGitleaksResponse: nil,
			mockGitleaksError:    gorm.ErrInvalidDB,
			wantErr:              true,
		},
		{
			name:  "Invalid DB error when getDependencySetting",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockDependencyError:  gorm.ErrInvalidDB,
			wantErr:              true,
		},
		{
			name:  "Invalid DB error when getCodeScanSetting",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			},
			mockGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockDependencyResponse: &[]model.CodeDependencySetting{},
			mockCodeScanError:      gorm.ErrInvalidDB,
			wantErr:                true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListGitHubSetting", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			if c.mockGitleaksResponse != nil || c.mockGitleaksError != nil {
				mockDB.On("ListGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockGitleaksResponse, c.mockGitleaksError).Once()
			}
			if c.mockDependencyResponse != nil || c.mockDependencyError != nil {
				mockDB.On("ListDependencySetting", test.RepeatMockAnything(3)...).Return(c.mockDependencyResponse, c.mockDependencyError).Once()
			}
			if c.mockCodeScanResponse != nil || c.mockCodeScanError != nil {
				mockDB.On("ListCodeScanSetting", test.RepeatMockAnything(3)...).Return(c.mockCodeScanResponse, c.mockCodeScanError).Once()
			}

			got, err := svc.ListGitHubSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetGitHubSetting(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                   string
		input                  *code.GetGitHubSettingRequest
		want                   *code.GetGitHubSettingResponse
		mockResponse           *model.CodeGitHubSetting
		mockError              error
		mockGitleaksResponse   *model.CodeGitleaksSetting
		mockGitleaksError      error
		mockDependencyResponse *model.CodeDependencySetting
		mockDependencyError    error
		mockCodeScanResponse   *model.CodeCodeScanSetting
		mockCodeScanError      error
		wantErr                bool
	}{
		{
			name:  "OK",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			want: &code.GetGitHubSettingResponse{GithubSetting: &code.GitHubSetting{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				CodeScanSetting:   &code.CodeScanSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockCodeScanResponse: &model.CodeCodeScanSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:      "OK empty",
			input:     &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			want:      &code.GetGitHubSettingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:  "OK gitleaks setting empty",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			want: &code.GetGitHubSettingResponse{GithubSetting: &code.GitHubSetting{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				CodeScanSetting:   &code.CodeScanSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			mockGitleaksError: gorm.ErrRecordNotFound,
			mockDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockCodeScanResponse: &model.CodeCodeScanSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:  "OK dependency setting empty",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			want: &code.GetGitHubSettingResponse{GithubSetting: &code.GitHubSetting{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				GitleaksSetting: &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				CodeScanSetting: &code.CodeScanSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyError: gorm.ErrRecordNotFound,
			mockCodeScanResponse: &model.CodeCodeScanSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:  "OK code scan setting empty",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			want: &code.GetGitHubSettingResponse{GithubSetting: &code.GitHubSetting{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockCodeScanError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.GetGitHubSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
		{
			name:  "Invalid DB error when GetGitleaksSetting",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
			mockGitleaksError: gorm.ErrInvalidDB,
			wantErr:           true,
		},
		{
			name:  "Invalid DB error when GetDependencySetting",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyError: gorm.ErrInvalidDB,
			wantErr:             true,
		},
		{
			name:  "Invalid DB error when GetCodeScanSetting",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockCodeScanError: gorm.ErrInvalidDB,
			wantErr:           true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("GetGitHubSetting", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			if c.mockGitleaksResponse != nil || c.mockGitleaksError != nil {
				mockDB.On("GetGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockGitleaksResponse, c.mockGitleaksError).Once()
			}
			if c.mockDependencyResponse != nil || c.mockDependencyError != nil {
				mockDB.On("GetDependencySetting", test.RepeatMockAnything(3)...).Return(c.mockDependencyResponse, c.mockDependencyError).Once()
			}
			if c.mockCodeScanResponse != nil || c.mockCodeScanError != nil {
				mockDB.On("GetCodeScanSetting", test.RepeatMockAnything(3)...).Return(c.mockCodeScanResponse, c.mockCodeScanError).Once()
			}

			got, err := svc.GetGitHubSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutGitHubSetting(t *testing.T) {
	now := time.Now()
	key := []byte("1234567890123456")
	block, err := aes.NewCipher(key)
	assert.NoError(t, err)
	cases := []struct {
		name         string
		input        *code.PutGitHubSettingRequest
		want         *code.PutGitHubSettingResponse
		mockResponse *model.CodeGitHubSetting
		mockError    error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &code.PutGitHubSettingRequest{ProjectId: 1, GithubSetting: &code.GitHubSettingForUpsert{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData},
			},
			want: &code.PutGitHubSettingResponse{GithubSetting: &code.GitHubSetting{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name: "OK(empty)",
			input: &code.PutGitHubSettingRequest{ProjectId: 1, GithubSetting: &code.GitHubSettingForUpsert{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData},
			},
			want: &code.PutGitHubSettingResponse{GithubSetting: &code.GitHubSetting{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:    "NG invalid param",
			input:   &code.PutGitHubSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &code.PutGitHubSettingRequest{ProjectId: 1, GithubSetting: &code.GitHubSettingForUpsert{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData},
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB, cipherBlock: block}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("UpsertGitHubSetting", test.RepeatMockAnything(2)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.PutGitHubSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteGitHubSetting(t *testing.T) {
	cases := []struct {
		name    string
		input   *code.DeleteGitHubSettingRequest
		wantErr bool

		callDeleteGitleaksCache     bool
		mockDeleteGitleaksCacheResp error

		callDeleteGitleaks     bool
		mockDeleteGitleaksResp error

		callDeleteDependency     bool
		mockDeleteDependencyResp error

		callDeleteCodeScan     bool
		mockDeleteCodeScanResp error

		callDeleteGithubSetting     bool
		mockDeleteGithubSettingResp error
	}{
		{
			name:    "OK",
			input:   &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: false,

			callDeleteGitleaksCache:     true,
			mockDeleteGitleaksCacheResp: nil,

			callDeleteGitleaks:     true,
			mockDeleteGitleaksResp: nil,

			callDeleteDependency:     true,
			mockDeleteDependencyResp: nil,

			callDeleteCodeScan:     true,
			mockDeleteCodeScanResp: nil,

			callDeleteGithubSetting:     true,
			mockDeleteGithubSettingResp: nil,
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteGitHubSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error (delete gitleaks cache)",
			input:   &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: true,

			callDeleteGitleaksCache:     true,
			mockDeleteGitleaksCacheResp: errors.New("something error"),
		},
		{
			name:    "NG DB error (delete gitleaks)",
			input:   &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: true,

			callDeleteGitleaksCache:     true,
			mockDeleteGitleaksCacheResp: nil,

			callDeleteGitleaks:     true,
			mockDeleteGitleaksResp: errors.New("something error"),
		},
		{
			name:    "NG DB error (delete dependency)",
			input:   &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: true,

			callDeleteGitleaksCache:     true,
			mockDeleteGitleaksCacheResp: nil,

			callDeleteGitleaks:     true,
			mockDeleteGitleaksResp: nil,

			callDeleteDependency:     true,
			mockDeleteDependencyResp: errors.New("something error"),
		},
		{
			name:    "NG DB error (delete code scan setting)",
			input:   &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: true,

			callDeleteGitleaksCache:     true,
			mockDeleteGitleaksCacheResp: nil,

			callDeleteGitleaks:     true,
			mockDeleteGitleaksResp: nil,

			callDeleteDependency:     true,
			mockDeleteDependencyResp: nil,

			callDeleteCodeScan:     true,
			mockDeleteCodeScanResp: errors.New("something error"),
		},
		{
			name:    "NG DB error (delete github setting)",
			input:   &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: true,

			callDeleteGitleaksCache:     true,
			mockDeleteGitleaksCacheResp: nil,

			callDeleteGitleaks:     true,
			mockDeleteGitleaksResp: nil,

			callDeleteDependency:     true,
			mockDeleteDependencyResp: nil,

			callDeleteCodeScan:     true,
			mockDeleteCodeScanResp: nil,

			callDeleteGithubSetting:     true,
			mockDeleteGithubSettingResp: errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.callDeleteGitleaksCache {
				mockDB.On("DeleteGitleaksCache", test.RepeatMockAnything(3)...).Return(c.mockDeleteGitleaksCacheResp).Once()
			}
			if c.callDeleteGitleaks {
				mockDB.On("DeleteGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockDeleteGitleaksResp).Once()
			}
			if c.callDeleteDependency {
				mockDB.On("DeleteDependencySetting", test.RepeatMockAnything(3)...).Return(c.mockDeleteDependencyResp).Once()
			}
			if c.callDeleteCodeScan {
				mockDB.On("DeleteCodeScanSetting", test.RepeatMockAnything(3)...).Return(c.mockDeleteCodeScanResp).Once()
			}
			if c.callDeleteGithubSetting {
				mockDB.On("DeleteGitHubSetting", test.RepeatMockAnything(3)...).Return(c.mockDeleteGithubSettingResp).Once()
			}
			_, err := svc.DeleteGitHubSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestPutGitleaksSetting(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *code.PutGitleaksSettingRequest
		want         *code.PutGitleaksSettingResponse
		mockResponse *model.CodeGitleaksSetting
		mockError    error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &code.PutGitleaksSettingRequest{ProjectId: 1, GitleaksSetting: &code.GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want: &code.PutGitleaksSettingResponse{GitleaksSetting: &code.GitleaksSetting{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name: "OK(empty)",
			input: &code.PutGitleaksSettingRequest{ProjectId: 1, GitleaksSetting: &code.GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want: &code.PutGitleaksSettingResponse{GitleaksSetting: &code.GitleaksSetting{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:    "NG invalid param",
			input:   &code.PutGitleaksSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &code.PutGitleaksSettingRequest{ProjectId: 1, GitleaksSetting: &code.GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("UpsertGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.PutGitleaksSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteGitleaksSetting(t *testing.T) {
	cases := []struct {
		name    string
		input   *code.DeleteGitleaksSettingRequest
		wantErr bool

		callDeleteGitleaksRepository     bool
		mockDeleteGitleaksRepositoryResp error

		callDeleteGitleaksCache     bool
		mockDeleteGitleaksCacheResp error

		callDeleteGitleaksSetting     bool
		mockDeleteGitleaksSettingResp error
	}{
		{
			name:                         "OK",
			input:                        &code.DeleteGitleaksSettingRequest{ProjectId: 1, GithubSettingId: 1},
			callDeleteGitleaksRepository: true,
			callDeleteGitleaksCache:      true,
			callDeleteGitleaksSetting:    true,
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteGitleaksSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:                         "NG(DeleteGitleaksCache error)",
			input:                        &code.DeleteGitleaksSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr:                      true,
			callDeleteGitleaksRepository: true,
			callDeleteGitleaksCache:      true,
			mockDeleteGitleaksCacheResp:  gorm.ErrInvalidDB,
		},
		{
			name:                          "NG(DeleteGitleaksSetting error)",
			input:                         &code.DeleteGitleaksSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr:                       true,
			callDeleteGitleaksRepository:  true,
			callDeleteGitleaksCache:       true,
			callDeleteGitleaksSetting:     true,
			mockDeleteGitleaksSettingResp: gorm.ErrInvalidDB,
		},
		{
			name:                             "NG(DeleteGitleaksRepository error)",
			input:                            &code.DeleteGitleaksSettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr:                          true,
			callDeleteGitleaksRepository:     true,
			mockDeleteGitleaksRepositoryResp: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.callDeleteGitleaksRepository {
				mockDB.On("DeleteGitleaksRepository", test.RepeatMockAnything(3)...).Return(c.mockDeleteGitleaksRepositoryResp).Once()
			}
			if c.callDeleteGitleaksCache {
				mockDB.On("DeleteGitleaksCache", test.RepeatMockAnything(2)...).Return(c.mockDeleteGitleaksCacheResp).Once()
			}
			if c.callDeleteGitleaksSetting {
				mockDB.On("DeleteGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockDeleteGitleaksSettingResp).Once()
			}
			_, err := svc.DeleteGitleaksSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestListGitleaksCache(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *code.ListGitleaksCacheRequest
		want    *code.ListGitleaksCacheResponse
		wantErr bool

		mockResp *[]model.CodeGitleaksCache
		mockErr  error
	}{
		{
			name: "OK",
			input: &code.ListGitleaksCacheRequest{
				ProjectId: 1, GithubSettingId: 1,
			},
			want: &code.ListGitleaksCacheResponse{
				GitleaksCache: []*code.GitleaksCache{
					{GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					{GithubSettingId: 2, RepositoryFullName: "owner/repo2", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResp: &[]model.CodeGitleaksCache{
				{CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, RepositoryFullName: "owner/repo2", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
		},
		{
			name: "NG(invalid param)",
			input: &code.ListGitleaksCacheRequest{
				GithubSettingId: 1,
			},
			wantErr: true,
		},
		{
			name: "NG(DB error)",
			input: &code.ListGitleaksCacheRequest{
				ProjectId: 1, GithubSettingId: 1,
			},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListGitleaksCache", test.RepeatMockAnything(3)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListGitleaksCache(context.TODO(), c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetGitleaksCache(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *code.GetGitleaksCacheRequest
		want    *code.GetGitleaksCacheResponse
		wantErr bool

		mockResp *model.CodeGitleaksCache
		mockErr  error
	}{
		{
			name: "OK",
			input: &code.GetGitleaksCacheRequest{
				ProjectId: 1, GithubSettingId: 1, RepositoryFullName: "owner/repo",
			},
			want: &code.GetGitleaksCacheResponse{
				GitleaksCache: &code.GitleaksCache{GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResp: &model.CodeGitleaksCache{CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr:  false,
		},
		{
			name: "OK(RecordNotFound)",
			input: &code.GetGitleaksCacheRequest{
				ProjectId: 1, GithubSettingId: 1, RepositoryFullName: "owner/repo",
			},
			want:    &code.GetGitleaksCacheResponse{},
			mockErr: gorm.ErrRecordNotFound,
			wantErr: false,
		},
		{
			name: "NG(invalid param)",
			input: &code.GetGitleaksCacheRequest{
				GithubSettingId: 1, RepositoryFullName: "owner/repo",
			},
			wantErr: true,
		},
		{
			name: "NG(DB error)",
			input: &code.GetGitleaksCacheRequest{
				ProjectId: 1, GithubSettingId: 1, RepositoryFullName: "owner/repo",
			},
			mockErr: gorm.ErrInvalidDB,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("GetGitleaksCache", test.RepeatMockAnything(5)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.GetGitleaksCache(context.TODO(), c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutGitleaksCache(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *code.PutGitleaksCacheRequest
		want    *code.PutGitleaksCacheResponse
		wantErr bool

		mockGetGitleaksSettingResp *model.CodeGitleaksSetting
		mockGetGitleaksSettingErr  error

		mockUpsertGitleaksCacheResp *model.CodeGitleaksCache
		mockUpsertGitleaksCacheErr  error
	}{
		{
			name: "OK",
			input: &code.PutGitleaksCacheRequest{
				ProjectId: 1,
				GitleaksCache: &code.GitleaksCacheForUpsert{
					GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix(),
				},
			},
			want: &code.PutGitleaksCacheResponse{
				GitleaksCache: &code.GitleaksCache{GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			wantErr:                     false,
			mockGetGitleaksSettingResp:  &model.CodeGitleaksSetting{CodeGitHubSettingID: 1},
			mockUpsertGitleaksCacheResp: &model.CodeGitleaksCache{CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", ScanAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name: "NG(invalid param)",
			input: &code.PutGitleaksCacheRequest{
				// ProjectId: 1, // required param
				GitleaksCache: &code.GitleaksCacheForUpsert{
					GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NG(No GitHub setting)",
			input: &code.PutGitleaksCacheRequest{
				ProjectId: 1,
				GitleaksCache: &code.GitleaksCacheForUpsert{
					GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix(),
				},
			},
			want:                      nil,
			wantErr:                   true,
			mockGetGitleaksSettingErr: gorm.ErrRecordNotFound,
		},
		{
			name: "NG(PutGitleaksCache error)",
			input: &code.PutGitleaksCacheRequest{
				ProjectId: 1,
				GitleaksCache: &code.GitleaksCacheForUpsert{
					GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix(),
				},
			},
			want:                       nil,
			wantErr:                    true,
			mockGetGitleaksSettingResp: &model.CodeGitleaksSetting{CodeGitHubSettingID: 1},
			mockUpsertGitleaksCacheErr: errors.New("something error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}
			if c.mockGetGitleaksSettingResp != nil || c.mockGetGitleaksSettingErr != nil {
				mockDB.On("GetGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockGetGitleaksSettingResp, c.mockGetGitleaksSettingErr).Once()
			}
			if c.mockUpsertGitleaksCacheResp != nil || c.mockUpsertGitleaksCacheErr != nil {
				mockDB.On("UpsertGitleaksCache", test.RepeatMockAnything(3)...).Return(c.mockUpsertGitleaksCacheResp, c.mockUpsertGitleaksCacheErr).Once()
			}
			got, err := svc.PutGitleaksCache(context.TODO(), c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutDependencySetting(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *code.PutDependencySettingRequest
		want         *code.PutDependencySettingResponse
		mockResponse *model.CodeDependencySetting
		mockError    error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &code.PutDependencySettingRequest{ProjectId: 1, DependencySetting: &code.DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want: &code.PutDependencySettingResponse{DependencySetting: &code.DependencySetting{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name: "OK(empty)",
			input: &code.PutDependencySettingRequest{ProjectId: 1, DependencySetting: &code.DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want: &code.PutDependencySettingResponse{DependencySetting: &code.DependencySetting{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:    "NG invalid param",
			input:   &code.PutDependencySettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &code.PutDependencySettingRequest{ProjectId: 1, DependencySetting: &code.DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("UpsertDependencySetting", test.RepeatMockAnything(2)...).Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.PutDependencySetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteDependencySetting(t *testing.T) {
	cases := []struct {
		name    string
		input   *code.DeleteDependencySettingRequest
		wantErr bool

		mockCall  bool
		mockError error
	}{
		{
			name:     "OK",
			input:    &code.DeleteDependencySettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockCall: true,
		},
		{
			name:     "NG invalid param",
			input:    &code.DeleteDependencySettingRequest{ProjectId: 1},
			wantErr:  true,
			mockCall: false,
		},
		{
			name:      "Invalid DB error",
			input:     &code.DeleteDependencySettingRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr:   true,
			mockCall:  true,
			mockError: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB}

			if c.mockCall {
				// First, repositories are deleted in bulk
				mockDB.On("DeleteDependencyRepository", test.RepeatMockAnything(3)...).Return(c.mockError).Once()
				// If the first call errors, the second should not be invoked; set it only when no error expected
				if c.mockError == nil {
					mockDB.On("DeleteDependencySetting", test.RepeatMockAnything(3)...).Return(c.mockError).Once()
				}
			}
			_, err := svc.DeleteDependencySetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestInvokeScan(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                       string
		input                      *code.InvokeScanGitleaksRequest
		mockGetGitleaksResponse    *model.CodeGitleaksSetting
		mockGetGitleaksError       error
		mockQueue                  CodeQueue
		mockUpsertGitleaksResponse *model.CodeGitleaksSetting
		mockUpsertGitleaksError    error
		wantErr                    bool
	}{
		{
			name:  "OK",
			input: &code.InvokeScanGitleaksRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockQueue:                  newFakeCodeQueue("succeed", nil),
			mockUpsertGitleaksResponse: &model.CodeGitleaksSetting{},
		},
		{
			name:    "NG invalid param",
			input:   &code.InvokeScanGitleaksRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:                 "NG db error when GetGitHubSetting",
			input:                &code.InvokeScanGitleaksRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetGitleaksError: gorm.ErrRecordNotFound,
			wantErr:              true,
		},
		{
			name:  "NG fail sending queue",
			input: &code.InvokeScanGitleaksRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockQueue: newFakeCodeQueue("failure", errors.New("something error")),
			wantErr:   true,
		},
		{
			name:  "NG NG db error when UpsertGitleaksSetting",
			input: &code.InvokeScanGitleaksRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockQueue:               newFakeCodeQueue("succeed", nil),
			mockUpsertGitleaksError: gorm.ErrInvalidDB,
			wantErr:                 true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB, sqs: c.mockQueue, logger: logging.NewLogger()}
			if c.mockGetGitleaksResponse != nil || c.mockGetGitleaksError != nil {
				mockDB.On("GetGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockGetGitleaksResponse, c.mockGetGitleaksError).Once()
			}
			if c.mockUpsertGitleaksResponse != nil || c.mockUpsertGitleaksError != nil {
				mockDB.On("UpsertGitleaksSetting", test.RepeatMockAnything(2)...).Return(c.mockUpsertGitleaksResponse, c.mockUpsertGitleaksError).Once()
			}
			_, err := svc.InvokeScanGitleaks(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
		})
	}
}

func TestInvokeScanDependency(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                         string
		input                        *code.InvokeScanDependencyRequest
		mockGetDependencyResponse    *model.CodeDependencySetting
		mockGetDependencyError       error
		mockQueue                    CodeQueue
		mockUpsertDependencyResponse *model.CodeDependencySetting
		mockUpsertDependencyError    error
		wantErr                      bool
	}{
		{
			name:  "OK",
			input: &code.InvokeScanDependencyRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockQueue:                    newFakeCodeQueue("succeed", nil),
			mockUpsertDependencyResponse: &model.CodeDependencySetting{},
		},
		{
			name:    "NG invalid param",
			input:   &code.InvokeScanDependencyRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:                   "NG db error when GetGitHubSetting",
			input:                  &code.InvokeScanDependencyRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetDependencyError: gorm.ErrRecordNotFound,
			wantErr:                true,
		},
		{
			name:  "NG fail sending queue",
			input: &code.InvokeScanDependencyRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockQueue: newFakeCodeQueue("failure", errors.New("something error")),
			wantErr:   true,
		},
		{
			name:  "NG NG db error when UpsertDependencySetting",
			input: &code.InvokeScanDependencyRequest{ProjectId: 1, GithubSettingId: 1},
			mockGetDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockQueue:                 newFakeCodeQueue("succeed", nil),
			mockUpsertDependencyError: gorm.ErrInvalidDB,
			wantErr:                   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			svc := CodeService{repository: mockDB, sqs: c.mockQueue, logger: logging.NewLogger()}
			if c.mockGetDependencyResponse != nil || c.mockGetDependencyError != nil {
				mockDB.On("GetDependencySetting", test.RepeatMockAnything(3)...).Return(c.mockGetDependencyResponse, c.mockGetDependencyError).Once()
			}
			if c.mockUpsertDependencyResponse != nil || c.mockUpsertDependencyError != nil {
				mockDB.On("UpsertDependencySetting", test.RepeatMockAnything(2)...).Return(c.mockUpsertDependencyResponse, c.mockUpsertDependencyError).Once()
			}
			_, err := svc.InvokeScanDependency(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
		})
	}
}

func TestInvokeScanAll(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name                         string
		ProjectID                    uint32
		mockListGitleaksResponse     *[]model.CodeGitleaksSetting
		mockListGitleaksError        error
		mockListDependencyResponse   *[]model.CodeDependencySetting
		mockListDependencyError      error
		mockListCodeScanResponse     *[]model.CodeCodeScanSetting
		mockListCodeScanError        error
		mockGetGitleaksResponse      *model.CodeGitleaksSetting
		mockGetGitleaksError         error
		mockGetDependencyResponse    *model.CodeDependencySetting
		mockGetDependencyError       error
		mockGetCodeScanResponse      *model.CodeCodeScanSetting
		mockGetCodeScanError         error
		mockGetGitHubSettingResponse *model.CodeGitHubSetting
		mockGetGitHubSettingError    error
		mockGithubClient             *FakeGithubClient
		mockIsActiveResponse         *project.IsActiveResponse
		mockIsActiveError            error
		mockQueue                    CodeQueue
		mockUpsertGitleaksResponse   *model.CodeGitleaksSetting
		mockUpsertGitleaksError      error
		mockUpsertDependencyResponse *model.CodeDependencySetting
		mockUpsertDependencyError    error
		mockUpsertCodeScanResponse   *model.CodeCodeScanSetting
		mockUpsertCodeScanError      error
		wantErr                      bool
	}{
		{
			name:                       "OK no data",
			ProjectID:                  1,
			mockListGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse:   &[]model.CodeCodeScanSetting{},
		},
		{
			name:      "OK scan gitleaks",
			ProjectID: 1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse:   &[]model.CodeCodeScanSetting{},
			mockIsActiveResponse:       &project.IsActiveResponse{Active: true},
			mockGetGitleaksResponse:    &model.CodeGitleaksSetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			mockQueue:                  newFakeCodeQueue("succeed", nil),
			mockUpsertGitleaksResponse: &model.CodeGitleaksSetting{},
		},
		{
			name:      "OK found gitleaks setting but projectID is zero",
			ProjectID: 1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 0, CodeDataSourceID: 1, ProjectID: 0, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse:   &[]model.CodeCodeScanSetting{},
		},
		{
			name:      "OK found gitleaks setting but project isn't active",
			ProjectID: 1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveResponse:       &project.IsActiveResponse{Active: false},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse:   &[]model.CodeCodeScanSetting{},
		},
		{
			name:                     "OK scan dependency",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockListCodeScanResponse:     &[]model.CodeCodeScanSetting{},
			mockIsActiveResponse:         &project.IsActiveResponse{Active: true},
			mockGetDependencyResponse:    &model.CodeDependencySetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			mockQueue:                    newFakeCodeQueue("succeed", nil),
			mockUpsertDependencyResponse: &model.CodeDependencySetting{},
		},
		{
			name:                     "OK found dependency setting but projectID is zero",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 0, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockListCodeScanResponse: &[]model.CodeCodeScanSetting{},
		},
		{
			name:                     "OK found dependency setting but project isn't active",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockListCodeScanResponse: &[]model.CodeCodeScanSetting{},
			mockIsActiveResponse:     &project.IsActiveResponse{Active: false},
		},
		{
			name:                       "OK scan CodeScan",
			ProjectID:                  1,
			mockListGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse: &[]model.CodeCodeScanSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: true, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveResponse:         &project.IsActiveResponse{Active: true},
			mockGetCodeScanResponse:      &model.CodeCodeScanSetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: true, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			mockGetGitHubSettingResponse: &model.CodeGitHubSetting{CodeGitHubSettingID: 1, ProjectID: 1, Type: "ORGANIZATION", TargetResource: "ca-risken", GitHubUser: "user", PersonalAccessToken: "", CreatedAt: now, UpdatedAt: now},
			mockQueue:                    newFakeCodeQueue("succeed", nil),
			mockUpsertCodeScanResponse:   &model.CodeCodeScanSetting{},
		},
		{
			name:                       "OK found CodeScan setting but projectID is zero",
			ProjectID:                  1,
			mockListGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse: &[]model.CodeCodeScanSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 0, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:                       "OK found CodeScan setting but project isn't active",
			ProjectID:                  1,
			mockListGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse: &[]model.CodeCodeScanSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveResponse: &project.IsActiveResponse{Active: false},
		},

		{
			name:                  "NG db error when ListGitleaksSetting",
			ProjectID:             1,
			mockListGitleaksError: gorm.ErrRecordNotFound,
			wantErr:               true,
		},
		{
			name:                     "NG db error when ListDependencySetting",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyError:  gorm.ErrRecordNotFound,
			wantErr:                  true,
		},
		{
			name:                       "NG db error when ListCodeScanSetting",
			ProjectID:                  1,
			mockListGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanError:      gorm.ErrRecordNotFound,
			wantErr:                    true,
		},

		{
			name:      "NG project client error when scanning gitleaks",
			ProjectID: 1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveError: errors.New("something error"),
			wantErr:           true,
		},
		{
			name:                     "NG project client error when scanning dependency",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveError: errors.New("something error"),
			wantErr:           true,
		},
		{
			name:      "NG error InvokeScanGitleaks",
			ProjectID: 1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveResponse: &project.IsActiveResponse{Active: true},
			mockGetGitleaksError: gorm.ErrInvalidDB,
			wantErr:              true,
		},
		{
			name:                     "NG error InvokeScanDependency",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveResponse:   &project.IsActiveResponse{Active: true},
			mockGetDependencyError: gorm.ErrInvalidDB,
			wantErr:                true,
		},
		{
			name:                       "NG error GetCodeScan",
			ProjectID:                  1,
			mockListGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
			mockListCodeScanResponse: &[]model.CodeCodeScanSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveResponse: &project.IsActiveResponse{Active: true},
			mockGetCodeScanError: gorm.ErrInvalidDB,
			wantErr:              true, // If GetCodeScanSetting returns an error, InvokeScanCodeScan returns an error
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewCodeRepoInterface(t)
			mockProject := projectmock.NewProjectServiceClient(t)
			githubClient := c.mockGithubClient
			if githubClient == nil {
				githubClient = newFakeGithubClient([]*ghub.Repository{
					{
						Name:       ghub.String("sample"),
						FullName:   ghub.String("ca-risken/sample"),
						Visibility: ghub.String("public"),
					},
				}, nil)
			}
			block, _ := aes.NewCipher([]byte("1234567890123456"))
			svc := CodeService{repository: mockDB, sqs: c.mockQueue, projectClient: mockProject, logger: logging.NewLogger(), githubClient: githubClient, cipherBlock: block}
			if c.mockListGitleaksResponse != nil || c.mockListGitleaksError != nil {
				mockDB.On("ListGitleaksSetting", test.RepeatMockAnything(2)...).Return(c.mockListGitleaksResponse, c.mockListGitleaksError).Once()
			}
			if c.mockListDependencyResponse != nil || c.mockListDependencyError != nil {
				mockDB.On("ListDependencySetting", test.RepeatMockAnything(2)...).Return(c.mockListDependencyResponse, c.mockListDependencyError).Once()
			}
			if c.mockListCodeScanResponse != nil || c.mockListCodeScanError != nil {
				mockDB.On("ListCodeScanSetting", test.RepeatMockAnything(2)...).Return(c.mockListCodeScanResponse, c.mockListCodeScanError).Once()
			}
			if c.mockGetGitleaksResponse != nil || c.mockGetGitleaksError != nil {
				mockDB.On("GetGitleaksSetting", test.RepeatMockAnything(3)...).Return(c.mockGetGitleaksResponse, c.mockGetGitleaksError).Once()
			}
			if c.mockUpsertGitleaksResponse != nil || c.mockUpsertGitleaksError != nil {
				mockDB.On("UpsertGitleaksSetting", test.RepeatMockAnything(2)...).Return(c.mockUpsertGitleaksResponse, c.mockUpsertGitleaksError).Once()
			}
			if c.mockGetDependencyResponse != nil || c.mockGetDependencyError != nil {
				mockDB.On("GetDependencySetting", test.RepeatMockAnything(3)...).Return(c.mockGetDependencyResponse, c.mockGetDependencyError).Once()
			}
			if c.mockUpsertDependencyResponse != nil || c.mockUpsertDependencyError != nil {
				mockDB.On("UpsertDependencySetting", test.RepeatMockAnything(2)...).Return(c.mockUpsertDependencyResponse, c.mockUpsertDependencyError).Once()
			}
			if c.mockGetCodeScanResponse != nil || c.mockGetCodeScanError != nil {
				// GetCodeScanSetting is called twice: once in InvokeScanCodeScan and once in listCodescanTargetRepository.
				// If the first call returns an error, the second call never happens.
				callCount := 2
				if c.mockGetCodeScanError != nil {
					callCount = 1
				}
				mockDB.On("GetCodeScanSetting", test.RepeatMockAnything(3)...).Return(c.mockGetCodeScanResponse, c.mockGetCodeScanError).Times(callCount)
			}
			if c.mockGetGitHubSettingResponse != nil || c.mockGetGitHubSettingError != nil {
				mockDB.On("GetGitHubSetting", test.RepeatMockAnything(3)...).Return(c.mockGetGitHubSettingResponse, c.mockGetGitHubSettingError).Once()
			}
			if c.mockUpsertCodeScanResponse != nil || c.mockUpsertCodeScanError != nil {
				mockDB.On("UpsertCodeScanSetting", test.RepeatMockAnything(2)...).Return(c.mockUpsertCodeScanResponse, c.mockUpsertCodeScanError).Once()
			}

			if c.mockIsActiveResponse != nil || c.mockIsActiveError != nil {
				mockProject.On("IsActive", ctx, &project.IsActiveRequest{ProjectId: c.ProjectID}).Return(c.mockIsActiveResponse, c.mockIsActiveError).Once()
			}
			_, err := svc.InvokeScanAll(ctx, &empty.Empty{})
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
		})
	}
}

type FakeCodeQueue struct {
	resp *sqs.SendMessageOutput
	err  error
}

func newFakeCodeQueue(msg string, err error) *FakeCodeQueue {
	return &FakeCodeQueue{
		resp: &sqs.SendMessageOutput{MessageId: &msg},
		err:  err,
	}
}

func (c *FakeCodeQueue) Send(ctx context.Context, url string, msg interface{}) (*sqs.SendMessageOutput, error) {
	return c.resp, c.err
}

type FakeGithubClient struct {
	repos []*ghub.Repository
	err   error
}

func newFakeGithubClient(repos []*ghub.Repository, err error) *FakeGithubClient {
	return &FakeGithubClient{
		repos: repos,
		err:   err,
	}
}

func (g *FakeGithubClient) ListRepository(ctx context.Context, config *code.GitHubSetting, repoName string) ([]*ghub.Repository, error) {
	return g.repos, g.err
}

func (g *FakeGithubClient) Clone(ctx context.Context, token string, cloneURL string, dstDir string) error {
	return g.err
}

func TestIsGitHubAuthError(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-GitHub error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name: "GitHub authentication error (401)",
			err: &ghub.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusUnauthorized,
				},
				Message: "Bad credentials",
			},
			expected: true,
		},
		{
			name: "GitHub error with 404 status",
			err: &ghub.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
				Message: "Not Found",
			},
			expected: false,
		},
		{
			name: "GitHub error with 403 status",
			err: &ghub.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
				Message: "Forbidden",
			},
			expected: false,
		},
		{
			name: "GitHub error with 500 status",
			err: &ghub.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
				Message: "Internal Server Error",
			},
			expected: false,
		},
		{
			name: "GitHub ErrorResponse with nil Response",
			err: &ghub.ErrorResponse{
				Response: nil,
				Message:  "Some error",
			},
			expected: false,
		},
		{
			name: "wrapped GitHub authentication error",
			err: fmt.Errorf("wrapped error: %w", &ghub.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusUnauthorized,
				},
				Message: "Bad credentials",
			}),
			expected: true,
		},
		{
			name: "wrapped GitHub error with 404 status",
			err: fmt.Errorf("wrapped error: %w", &ghub.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
				Message: "Not Found",
			}),
			expected: false,
		},
		{
			name: "double wrapped GitHub authentication error",
			err: fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", &ghub.ErrorResponse{
				Response: &http.Response{
					StatusCode: http.StatusUnauthorized,
				},
				Message: "Bad credentials",
			})),
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := isGitHubAuthError(c.err)
			if got != c.expected {
				t.Errorf("isGitHubAuthError(%v) = %v, want %v", c.err, got, c.expected)
			}
		})
	}
}
