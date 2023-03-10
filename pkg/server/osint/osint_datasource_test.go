package osint

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/pkg/db/mocks"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/pkg/test"
	"github.com/ca-risken/datasource-api/proto/osint"
	"gorm.io/gorm"
)

func TestListOsintDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *osint.ListOsintDataSourceRequest
		want         *osint.ListOsintDataSourceResponse
		mockResponce *[]model.OsintDataSource
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.ListOsintDataSourceRequest{ProjectId: 1001, Name: "test"},
			want: &osint.ListOsintDataSourceResponse{OsintDataSource: []*osint.OsintDataSource{
				{OsintDataSourceId: 1001, Name: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{OsintDataSourceId: 1002, Name: "test_name2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.OsintDataSource{
				{OsintDataSourceID: 1001, Name: "test_name", CreatedAt: now, UpdatedAt: now},
				{OsintDataSourceID: 1002, Name: "test_name2", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK Record not found",
			input:     &osint.ListOsintDataSourceRequest{ProjectId: 1001, Name: "test"},
			want:      &osint.ListOsintDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOSINTRepoInterface(t)
			svc := OsintService{repository: mockDB, logger: logging.NewLogger()}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListOsintDataSource", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOsintDataSource(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOsintDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *osint.GetOsintDataSourceRequest
		want         *osint.GetOsintDataSourceResponse
		mockResponce *model.OsintDataSource
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.GetOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			want: &osint.GetOsintDataSourceResponse{OsintDataSource: &osint.OsintDataSource{
				OsintDataSourceId: 1001, Name: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
			}},
			mockResponce: &model.OsintDataSource{
				OsintDataSourceID: 1001, Name: "test_name", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record not found",
			input:     &osint.GetOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			want:      &osint.GetOsintDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOSINTRepoInterface(t)
			svc := OsintService{repository: mockDB, logger: logging.NewLogger()}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetOsintDataSource", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOsintDataSource(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOsintDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name        string
		input       *osint.PutOsintDataSourceRequest
		want        *osint.PutOsintDataSourceResponse
		wantErr     bool
		mockUpdResp *model.OsintDataSource
		mockUpdErr  error
	}{
		{
			name:        "OK Update",
			input:       &osint.PutOsintDataSourceRequest{ProjectId: 1001, OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0, OsintDataSourceId: 1001}},
			want:        &osint.PutOsintDataSourceResponse{OsintDataSource: &osint.OsintDataSource{OsintDataSourceId: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.OsintDataSource{OsintDataSourceID: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Insert",
			input:       &osint.PutOsintDataSourceRequest{ProjectId: 1001, OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0}},
			want:        &osint.PutOsintDataSourceResponse{OsintDataSource: &osint.OsintDataSource{OsintDataSourceId: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.OsintDataSource{OsintDataSourceID: 1001, Name: "test_name", Description: "test_desc", MaxScore: 10.0, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter(required project_id)",
			input:   &osint.PutOsintDataSourceRequest{OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0, OsintDataSourceId: 1001}},
			wantErr: true,
		},
		{
			name:        "NG DB error(UpsertOsintDataSource)",
			input:       &osint.PutOsintDataSourceRequest{ProjectId: 1001, OsintDataSource: &osint.OsintDataSourceForUpsert{Name: "test_name", Description: "test_desc", MaxScore: 10.0, OsintDataSourceId: 1001}},
			mockUpdResp: nil,
			mockUpdErr:  errors.New("Something error"),
			wantErr:     true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOSINTRepoInterface(t)
			svc := OsintService{repository: mockDB, logger: logging.NewLogger()}

			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mockDB.On("UpsertOsintDataSource", test.RepeatMockAnything(2)...).Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOsintDataSource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOsintDataSource(t *testing.T) {
	cases := []struct {
		name     string
		input    *osint.DeleteOsintDataSourceRequest
		wantErr  bool
		mockCall bool
		mockResp error
	}{
		{
			name:     "OK",
			input:    &osint.DeleteOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			wantErr:  false,
			mockCall: true,
		},
		{
			name:     "Invalid DB error",
			input:    &osint.DeleteOsintDataSourceRequest{ProjectId: 1001, OsintDataSourceId: 1001},
			wantErr:  true,
			mockCall: true,
			mockResp: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewOSINTRepoInterface(t)
			svc := OsintService{repository: mockDB, logger: logging.NewLogger()}

			if c.mockCall {
				mockDB.On("DeleteOsintDataSource", test.RepeatMockAnything(3)...).Return(c.mockResp).Once()
			}
			_, err := svc.DeleteOsintDataSource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
		})
	}
}
