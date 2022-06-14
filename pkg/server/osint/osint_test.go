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

func TestListOsint(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name         string
		input        *osint.ListOsintRequest
		want         *osint.ListOsintResponse
		mockResponce *[]model.Osint
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.ListOsintRequest{ProjectId: 1001},
			want: &osint.ListOsintResponse{Osint: []*osint.Osint{
				{OsintId: 1001, ProjectId: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{OsintId: 1002, ProjectId: 1001, ResourceType: "test_type", ResourceName: "test_name2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.Osint{
				{OsintID: 1001, ProjectID: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now, UpdatedAt: now},
				{OsintID: 1002, ProjectID: 1001, ResourceType: "test_type", ResourceName: "test_name2", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK Record not found",
			input:     &osint.ListOsintRequest{ProjectId: 1001},
			want:      &osint.ListOsintResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListOsint").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListOsint(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetOsint(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name         string
		input        *osint.GetOsintRequest
		want         *osint.GetOsintResponse
		mockResponce *model.Osint
		mockError    error
	}{
		{
			name:  "OK",
			input: &osint.GetOsintRequest{ProjectId: 1001, OsintId: 1001},
			want: &osint.GetOsintResponse{Osint: &osint.Osint{
				OsintId: 1001, ProjectId: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now.Unix(), UpdatedAt: now.Unix(),
			}},
			mockResponce: &model.Osint{
				OsintID: 1001, ProjectID: 1001, ResourceType: "test_type", ResourceName: "test_name", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK Record not found",
			input:     &osint.GetOsintRequest{ProjectId: 1001, OsintId: 1001},
			want:      &osint.GetOsintResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetOsint").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetOsint(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutOsint(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name        string
		input       *osint.PutOsintRequest
		want        *osint.PutOsintResponse
		wantErr     bool
		mockUpdResp *model.Osint
		mockUpdErr  error
	}{
		{
			name:        "OK Update",
			input:       &osint.PutOsintRequest{ProjectId: 1001, Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, OsintId: 1001}},
			want:        &osint.PutOsintResponse{Osint: &osint.Osint{OsintId: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.Osint{OsintID: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Insert",
			input:       &osint.PutOsintRequest{ProjectId: 1001, Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001}},
			want:        &osint.PutOsintResponse{Osint: &osint.Osint{OsintId: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockUpdResp: &model.Osint{OsintID: 1001, ResourceType: "test_type", ResourceName: "test_name", ProjectID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid parameter(required project_id)",
			input:   &osint.PutOsintRequest{Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, OsintId: 1001}},
			wantErr: true,
		},
		{
			name:        "NG DB error(UpsertOsint)",
			input:       &osint.PutOsintRequest{Osint: &osint.OsintForUpsert{ResourceType: "test_type", ResourceName: "test_name", ProjectId: 1001, OsintId: 1001}},
			wantErr:     true,
			mockUpdResp: nil,
			mockUpdErr:  errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockUpdResp != nil || c.mockUpdErr != nil {
				mockDB.On("UpsertOsint").Return(c.mockUpdResp, c.mockUpdErr).Once()
			}
			got, err := svc.PutOsint(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteOsint(t *testing.T) {
	var ctx context.Context
	mockDB := dbmock.MockOsintRepository{}
	svc := OsintService{repository: &mockDB, logger: logging.NewLogger()}
	cases := []struct {
		name                    string
		input                   *osint.DeleteOsintRequest
		wantErr                 bool
		mockResp                error
		mockListOSINTDataSource *[]model.RelOsintDataSource
		mockListOsintDetectWord *[]model.OsintDetectWord
	}{
		{
			name:     "OK",
			input:    &osint.DeleteOsintRequest{ProjectId: 1001, OsintId: 1001},
			wantErr:  false,
			mockResp: nil,
		},
		{
			name:     "NG DB error",
			input:    &osint.DeleteOsintRequest{ProjectId: 1001, OsintId: 1001},
			wantErr:  true,
			mockResp: errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("ListRelOsintDataSource").Return(&[]model.RelOsintDataSource{{RelOsintDataSourceID: 1, ProjectID: 1}}, nil)
			mockDB.On("ListOsintDetectWord").Return(&[]model.OsintDetectWord{{OsintDetectWordID: 1, ProjectID: 1}}, nil)
			mockDB.On("DeleteRelOsintDataSource").Return(nil)
			mockDB.On("DeleteOsintDetectWord").Return(nil)
			mockDB.On("DeleteOsint").Return(c.mockResp).Once()
			_, err := svc.DeleteOsint(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
		})
	}
}
