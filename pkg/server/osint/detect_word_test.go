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

func TestListOsintDetectWord(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name         string
		input        *osint.ListOsintDetectWordRequest
		want         *osint.ListOsintDetectWordResponse
		mockResponce *[]model.OsintDetectWord
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.ListOsintDetectWordRequest{ProjectId: 1001},
			want: &osint.ListOsintDetectWordResponse{OsintDetectWord: []*osint.OsintDetectWord{
				{OsintDetectWordId: 1001, RelOsintDataSourceId: 1001, Word: "hoge", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{OsintDetectWordId: 1002, RelOsintDataSourceId: 1002, Word: "hoge", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.OsintDetectWord{
				{OsintDetectWordID: 1001, RelOsintDataSourceID: 1001, Word: "hoge", CreatedAt: now, UpdatedAt: now},
				{OsintDetectWordID: 1002, RelOsintDataSourceID: 1002, Word: "hoge", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK Record not found",
			input:     &osint.ListOsintDetectWordRequest{ProjectId: 1001},
			want:      &osint.ListOsintDetectWordResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListOsintDetectWord").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOsintDetectWord(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOsintDetectWord(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name         string
		input        *osint.GetOsintDetectWordRequest
		want         *osint.GetOsintDetectWordResponse
		mockResponce *model.OsintDetectWord
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.GetOsintDetectWordRequest{ProjectId: 1001, OsintDetectWordId: 1001},
			want: &osint.GetOsintDetectWordResponse{OsintDetectWord: &osint.OsintDetectWord{
				OsintDetectWordId: 1001, RelOsintDataSourceId: 1001, Word: "hoge", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
			}},
			mockResponce: &model.OsintDetectWord{OsintDetectWordID: 1001, RelOsintDataSourceID: 1001, Word: "hoge", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record not found",
			input:     &osint.GetOsintDetectWordRequest{ProjectId: 1001, OsintDetectWordId: 1001},
			want:      &osint.GetOsintDetectWordResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetOsintDetectWord").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOsintDetectWord(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOsintDetectWord(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name        string
		input       *osint.PutOsintDetectWordRequest
		want        *osint.PutOsintDetectWordResponse
		wantErr     bool
		mockUpdResp *model.OsintDetectWord
		mockUpdErr  error
	}{
		{
			name:        "OK Update",
			input:       &osint.PutOsintDetectWordRequest{ProjectId: 1001, OsintDetectWord: &osint.OsintDetectWordForUpsert{OsintDetectWordId: 1001, RelOsintDataSourceId: 1001, ProjectId: 1001, Word: "hoge"}},
			want:        &osint.PutOsintDetectWordResponse{OsintDetectWord: &osint.OsintDetectWord{OsintDetectWordId: 1001, RelOsintDataSourceId: 1001, Word: "hoge", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.OsintDetectWord{OsintDetectWordID: 1001, RelOsintDataSourceID: 1001, Word: "hoge", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Insert",
			input:       &osint.PutOsintDetectWordRequest{ProjectId: 1001, OsintDetectWord: &osint.OsintDetectWordForUpsert{RelOsintDataSourceId: 1001, ProjectId: 1001, Word: "hoge"}},
			want:        &osint.PutOsintDetectWordResponse{OsintDetectWord: &osint.OsintDetectWord{OsintDetectWordId: 1001, RelOsintDataSourceId: 1001, Word: "hoge", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.OsintDetectWord{OsintDetectWordID: 1001, RelOsintDataSourceID: 1001, Word: "hoge", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter(required project_id)",
			input:   &osint.PutOsintDetectWordRequest{OsintDetectWord: &osint.OsintDetectWordForUpsert{RelOsintDataSourceId: 1001, ProjectId: 1001, Word: "hoge"}},
			wantErr: true,
		},
		{
			name:        "NG DB error(UpsertOsintDetectWord)",
			input:       &osint.PutOsintDetectWordRequest{ProjectId: 1001, OsintDetectWord: &osint.OsintDetectWordForUpsert{OsintDetectWordId: 1001, RelOsintDataSourceId: 1001, Word: "hoge", ProjectId: 1001}},
			mockUpdResp: nil,
			mockUpdErr:  errors.New("Something error"),
			wantErr:     true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mockDB.On("UpsertOsintDetectWord").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOsintDetectWord(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOsintDetectWord(t *testing.T) {
	var ctx context.Context
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name     string
		input    *osint.DeleteOsintDetectWordRequest
		wantErr  bool
		mockResp error
	}{
		{
			name:    "OK",
			input:   &osint.DeleteOsintDetectWordRequest{ProjectId: 1001, OsintDetectWordId: 1001},
			wantErr: false,
		},
		{
			name:     "Invalid DB error",
			input:    &osint.DeleteOsintDetectWordRequest{ProjectId: 1001, OsintDetectWordId: 1001},
			wantErr:  true,
			mockResp: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteOsintDetectWord").Return(c.mockResp).Once()
			_, err := svc.DeleteOsintDetectWord(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
		})
	}
}
