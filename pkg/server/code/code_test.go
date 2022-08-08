package code

import (
	"context"
	"crypto/aes"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	projectmock "github.com/ca-risken/core/proto/project/mocks"
	mockdb "github.com/ca-risken/datasource-api/pkg/db/mock"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	length65string = "12345678901234567890123456789012345678901234567890123456789012345"
)

func TestListDataSource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
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
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListCodeDataSource").Return(c.mockResponse, c.mockError).Once()
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
	var ctx context.Context
	now := time.Now()

	cases := []struct {
		name                   string
		input                  *code.ListGitHubSettingRequest
		want                   *code.ListGitHubSettingResponse
		mockResponse           *[]model.CodeGitHubSetting
		mockGitleaksResponse   *[]model.CodeGitleaksSetting
		mockDependencyResponse *[]model.CodeDependencySetting
		mockError              error
		mockGitleaksError      error
		mockDependencyError    error
		wantErr                bool
	}{
		{
			name:  "OK",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			want: &code.ListGitHubSettingResponse{GithubSetting: []*code.GitHubSetting{
				{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
				{GithubSettingId: 2, Name: "two", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 2, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo2", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					DependencySetting: &code.DependencySetting{GithubSettingId: 2, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			}},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, Name: "two", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo2", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:                   "OK empty",
			input:                  &code.ListGitHubSettingRequest{ProjectId: 1},
			want:                   &code.ListGitHubSettingResponse{},
			mockResponse:           &[]model.CodeGitHubSetting{},
			mockGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockDependencyResponse: &[]model.CodeDependencySetting{},
		},
		{
			name:  "OK gitleaks setting empty",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			want: &code.ListGitHubSettingResponse{GithubSetting: []*code.GitHubSetting{
				{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				}}},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now}},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:  "OK dependency setting empty",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			want: &code.ListGitHubSettingResponse{GithubSetting: []*code.GitHubSetting{
				{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
					GitleaksSetting: &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				}}},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now}},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockDependencyResponse: &[]model.CodeDependencySetting{},
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
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			},
			mockGitleaksResponse: nil,
			mockGitleaksError:    gorm.ErrInvalidDB,
			wantErr:              true,
		},
		{
			name:  "Invalid DB error when getDependencySetting",
			input: &code.ListGitHubSettingRequest{ProjectId: 1},
			mockResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			},
			mockGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockDependencyError:  gorm.ErrInvalidDB,
			wantErr:              true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListGitHubSetting").Return(c.mockResponse, c.mockError).Once()
			}
			if c.mockGitleaksResponse != nil || c.mockGitleaksError != nil {
				mockDB.On("ListGitleaksSetting").Return(c.mockGitleaksResponse, c.mockGitleaksError).Once()
			}
			if c.mockDependencyResponse != nil || c.mockDependencyError != nil {
				mockDB.On("ListDependencySetting").Return(c.mockDependencyResponse, c.mockDependencyError).Once()
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
	var ctx context.Context
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
		wantErr                bool
	}{
		{
			name:  "OK",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			want: &code.GetGitHubSettingResponse{GithubSetting: &code.GitHubSetting{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				GitleaksSetting:   &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyResponse: &model.CodeDependencySetting{
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
			want: &code.GetGitHubSettingResponse{GithubSetting: &code.GitHubSetting{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				DependencySetting: &code.DependencySetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			mockGitleaksError: gorm.ErrRecordNotFound,
			mockDependencyResponse: &model.CodeDependencySetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:  "OK dependency setting empty",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			want: &code.GetGitHubSettingResponse{GithubSetting: &code.GitHubSetting{GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
				GitleaksSetting: &code.GitleaksSetting{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyError: gorm.ErrRecordNotFound,
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
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
			mockGitleaksError: gorm.ErrInvalidDB,
			wantErr:           true,
		},
		{
			name:  "Invalid DB error when GetDependencySetting",
			input: &code.GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
			mockGitleaksResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			mockDependencyError: gorm.ErrInvalidDB,
			wantErr:             true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("GetGitHubSetting").Return(c.mockResponse, c.mockError).Once()
			}
			if c.mockGitleaksResponse != nil || c.mockGitleaksError != nil {
				mockDB.On("GetGitleaksSetting").Return(c.mockGitleaksResponse, c.mockGitleaksError).Once()
			}
			if c.mockDependencyResponse != nil || c.mockDependencyError != nil {
				mockDB.On("GetDependencySetting").Return(c.mockDependencyResponse, c.mockDependencyError).Once()
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
	var ctx context.Context
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
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData},
			},
			want: &code.PutGitHubSettingResponse{GithubSetting: &code.GitHubSetting{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name: "OK(empty)",
			input: &code.PutGitHubSettingRequest{ProjectId: 1, GithubSetting: &code.GitHubSettingForUpsert{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData},
			},
			want: &code.PutGitHubSettingResponse{GithubSetting: &code.GitHubSetting{
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now,
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
				GithubSettingId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: maskData},
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{
				repository:  &mockDB,
				cipherBlock: block,
			}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("UpsertGitHubSetting").Return(c.mockResponse, c.mockError).Once()
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
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name                         string
		input                        *code.DeleteGitHubSettingRequest
		isCalledDeleteGithubSetting  bool
		mockError                    error
		isCalledDeleteGitleaks       bool
		mockGitleaksError            error
		isCalledDeleteDependency     bool
		mockDependencyError          error
		isCalledListEnterpriseOrg    bool
		mockListEnterpriseOrg        *[]model.CodeGitHubEnterpriseOrg
		mockEnterpriserOrgError      error
		isCalledDeleteEnterpriseOrg  bool
		mockDeleteEnterpriseOrgError error
		wantErr                      bool
	}{
		{
			name:                        "OK",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			isCalledDeleteGithubSetting: true,
			mockError:                   nil,
			isCalledDeleteGitleaks:      true,
			mockGitleaksError:           nil,
			isCalledDeleteDependency:    true,
			mockDependencyError:         nil,
			isCalledListEnterpriseOrg:   true,
			mockListEnterpriseOrg: &[]model.CodeGitHubEnterpriseOrg{
				{CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			},
			mockEnterpriserOrgError:      nil,
			isCalledDeleteEnterpriseOrg:  true,
			mockDeleteEnterpriseOrgError: nil,
		},
		{
			name:                        "OK enterprise org empty",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			isCalledDeleteGithubSetting: true,
			mockError:                   nil,
			isCalledDeleteGitleaks:      true,
			mockGitleaksError:           nil,
			isCalledDeleteDependency:    true,
			mockDependencyError:         nil,
			isCalledListEnterpriseOrg:   true,
			mockListEnterpriseOrg:       &[]model.CodeGitHubEnterpriseOrg{},
		},
		{
			name:                        "NG invalid param",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1},
			isCalledDeleteGithubSetting: false,
			wantErr:                     true,
		},
		{
			name:                        "Invalid DB error",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockError:                   gorm.ErrInvalidDB,
			isCalledDeleteGithubSetting: true,
			wantErr:                     true,
		},
		{
			name:                        "Invalid DB error (deleteGitleaks)",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			isCalledDeleteGithubSetting: true,
			mockError:                   nil,
			isCalledDeleteGitleaks:      true,
			mockGitleaksError:           gorm.ErrInvalidDB,
			wantErr:                     true,
		},
		{
			name:                        "Invalid DB error (deleteGitleaks)",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			isCalledDeleteGithubSetting: true,
			mockError:                   nil,
			isCalledDeleteGitleaks:      true,
			mockGitleaksError:           nil,
			isCalledDeleteDependency:    true,
			mockDependencyError:         gorm.ErrInvalidDB,
			wantErr:                     true,
		},
		{
			name:                        "Invalid DB error (listEnterpriseOrg)",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			isCalledDeleteGithubSetting: true,
			mockError:                   nil,
			isCalledDeleteGitleaks:      true,
			mockGitleaksError:           nil,
			isCalledDeleteDependency:    true,
			mockDependencyError:         nil,
			isCalledListEnterpriseOrg:   true,
			mockEnterpriserOrgError:     gorm.ErrInvalidDB,
			wantErr:                     true,
		},
		{
			name:                        "Invalid DB error (deleteEnterpriseOrg)",
			input:                       &code.DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
			isCalledDeleteGithubSetting: true,
			mockError:                   nil,
			isCalledDeleteGitleaks:      true,
			mockGitleaksError:           nil,
			isCalledDeleteDependency:    true,
			mockDependencyError:         nil,
			isCalledListEnterpriseOrg:   true,
			mockListEnterpriseOrg: &[]model.CodeGitHubEnterpriseOrg{
				{CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			},
			isCalledDeleteEnterpriseOrg:  true,
			mockDeleteEnterpriseOrgError: gorm.ErrInvalidDB,
			wantErr:                      true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB}
			if c.isCalledDeleteGithubSetting {
				mockDB.On("DeleteGitHubSetting").Return(c.mockError).Once()
			}
			if c.isCalledDeleteGitleaks {
				mockDB.On("DeleteGitleaksSetting").Return(c.mockGitleaksError).Once()
			}
			if c.isCalledDeleteDependency {
				mockDB.On("DeleteDependencySetting").Return(c.mockDependencyError).Once()
			}
			if c.isCalledListEnterpriseOrg {
				mockDB.On("ListGitHubEnterpriseOrg").Return(c.mockListEnterpriseOrg, c.mockEnterpriserOrgError).Once()
			}
			if c.isCalledDeleteEnterpriseOrg {
				mockDB.On("DeleteGitHubEnterpriseOrg").Return(c.mockDeleteEnterpriseOrgError).Once()
			}
			_, err := svc.DeleteGitHubSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestPutGitleaksSetting(t *testing.T) {
	var ctx context.Context
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
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{
				repository: &mockDB,
			}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("UpsertGitleaksSetting").Return(c.mockResponse, c.mockError).Once()
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
	var ctx context.Context
	cases := []struct {
		name      string
		input     *code.DeleteGitleaksSettingRequest
		mockError error
		wantErr   bool
	}{
		{
			name:  "OK",
			input: &code.DeleteGitleaksSettingRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteGitleaksSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.DeleteGitleaksSettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB}
			mockDB.On("DeleteGitleaksSetting").Return(c.mockError).Once()
			_, err := svc.DeleteGitleaksSetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestPutDependencySetting(t *testing.T) {
	var ctx context.Context
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
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{
				repository: &mockDB,
			}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("UpsertDependencySetting").Return(c.mockResponse, c.mockError).Once()
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
	var ctx context.Context
	cases := []struct {
		name      string
		input     *code.DeleteDependencySettingRequest
		mockError error
		wantErr   bool
	}{
		{
			name:  "OK",
			input: &code.DeleteDependencySettingRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteDependencySettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.DeleteDependencySettingRequest{ProjectId: 1, GithubSettingId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB}
			mockDB.On("DeleteDependencySetting").Return(c.mockError).Once()
			_, err := svc.DeleteDependencySetting(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestListGitHubEnterpriseOrg(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *code.ListGitHubEnterpriseOrgRequest
		want         *code.ListGitHubEnterpriseOrgResponse
		mockResponse *[]model.CodeGitHubEnterpriseOrg
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &code.ListGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1},
			want: &code.ListGitHubEnterpriseOrgResponse{GithubEnterpriseOrg: []*code.GitHubEnterpriseOrg{
				{GithubSettingId: 1, Organization: "one", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{GithubSettingId: 2, Organization: "two", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &[]model.CodeGitHubEnterpriseOrg{
				{CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, Organization: "two", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty",
			input:     &code.ListGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1},
			want:      &code.ListGitHubEnterpriseOrgResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.ListGitHubEnterpriseOrgRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.ListGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListGitHubEnterpriseOrg").Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.ListGitHubEnterpriseOrg(ctx, c.input)
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

func TestPutGitHubEnterpriseOrg(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *code.PutGitHubEnterpriseOrgRequest
		want         *code.PutGitHubEnterpriseOrgResponse
		mockResponse *model.CodeGitHubEnterpriseOrg
		mockError    error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &code.PutGitHubEnterpriseOrgRequest{ProjectId: 1, GithubEnterpriseOrg: &code.GitHubEnterpriseOrgForUpsert{
				GithubSettingId: 1, Organization: "one", ProjectId: 1},
			},
			want: &code.PutGitHubEnterpriseOrgResponse{GithubEnterpriseOrg: &code.GitHubEnterpriseOrg{
				GithubSettingId: 1, Organization: "one", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponse: &model.CodeGitHubEnterpriseOrg{
				CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:    "NG invalid param",
			input:   &code.PutGitHubEnterpriseOrgRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &code.PutGitHubEnterpriseOrgRequest{ProjectId: 1, GithubEnterpriseOrg: &code.GitHubEnterpriseOrgForUpsert{
				GithubSettingId: 1, Organization: "one", ProjectId: 1},
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("UpsertGitHubEnterpriseOrg").Return(c.mockResponse, c.mockError).Once()
			}
			got, err := svc.PutGitHubEnterpriseOrg(ctx, c.input)
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

func TestDeleteGitHubEnterpriseOrg(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name      string
		input     *code.DeleteGitHubEnterpriseOrgRequest
		mockError error
		wantErr   bool
	}{
		{
			name:  "OK",
			input: &code.DeleteGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1, Organization: "Organization"},
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.DeleteGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1, Organization: "Organization"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		mockDB := mockdb.MockCodeRepository{}
		svc := CodeService{repository: &mockDB}
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteGitHubEnterpriseOrg").Return(c.mockError).Once()
			_, err := svc.DeleteGitHubEnterpriseOrg(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestInvokeScan(t *testing.T) {
	var ctx context.Context
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
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB, sqs: c.mockQueue, logger: logging.NewLogger()}
			if c.mockGetGitleaksResponse != nil || c.mockGetGitleaksError != nil {
				mockDB.On("GetGitleaksSetting").Return(c.mockGetGitleaksResponse, c.mockGetGitleaksError).Once()
			}
			if c.mockUpsertGitleaksResponse != nil || c.mockUpsertGitleaksError != nil {
				mockDB.On("UpsertGitleaksSetting").Return(c.mockUpsertGitleaksResponse, c.mockUpsertGitleaksError).Once()
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
	var ctx context.Context
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
			mockDB := mockdb.MockCodeRepository{}
			svc := CodeService{repository: &mockDB, sqs: c.mockQueue, logger: logging.NewLogger()}
			if c.mockGetDependencyResponse != nil || c.mockGetDependencyError != nil {
				mockDB.On("GetDependencySetting").Return(c.mockGetDependencyResponse, c.mockGetDependencyError).Once()
			}
			if c.mockUpsertDependencyResponse != nil || c.mockUpsertDependencyError != nil {
				mockDB.On("UpsertDependencySetting").Return(c.mockUpsertDependencyResponse, c.mockUpsertDependencyError).Once()
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
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name                         string
		ProjectID                    uint32
		mockListGitleaksResponse     *[]model.CodeGitleaksSetting
		mockListGitleaksError        error
		mockListDependencyResponse   *[]model.CodeDependencySetting
		mockListDependencyError      error
		mockGetGitleaksResponse      *model.CodeGitleaksSetting
		mockGetGitleaksError         error
		mockGetDependencyResponse    *model.CodeDependencySetting
		mockGetDependencyError       error
		mockIsActiveResponse         *project.IsActiveResponse
		mockIsActiveError            error
		mockQueue                    CodeQueue
		mockUpsertGitleaksResponse   *model.CodeGitleaksSetting
		mockUpsertGitleaksError      error
		mockUpsertDependencyResponse *model.CodeDependencySetting
		mockUpsertDependencyError    error
		wantErr                      bool
	}{
		{
			name:                       "OK no data",
			ProjectID:                  1,
			mockListGitleaksResponse:   &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
		},
		{
			name:      "OK scan gitleaks",
			ProjectID: 1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
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
		},
		{
			name:      "OK found gitleaks setting but project isn't active",
			ProjectID: 1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", ScanPublic: false, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			mockIsActiveResponse:       &project.IsActiveResponse{Active: false},
			mockListDependencyResponse: &[]model.CodeDependencySetting{},
		},
		{
			name:                     "OK scan dependency",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
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
			mockIsActiveResponse: &project.IsActiveResponse{Active: false},
		},
		{
			name:                     "OK found dependency setting but project isn't active",
			ProjectID:                1,
			mockListGitleaksResponse: &[]model.CodeGitleaksSetting{},
			mockListDependencyResponse: &[]model.CodeDependencySetting{
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
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mockdb.MockCodeRepository{}
			mockProject := &projectmock.ProjectServiceClient{}
			svc := CodeService{repository: &mockDB, sqs: c.mockQueue, projectClient: mockProject, logger: logging.NewLogger()}
			if c.mockListGitleaksResponse != nil || c.mockListGitleaksError != nil {
				mockDB.On("ListGitleaksSetting").Return(c.mockListGitleaksResponse, c.mockListGitleaksError).Once()
			}
			if c.mockListDependencyResponse != nil || c.mockListDependencyError != nil {
				mockDB.On("ListDependencySetting").Return(c.mockListDependencyResponse, c.mockListDependencyError).Once()
			}
			if c.mockGetGitleaksResponse != nil || c.mockGetGitleaksError != nil {
				mockDB.On("GetGitleaksSetting").Return(c.mockGetGitleaksResponse, c.mockGetGitleaksError).Once()
			}
			if c.mockUpsertGitleaksResponse != nil || c.mockUpsertGitleaksError != nil {
				mockDB.On("UpsertGitleaksSetting").Return(c.mockUpsertGitleaksResponse, c.mockUpsertGitleaksError).Once()
			}
			if c.mockGetDependencyResponse != nil || c.mockGetDependencyError != nil {
				mockDB.On("GetDependencySetting").Return(c.mockGetDependencyResponse, c.mockGetDependencyError).Once()
			}
			if c.mockUpsertDependencyResponse != nil || c.mockUpsertDependencyError != nil {
				mockDB.On("UpsertDependencySetting").Return(c.mockUpsertDependencyResponse, c.mockUpsertDependencyError).Once()
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
