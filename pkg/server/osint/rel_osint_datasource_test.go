package osint

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	dbmock "github.com/ca-risken/datasource-api/pkg/db/mock"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/osint"
	"gorm.io/gorm"
)

func TestListRelOsintDataSource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name         string
		input        *osint.ListRelOsintDataSourceRequest
		want         *osint.ListRelOsintDataSourceResponse
		mockResponce *[]model.RelOsintDataSource
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.ListRelOsintDataSourceRequest{ProjectId: 1001},
			want: &osint.ListRelOsintDataSourceResponse{RelOsintDataSource: []*osint.RelOsintDataSource{
				{RelOsintDataSourceId: 1001, OsintId: 1001, OsintDataSourceId: 1001, Status: osint.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{RelOsintDataSourceId: 1002, OsintId: 1002, OsintDataSourceId: 1001, Status: osint.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.RelOsintDataSource{
				{RelOsintDataSourceID: 1001, OsintID: 1001, OsintDataSourceID: 1001, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{RelOsintDataSourceID: 1002, OsintID: 1002, OsintDataSourceID: 1001, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK Record not found",
			input:     &osint.ListRelOsintDataSourceRequest{ProjectId: 1001},
			want:      &osint.ListRelOsintDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListRelOsintDataSource").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListRelOsintDataSource(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetRelOsintDataSource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name         string
		input        *osint.GetRelOsintDataSourceRequest
		want         *osint.GetRelOsintDataSourceResponse
		mockResponce *model.RelOsintDataSource
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.GetRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSourceId: 1001},
			want: &osint.GetRelOsintDataSourceResponse{RelOsintDataSource: &osint.RelOsintDataSource{
				RelOsintDataSourceId: 1001, OsintId: 1001, OsintDataSourceId: 1001, Status: osint.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
			}},
			mockResponce: &model.RelOsintDataSource{RelOsintDataSourceID: 1001, OsintID: 1001, OsintDataSourceID: 1001, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record not found",
			input:     &osint.GetRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSourceId: 1001},
			want:      &osint.GetRelOsintDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetRelOsintDataSource").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetRelOsintDataSource(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutRelOsintDataSource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name        string
		input       *osint.PutRelOsintDataSourceRequest
		want        *osint.PutRelOsintDataSourceResponse
		wantErr     bool
		mockUpdResp *model.RelOsintDataSource
		mockUpdErr  error
	}{
		{
			name:        "OK Update",
			input:       &osint.PutRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSource: &osint.RelOsintDataSourceForUpsert{RelOsintDataSourceId: 1001, OsintId: 1001, ProjectId: 1001, OsintDataSourceId: 1001, Status: osint.Status_OK, ScanAt: now.Unix()}},
			want:        &osint.PutRelOsintDataSourceResponse{RelOsintDataSource: &osint.RelOsintDataSource{RelOsintDataSourceId: 1001, OsintId: 1001, OsintDataSourceId: 1001, Status: osint.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.RelOsintDataSource{RelOsintDataSourceID: 1001, OsintID: 1001, OsintDataSourceID: 1001, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Insert",
			input:       &osint.PutRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSource: &osint.RelOsintDataSourceForUpsert{OsintId: 1001, OsintDataSourceId: 1001, ProjectId: 1001, Status: osint.Status_OK, ScanAt: now.Unix()}},
			want:        &osint.PutRelOsintDataSourceResponse{RelOsintDataSource: &osint.RelOsintDataSource{RelOsintDataSourceId: 1001, OsintId: 1001, OsintDataSourceId: 1001, Status: osint.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.RelOsintDataSource{RelOsintDataSourceID: 1001, OsintID: 1001, OsintDataSourceID: 1001, Status: "OK", StatusDetail: "", ScanAt: now, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter(required project_id)",
			input:   &osint.PutRelOsintDataSourceRequest{RelOsintDataSource: &osint.RelOsintDataSourceForUpsert{OsintId: 1001, OsintDataSourceId: 1001, ProjectId: 1001, Status: osint.Status_OK, ScanAt: now.Unix()}},
			wantErr: true,
		},
		{
			name:        "NG DB error(UpsertRelOsintDataSource)",
			input:       &osint.PutRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSource: &osint.RelOsintDataSourceForUpsert{RelOsintDataSourceId: 1001, OsintId: 1001, ProjectId: 1001, OsintDataSourceId: 1001, Status: osint.Status_OK, ScanAt: now.Unix()}},
			mockUpdResp: nil,
			mockUpdErr:  errors.New("Something error"),
			wantErr:     true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mockDB.On("UpsertRelOsintDataSource").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutRelOsintDataSource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteRelOsintDataSource(t *testing.T) {
	var ctx context.Context
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name     string
		input    *osint.DeleteRelOsintDataSourceRequest
		wantErr  bool
		mockResp error
	}{
		{
			name:    "OK",
			input:   &osint.DeleteRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSourceId: 1001},
			wantErr: false,
		},
		{
			name:     "Invalid DB error",
			input:    &osint.DeleteRelOsintDataSourceRequest{ProjectId: 1001, RelOsintDataSourceId: 1001},
			wantErr:  true,
			mockResp: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("ListOsintDetectWord").Return(&[]model.OsintDetectWord{{OsintDetectWordID: 1, ProjectID: 1}}, nil)
			mockDB.On("DeleteOsintDetectWord").Return(nil)
			mockDB.On("DeleteRelOsintDataSource").Return(c.mockResp).Once()
			_, err := svc.DeleteRelOsintDataSource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
		})
	}
}
