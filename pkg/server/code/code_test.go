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
		mockResponce *[]model.CodeDataSource
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
			mockResponce: &[]model.CodeDataSource{
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListCodeDataSource").Return(c.mockResponce, c.mockError).Once()
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
		name         string
		input        *code.ListGitleaksRequest
		want         *code.ListGitleaksResponse
		mockResponce *[]model.CodeGitleaks
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			want: &code.ListGitleaksResponse{Gitleaks: []*code.Gitleaks{
				{GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, ScanPublic: true, ScanInternal: false, ScanPrivate: false, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{GitleaksId: 2, CodeDataSourceId: 1, Name: "two", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, ScanPublic: true, ScanInternal: false, ScanPrivate: false, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.CodeGitleaks{
				{GitleaksID: 1, CodeDataSourceID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: "token", ScanPublic: true, ScanInternal: false, ScanPrivate: false, GitleaksConfig: "conf", Status: "OK", StatusDetail: "detail", ScanAt: now, ScanSucceededAt: &now, CreatedAt: now, UpdatedAt: now},
				{GitleaksID: 2, CodeDataSourceID: 1, Name: "two", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: "token", ScanPublic: true, ScanInternal: false, ScanPrivate: false, GitleaksConfig: "conf", Status: "OK", StatusDetail: "detail", ScanAt: now, ScanSucceededAt: &now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty",
			input:     &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			want:      &code.ListGitleaksResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.ListGitleaksRequest{CodeDataSourceId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.ListGitleaksRequest{ProjectId: 1, CodeDataSourceId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListGitleaks").Return(c.mockResponce, c.mockError).Once()
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
		name         string
		input        *code.GetGitleaksRequest
		want         *code.GetGitleaksResponse
		mockResponce *model.CodeGitleaks
		mockError    error
		wantErr      bool
	}{
		{
			name:         "OK",
			input:        &code.GetGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			want:         &code.GetGitleaksResponse{Gitleaks: &code.Gitleaks{GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: "token", ScanPublic: true, ScanInternal: false, ScanPrivate: false, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.CodeGitleaks{GitleaksID: 1, CodeDataSourceID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: "token", ScanPublic: true, ScanInternal: false, ScanPrivate: false, GitleaksConfig: "conf", Status: "OK", StatusDetail: "detail", ScanAt: now, ScanSucceededAt: &now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK empty",
			input:     &code.GetGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			want:      &code.GetGitleaksResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.GetGitleaksRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.GetGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetGitleaks").Return(c.mockResponce, c.mockError).Once()
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
		name         string
		input        *code.PutGitleaksRequest
		want         *code.PutGitleaksResponse
		mockResponce *model.CodeGitleaks
		mockError    error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &code.PutGitleaksRequest{ProjectId: 1, Gitleaks: &code.GitleaksForUpsert{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: now.Unix()},
			},
			want: &code.PutGitleaksResponse{Gitleaks: &code.Gitleaks{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponce: &model.CodeGitleaks{
				GitleaksID: 1, CodeDataSourceID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: "token", GitleaksConfig: "conf", Status: "OK", StatusDetail: "detail", ScanAt: now, ScanSucceededAt: &now, CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name: "OK(empty)",
			input: &code.PutGitleaksRequest{ProjectId: 1, Gitleaks: &code.GitleaksForUpsert{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want: &code.PutGitleaksResponse{Gitleaks: &code.Gitleaks{
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix(), ScanSucceededAt: 0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponce: &model.CodeGitleaks{
				GitleaksID: 1, CodeDataSourceID: 1, Name: "one", ProjectID: 1, Type: "ENTERPRISE", TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: "token", GitleaksConfig: "conf", Status: "OK", StatusDetail: "detail", ScanAt: now, ScanSucceededAt: nil, CreatedAt: now, UpdatedAt: now,
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
				GitleaksId: 1, CodeDataSourceId: 1, Name: "one", ProjectId: 1, Type: code.Type_ENTERPRISE, TargetResource: "target", RepositoryPattern: "repo", GithubUser: "user", PersonalAccessToken: maskData, GitleaksConfig: "conf", Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("UpsertGitleaks").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockdb.MockCodeRepository{}
	svc := CodeService{repository: &mockDB}
	cases := []struct {
		name      string
		input     *code.DeleteGitleaksRequest
		mockError error
		wantErr   bool
	}{
		{
			name:  "OK",
			input: &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
		},
		{
			name:    "NG invalid param",
			input:   &code.DeleteGitleaksRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.DeleteGitleaksRequest{ProjectId: 1, GitleaksId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteGitleaks").Return(c.mockError).Once()
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
		name         string
		input        *code.ListEnterpriseOrgRequest
		want         *code.ListEnterpriseOrgResponse
		mockResponce *[]model.CodeEnterpriseOrg
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &code.ListEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1},
			want: &code.ListEnterpriseOrgResponse{EnterpriseOrg: []*code.EnterpriseOrg{
				{GitleaksId: 1, Login: "one", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{GitleaksId: 2, Login: "two", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.CodeEnterpriseOrg{
				{GitleaksID: 1, Login: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
				{GitleaksID: 2, Login: "two", ProjectID: 1, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty",
			input:     &code.ListEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1},
			want:      &code.ListEnterpriseOrgResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &code.ListEnterpriseOrgRequest{GitleaksId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &code.ListEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListEnterpriseOrg").Return(c.mockResponce, c.mockError).Once()
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
		name         string
		input        *code.PutEnterpriseOrgRequest
		want         *code.PutEnterpriseOrgResponse
		mockResponce *model.CodeEnterpriseOrg
		mockError    error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &code.PutEnterpriseOrgRequest{ProjectId: 1, EnterpriseOrg: &code.EnterpriseOrgForUpsert{
				GitleaksId: 1, Login: "one", ProjectId: 1},
			},
			want: &code.PutEnterpriseOrgResponse{EnterpriseOrg: &code.EnterpriseOrg{
				GitleaksId: 1, Login: "one", ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponce: &model.CodeEnterpriseOrg{
				GitleaksID: 1, Login: "one", ProjectID: 1, CreatedAt: now, UpdatedAt: now,
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
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("UpsertEnterpriseOrg").Return(c.mockResponce, c.mockError).Once()
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
		name      string
		input     *code.DeleteEnterpriseOrgRequest
		mockError error
		wantErr   bool
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
			name:      "Invalid DB error",
			input:     &code.DeleteEnterpriseOrgRequest{ProjectId: 1, GitleaksId: 1, Login: "login"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteEnterpriseOrg").Return(c.mockError).Once()
			_, err := svc.DeleteEnterpriseOrg(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}
