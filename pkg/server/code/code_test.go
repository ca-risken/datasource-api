package code

import (
	"context"
	"crypto/aes"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/ca-risken/datasource-api/pkg/db/mock"
	"github.com/ca-risken/datasource-api/pkg/model"

	"github.com/ca-risken/datasource-api/proto/code"
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
				{CodeDataSourceId: 1, Name: "one", Description: "desc", MaxScore: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{CodeDataSourceId: 2, Name: "two", Description: "desc", MaxScore: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponse: &[]model.CodeDataSource{
				{CodeDataSourceID: 1, Name: "one", Description: "desc", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
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

func TestListGitleaks(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
	cases := []struct {
		name                        string
		input                       *code.ListGitleaksRequest
		want                        *code.ListGitleaksResponse
		mockGitHubSettingResponse   *[]model.CodeGitHubSetting
		mockGitHubSettingError      error
		mockGitleaksSettingResponse *[]model.CodeGitleaksSetting
		mockGitleaksSettingError    error
		wantErr                     bool
	}{
		{
			name:  "OK",
			input: &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			want: &code.ListGitleaksResponse{Gitleaks: []*code.Gitleaks{
				{GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, ScanPublic: true, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{GitleaksId: 2, CodeDataSourceId: 1, Name: "two", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo2", GithubUser: "user", PersonalAccessToken: maskData, ScanPublic: true, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockGitHubSettingResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token"},
				{CodeGitHubSettingID: 2, Name: "two", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token"},
			},
			mockGitleaksSettingResponse: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, RepositoryPattern: "repo", ScanPublic: true, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, RepositoryPattern: "repo2", ScanPublic: true, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:                      "OK empty",
			input:                     &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			want:                      &code.ListGitleaksResponse{},
			mockGitHubSettingResponse: &[]model.CodeGitHubSetting{},
			mockGitHubSettingError:    nil,
		},
		{
			name:  "OK return empty when gitleaks_setting is not found",
			input: &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			want:  &code.ListGitleaksResponse{},
			mockGitHubSettingResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token"},
			},
			mockGitleaksSettingResponse: &[]model.CodeGitleaksSetting{},
		},
		{
			name:    "NG invalid param",
			input:   &code.ListGitleaksRequest{CodeDataSourceId: 1},
			wantErr: true,
		},
		{
			name:                   "Invalid DB error",
			input:                  &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			mockGitHubSettingError: gorm.ErrInvalidDB,
			wantErr:                true,
		},
		{
			name:  "Invalid DB error (getGitleaks)",
			input: &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			mockGitHubSettingResponse: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token"},
			},
			mockGitleaksSettingError: gorm.ErrInvalidDB,
			wantErr:                  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGitHubSettingResponse != nil || c.mockGitHubSettingError != nil {
				mockDB.On("ListGitHubSetting").Return(c.mockGitHubSettingResponse, c.mockGitHubSettingError).Once()
			}
			if c.mockGitleaksSettingResponse != nil || c.mockGitleaksSettingError != nil {
				mockDB.On("ListGitleaksSetting").Return(c.mockGitleaksSettingResponse, c.mockGitleaksSettingError).Once()
			}
			got, err := svc.ListGitleaks(ctx, c.input)
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

func TestGetGitleaks(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
	cases := []struct {
		name                        string
		input                       *code.GetGitleaksRequest
		want                        *code.GetGitleaksResponse
		mockGitHubSettingResponse   *model.CodeGitHubSetting
		mockGitHubSettingError      error
		mockGitleaksSettingResponse *model.CodeGitleaksSetting
		mockGitleaksSettingError    error
		wantErr                     bool
	}{
		{
			name:                        "OK",
			input:                       &code.GetGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			want:                        &code.GetGitleaksResponse{Gitleaks: &code.Gitleaks{GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: "token", ScanPublic: true, ScanInternal: false, ScanPrivate: false, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGitHubSettingResponse:   &model.CodeGitHubSetting{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token"},
			mockGitleaksSettingResponse: &model.CodeGitleaksSetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, RepositoryPattern: "repo", ScanPublic: true, ScanInternal: false, ScanPrivate: false, Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:                   "OK empty",
			input:                  &code.GetGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			want:                   &code.GetGitleaksResponse{},
			mockGitHubSettingError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.GetGitleaksRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:                      "NG not found gitleaks_setting",
			input:                     &code.GetGitleaksRequest{ProjectId: 1},
			mockGitHubSettingResponse: &model.CodeGitHubSetting{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token"},
			mockGitleaksSettingError:  gorm.ErrRecordNotFound,
			wantErr:                   true,
		},
		{
			name:                   "Invalid DB error",
			input:                  &code.GetGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			mockGitHubSettingError: gorm.ErrInvalidDB,
			wantErr:                true,
		},
		{
			name:                      "Invalid DB error (getGitleaks)",
			input:                     &code.GetGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			mockGitHubSettingResponse: &model.CodeGitHubSetting{CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token"},
			mockGitleaksSettingError:  gorm.ErrInvalidDB,
			wantErr:                   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGitHubSettingResponse != nil || c.mockGitHubSettingError != nil {
				mockDB.On("GetGitHubSetting").Return(c.mockGitHubSettingResponse, c.mockGitHubSettingError).Once()
			}
			if c.mockGitleaksSettingResponse != nil || c.mockGitleaksSettingError != nil {
				mockDB.On("GetGitleaksSetting").Return(c.mockGitleaksSettingResponse, c.mockGitleaksSettingError).Once()
			}
			got, err := svc.GetGitleaks(ctx, c.input)
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

func TestPutGitleaks(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	key := []byte("1234567890123456")
	block, err := aes.NewCipher(key)
	assert.NoError(t, err)
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{
		repository:  &mockDB,
		cipherBlock: block,
	}
	cases := []struct {
		name                        string
		input                       *code.PutGitleaksRequest
		want                        *code.PutGitleaksResponse
		mockGitHubSettingResponse   *model.CodeGitHubSetting
		mockGitHubSettingError      error
		mockGitleaksSettingResponse *model.CodeGitleaksSetting
		mockGitleaksSettingError    error
		wantErr                     bool
	}{
		{
			name: "OK",
			input: &code.PutGitleaksRequest{ProjectId: 1, Gitleaks: &code.GitleaksForUpsert{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: now.Unix()},
			},
			want: &code.PutGitleaksResponse{Gitleaks: &code.Gitleaks{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockGitHubSettingResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token",
			},
			mockGitleaksSettingResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name: "OK(empty)",
			input: &code.PutGitleaksRequest{ProjectId: 1, Gitleaks: &code.GitleaksForUpsert{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want: &code.PutGitleaksResponse{Gitleaks: &code.Gitleaks{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: 0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockGitHubSettingResponse: &model.CodeGitHubSetting{
				CodeGitHubSettingID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token",
			},
			mockGitleaksSettingResponse: &model.CodeGitleaksSetting{
				CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, RepositoryPattern: "repo", Status: "OK", StatusDetail: "detail", ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:    "NG invalid param",
			input:   &code.PutGitleaksRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &code.PutGitleaksRequest{ProjectId: 1, Gitleaks: &code.GitleaksForUpsert{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			mockGitHubSettingError: gorm.ErrInvalidDB,
			wantErr:                true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGitHubSettingResponse != nil || c.mockGitHubSettingError != nil {
				mockDB.On("UpsertGitHubSetting").Return(c.mockGitHubSettingResponse, c.mockGitHubSettingError).Once()
			}
			if c.mockGitleaksSettingResponse != nil || c.mockGitleaksSettingError != nil {
				mockDB.On("UpsertGitleaksSetting").Return(c.mockGitleaksSettingResponse, c.mockGitleaksSettingError).Once()
			}
			got, err := svc.PutGitleaks(ctx, c.input)
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

func TestDeleteGitleaks(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
	cases := []struct {
		name                               string
		input                              *code.DeleteGitleaksRequest
		isCalledDeleteGitHubSetting        bool
		mockGitHubSettingError             error
		isCalledDeleteGitleaksSetting      bool
		mockGitleaksSettingError           error
		mockListEnterpriseOrgResponse      *[]model.CodeGitHubEnterpriseOrg
		mockListEnterpriseOrgError         error
		isCalledDeleteEnterpriseOrgSetting bool
		mockDeleteEnterpriseOrgError       error
		wantErr                            bool
	}{
		{
			name:                          "OK",
			input:                         &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			isCalledDeleteGitHubSetting:   true,
			isCalledDeleteGitleaksSetting: true,
			mockListEnterpriseOrgResponse: &[]model.CodeGitHubEnterpriseOrg{
				{CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			},
			isCalledDeleteEnterpriseOrgSetting: true,
		},
		{
			name:                          "OK (enterprise_org empty)",
			input:                         &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			isCalledDeleteGitHubSetting:   true,
			isCalledDeleteGitleaksSetting: true,
			mockListEnterpriseOrgResponse: &[]model.CodeGitHubEnterpriseOrg{},
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteGitleaksRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:                        "Invalid DB error (DeleteGitHubSetting)",
			input:                       &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			isCalledDeleteGitHubSetting: true,
			mockGitHubSettingError:      gorm.ErrInvalidDB,
			wantErr:                     true,
		},
		{
			name:                          "Invalid DB error (DeleteGitleaksSetting)",
			input:                         &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			isCalledDeleteGitHubSetting:   true,
			isCalledDeleteGitleaksSetting: true,
			mockGitleaksSettingError:      gorm.ErrInvalidDB,
			wantErr:                       true,
		},
		{
			name:                          "Invalid DB error (ListGitHubEnterpriseOrg)",
			input:                         &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			isCalledDeleteGitHubSetting:   true,
			isCalledDeleteGitleaksSetting: true,
			mockListEnterpriseOrgError:    gorm.ErrInvalidDB,
			wantErr:                       true,
		},
		{
			name:                          "Invalid DB error (DeleteGitHubEnterpriseOrg)",
			input:                         &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			isCalledDeleteGitHubSetting:   true,
			isCalledDeleteGitleaksSetting: true,
			mockListEnterpriseOrgResponse: &[]model.CodeGitHubEnterpriseOrg{
				{CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			},
			isCalledDeleteEnterpriseOrgSetting: true,
			mockDeleteEnterpriseOrgError:       gorm.ErrInvalidDB,
			wantErr:                            true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.isCalledDeleteGitHubSetting {
				mockDB.On("DeleteGitHubSetting").Return(c.mockGitHubSettingError).Once()
			}
			if c.isCalledDeleteGitleaksSetting {
				mockDB.On("DeleteGitleaksSetting").Return(c.mockGitleaksSettingError).Once()
			}
			if c.mockListEnterpriseOrgResponse != nil || c.mockListEnterpriseOrgError != nil {
				mockDB.On("ListGitHubEnterpriseOrg").Return(c.mockListEnterpriseOrgResponse, c.mockListEnterpriseOrgError).Once()
			}
			if c.isCalledDeleteEnterpriseOrgSetting {
				mockDB.On("DeleteGitHubEnterpriseOrg").Return(c.mockDeleteEnterpriseOrgError).Once()
			}
			_, err := svc.DeleteGitleaks(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestListEnterpriseOrg(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
	cases := []struct {
		name                      string
		input                     *code.ListEnterpriseOrgRequest
		want                      *code.ListEnterpriseOrgResponse
		mockGitHubSettingResponse *[]model.CodeGitHubEnterpriseOrg
		mockGitHubSettingError    error
		wantErr                   bool
	}{
		{
			name:  "OK",
			input: &code.ListEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1},
			want: &code.ListEnterpriseOrgResponse{EnterpriseOrg: []*code.EnterpriseOrg{
				{GitleaksId: 1, Login: "one", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{GitleaksId: 2, Login: "two", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockGitHubSettingResponse: &[]model.CodeGitHubEnterpriseOrg{
				{CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, Organization: "two", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:                   "OK empty",
			input:                  &code.ListEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1},
			want:                   &code.ListEnterpriseOrgResponse{},
			mockGitHubSettingError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.ListEnterpriseOrgRequest{GitleaksId: 1},
			wantErr: true,
		},
		{
			name:                   "Invalid DB error",
			input:                  &code.ListEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1},
			mockGitHubSettingError: gorm.ErrInvalidDB,
			wantErr:                true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGitHubSettingResponse != nil || c.mockGitHubSettingError != nil {
				mockDB.On("ListGitHubEnterpriseOrg").Return(c.mockGitHubSettingResponse, c.mockGitHubSettingError).Once()
			}
			got, err := svc.ListEnterpriseOrg(ctx, c.input)
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

func TestPutEnterpriseOrg(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
	cases := []struct {
		name                      string
		input                     *code.PutEnterpriseOrgRequest
		want                      *code.PutEnterpriseOrgResponse
		mockGitHubSettingResponse *model.CodeGitHubEnterpriseOrg
		mockGitHubSettingError    error
		wantErr                   bool
	}{
		{
			name: "OK",
			input: &code.PutEnterpriseOrgRequest{ProjectId: 1, EnterpriseOrg: &code.EnterpriseOrgForUpsert{
				GitleaksId: 1, Login: "one", ProjectId: 1},
			},
			want: &code.PutEnterpriseOrgResponse{EnterpriseOrg: &code.EnterpriseOrg{
				GitleaksId: 1, Login: "one", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockGitHubSettingResponse: &model.CodeGitHubEnterpriseOrg{
				CodeGitHubSettingID: 1, Organization: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:    "NG invalid param",
			input:   &code.PutEnterpriseOrgRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &code.PutEnterpriseOrgRequest{ProjectId: 1, EnterpriseOrg: &code.EnterpriseOrgForUpsert{
				GitleaksId: 1, Login: "one", ProjectId: 1},
			},
			mockGitHubSettingError: gorm.ErrInvalidDB,
			wantErr:                true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGitHubSettingResponse != nil || c.mockGitHubSettingError != nil {
				mockDB.On("UpsertGitHubEnterpriseOrg").Return(c.mockGitHubSettingResponse, c.mockGitHubSettingError).Once()
			}
			got, err := svc.PutEnterpriseOrg(ctx, c.input)
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

func TestDeleteEnterpriseOrg(t *testing.T) {
	var ctx context.Context
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
	cases := []struct {
		name                   string
		input                  *code.DeleteEnterpriseOrgRequest
		mockGitHubSettingError error
		wantErr                bool
	}{
		{
			name:  "OK",
			input: &code.DeleteEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1, Login: "login"},
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1},
			wantErr: true,
		},
		{
			name:                   "Invalid DB error",
			input:                  &code.DeleteEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1, Login: "login"},
			mockGitHubSettingError: gorm.ErrInvalidDB,
			wantErr:                true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteGitHubEnterpriseOrg").Return(c.mockGitHubSettingError).Once()
			_, err := svc.DeleteEnterpriseOrg(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}
