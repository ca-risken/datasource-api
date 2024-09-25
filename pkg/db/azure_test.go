package db

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/azure"
)

func TestListAzureDataSource(t *testing.T) {
	now := time.Now()
	type args struct {
		AzureDataSourceID uint32
		Name              string
	}
	cases := []struct {
		name        string
		args        args
		want        *[]model.AzureDataSource
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK no param",
			args: args{},
			want: &[]model.AzureDataSource{
				{AzureDataSourceID: 1, Name: "azure:datasource1", Description: "description", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
				{AzureDataSourceID: 2, Name: "azure:datasource2", Description: "description", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure_data_source where 1=1")).WillReturnRows(sqlmock.NewRows([]string{
					"azure_data_source_id", "name", "description", "max_score", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure:datasource1", "description", 1.0, now, now).
					AddRow(uint32(2), "azure:datasource2", "description", 1.0, now, now))
			},
		},
		{
			name: "OK (azure_data_source_id)",
			args: args{AzureDataSourceID: 1},
			want: &[]model.AzureDataSource{
				{AzureDataSourceID: 1, Name: "azure:datasource1", Description: "description", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure_data_source where 1=1 and azure_data_source_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"azure_data_source_id", "name", "description", "max_score", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure:datasource1", "description", 1.0, now, now))
			},
		},
		{
			name: "OK (name)",
			args: args{Name: "azure:datasource1"},
			want: &[]model.AzureDataSource{
				{AzureDataSourceID: 1, Name: "azure:datasource1", Description: "description", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure_data_source where 1=1 and name = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"azure_data_source_id", "name", "description", "max_score", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure:datasource1", "description", 1.0, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure_data_source where 1=1")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListAzureDataSource(ctx, c.args.AzureDataSourceID, c.args.Name)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListAzure(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID      uint32
		AzureID        uint32
		SubscriptionID string
	}

	cases := []struct {
		name        string
		args        args
		want        *[]model.Azure
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK no param",
			args: args{},
			want: &[]model.Azure{
				{AzureID: 1, Name: "azure1", ProjectID: 1, SubscriptionID: "azure1", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
				{AzureID: 2, Name: "azure2", ProjectID: 1, SubscriptionID: "azure2", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure where 1=1")).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now).
					AddRow(uint32(2), "azure2", uint32(1), "azure2", "code", now, now))
			},
		},
		{
			name: "OK (project_id)",
			args: args{ProjectID: 1},
			want: &[]model.Azure{
				{AzureID: 1, Name: "azure1", ProjectID: 1, SubscriptionID: "azure1", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
				{AzureID: 2, Name: "azure2", ProjectID: 1, SubscriptionID: "azure2", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure where 1=1 and project_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now).
					AddRow(uint32(2), "azure2", uint32(1), "azure2", "code", now, now))
			},
		},
		{
			name: "OK (azure_id)",
			args: args{AzureID: 1},
			want: &[]model.Azure{
				{AzureID: 1, Name: "azure1", ProjectID: 1, SubscriptionID: "azure1", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
				{AzureID: 2, Name: "azure2", ProjectID: 1, SubscriptionID: "azure2", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure where 1=1 and azure_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now).
					AddRow(uint32(2), "azure2", uint32(1), "azure2", "code", now, now))
			},
		},
		{
			name: "OK (subscription_id)",
			args: args{SubscriptionID: "1234567890"},
			want: &[]model.Azure{
				{AzureID: 1, Name: "azure1", ProjectID: 1, SubscriptionID: "azure1", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure where 1=1 and subscription_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, AzureID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from azure where 1=1 and project_id = ? and azure_id = ?")).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListAzure(ctx, c.args.ProjectID, c.args.AzureID, c.args.SubscriptionID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAzure(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID uint32
		AzureID   uint32
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Azure
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{AzureID: 1, ProjectID: 1},
			want:    &model.Azure{AzureID: 1, Name: "azure1", ProjectID: 1, SubscriptionID: "azure1", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetAzure)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
					AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{AzureID: 1, ProjectID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetAzure)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.GetAzure(ctx, c.args.ProjectID, c.args.AzureID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpsertAzure(t *testing.T) {
	now := time.Now()
	type args struct {
		azure *azure.AzureForUpsert
	}
	cases := []struct {
		name        string
		args        args
		want        *model.Azure
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{azure: &azure.AzureForUpsert{AzureId: 1, Name: "azure1", ProjectId: 1, SubscriptionId: "azure1", VerificationCode: "code"}},
			want:    &model.Azure{AzureID: 1, Name: "azure1", ProjectID: 1, SubscriptionID: "azure1", VerificationCode: "code", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					regexp.QuoteMeta(insertUpsertAzure)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetAzureBySubscriptionID)).
					WillReturnRows(sqlmock.NewRows([]string{
						"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
						AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{azure: &azure.AzureForUpsert{Name: "azure1", ProjectId: 1, SubscriptionId: "azure1", VerificationCode: "code"}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					regexp.QuoteMeta(insertUpsertAzure)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.UpsertAzure(ctx, c.args.azure)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteAzure(t *testing.T) {
	type args struct {
		ProjectID uint32
		AzureID   uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{AzureID: 1, ProjectID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteAzure)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{AzureID: 1, ProjectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteAzure)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.DeleteAzure(ctx, c.args.ProjectID, c.args.AzureID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListRelAzureDataSource(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID uint32
		AzureID   uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *[]RelAzureDataSource
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK no param",
			args: args{},
			want: &[]RelAzureDataSource{
				{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: "status", StatusDetail: "status_detail",
					ScanAt: now, CreatedAt: now, UpdatedAt: now, Name: "azure:datasource1", Description: "description", MaxScore: 1.0,
					SubscriptionID: "0123456789", ErrorNotifiedAt: now,
				},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListRelAzureDataSource)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "azure_data_source_id", "project_id", "status", "status_detail", "scan_at", "created_at", "updated_at",
					"name", "description", "max_score", "subscription_id", "error_notified_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "status", "status_detail", now, now, now, "azure:datasource1", "description", 1.0, "0123456789", now))
			},
		},
		{
			name: "OK (project_id)",
			args: args{ProjectID: 1},
			want: &[]RelAzureDataSource{
				{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: "status", StatusDetail: "status_detail",
					ScanAt: now, CreatedAt: now, UpdatedAt: now, Name: "azure:datasource1", Description: "description", MaxScore: 1.0,
					SubscriptionID: "0123456789", ErrorNotifiedAt: now,
				},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListRelAzureDataSource)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "azure_data_source_id", "project_id", "status", "status_detail", "scan_at", "created_at", "updated_at",
					"name", "description", "max_score", "subscription_id", "error_notified_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "status", "status_detail", now, now, now, "azure:datasource1", "description", 1.0, "0123456789", now))
			},
		},
		{
			name: "OK (azure_id)",
			args: args{AzureID: 1},
			want: &[]RelAzureDataSource{
				{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: "status", StatusDetail: "status_detail",
					ScanAt: now, CreatedAt: now, UpdatedAt: now, Name: "azure:datasource1", Description: "description", MaxScore: 1.0,
					SubscriptionID: "0123456789", ErrorNotifiedAt: now,
				},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListRelAzureDataSource)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "azure_data_source_id", "project_id", "status", "status_detail", "scan_at", "created_at", "updated_at",
					"name", "description", "max_score", "subscription_id", "error_notified_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "status", "status_detail", now, now, now, "azure:datasource1", "description", 1.0, "0123456789", now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, AzureID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListRelAzureDataSource)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListRelAzureDataSource(ctx, c.args.ProjectID, c.args.AzureID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetRelAzureDataSource(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID         uint32
		AzureID           uint32
		AzureDataSourceID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *RelAzureDataSource
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{ProjectID: 1, AzureID: 1, AzureDataSourceID: 1},
			want: &RelAzureDataSource{
				AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: "status", StatusDetail: "status_detail",
				ScanAt: now, CreatedAt: now, UpdatedAt: now, Name: "azure:datasource1", Description: "description", MaxScore: 1.0,
				SubscriptionID: "0123456789", ErrorNotifiedAt: now,
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetRelAzureDataSource)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "azure_data_source_id", "project_id", "status", "status_detail", "scan_at", "created_at", "updated_at",
					"name", "description", "max_score", "subscription_id", "error_notified_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "status", "status_detail", now, now, now, "azure:datasource1", "description", 1.0, "0123456789", now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, AzureID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetRelAzureDataSource)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.GetRelAzureDataSource(ctx, c.args.ProjectID, c.args.AzureID, c.args.AzureDataSourceID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpsertRelAzureDataSource(t *testing.T) {
	now := time.Now()
	type args struct {
		azure *azure.RelAzureDataSourceForUpsert
	}
	cases := []struct {
		name        string
		args        args
		want        *RelAzureDataSource
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{azure: &azure.RelAzureDataSourceForUpsert{AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_CONFIGURED, StatusDetail: "status_detail", ScanAt: now.Unix()}},
			want: &RelAzureDataSource{
				AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: "configured", StatusDetail: "status_detail",
				Name: "azure1", SubscriptionID: "azure1", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetAzureDataSource)).
					WillReturnRows(sqlmock.NewRows([]string{
						"azure_datasource_id", "name", "description", "max_score", "created_at", "updated_at"}).
						AddRow(uint32(1), "azure:datasource1", "description", 1.0, now, now))
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetAzure)).
					WillReturnRows(sqlmock.NewRows([]string{
						"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
						AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now))
				mock.ExpectExec(
					regexp.QuoteMeta(insertUpsertRelAzureDataSource)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetRelAzureDataSource)).
					WillReturnRows(sqlmock.NewRows([]string{
						"azure_id", "azure_data_source_id", "name", "project_id", "status", "status_detail", "subscription_id", "verification_code", "scan_at", "created_at", "updated_at"}).
						AddRow(uint32(1), uint32(1), "azure1", uint32(1), "configured", "status_detail", "azure1", "code", now, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{azure: &azure.RelAzureDataSourceForUpsert{AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_CONFIGURED, StatusDetail: "status_detail", ScanAt: now.Unix()}},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetAzureDataSource)).
					WillReturnRows(sqlmock.NewRows([]string{
						"azure_datasource_id", "name", "description", "max_score", "created_at", "updated_at"}).
						AddRow(uint32(1), "azure:datasource1", "description", 1.0, now, now))
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetAzure)).
					WillReturnRows(sqlmock.NewRows([]string{
						"azure_id", "name", "project_id", "subscription_id", "verification_code", "created_at", "updated_at"}).
						AddRow(uint32(1), "azure1", uint32(1), "azure1", "code", now, now))
				mock.ExpectExec(
					regexp.QuoteMeta(insertUpsertRelAzureDataSource)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.UpsertRelAzureDataSource(ctx, c.args.azure)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteRelAzureDataSource(t *testing.T) {
	type args struct {
		ProjectID         uint32
		AzureID           uint32
		AzureDataSourceID uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{AzureID: 1, ProjectID: 1, AzureDataSourceID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteRelAzureDataSource)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{AzureID: 1, ProjectID: 1, AzureDataSourceID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteRelAzureDataSource)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.DeleteRelAzureDataSource(ctx, c.args.ProjectID, c.args.AzureID, c.args.AzureDataSourceID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListRelAzureDataSourceByDataSourceID(t *testing.T) {
	now := time.Now()
	type args struct {
		AzureDataSourceID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *[]RelAzureDataSource
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK no param",
			args: args{},
			want: &[]RelAzureDataSource{
				{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: "status", StatusDetail: "status_detail",
					ScanAt: now, CreatedAt: now, UpdatedAt: now, Name: "azure:datasource1", Description: "description", MaxScore: 1.0,
					SubscriptionID: "0123456789", ErrorNotifiedAt: now,
				},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListRelAzureDataSourceByDataSourceID)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "azure_data_source_id", "project_id", "status", "status_detail", "scan_at", "created_at", "updated_at",
					"name", "description", "max_score", "subscription_id", "error_notified_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "status", "status_detail", now, now, now, "azure:datasource1", "description", 1.0, "0123456789", now))
			},
		},
		{
			name: "OK (azure_data_source_id)",
			args: args{AzureDataSourceID: 1},
			want: &[]RelAzureDataSource{
				{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: "status", StatusDetail: "status_detail",
					ScanAt: now, CreatedAt: now, UpdatedAt: now, Name: "azure:datasource1", Description: "description", MaxScore: 1.0,
					SubscriptionID: "0123456789", ErrorNotifiedAt: now,
				},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListRelAzureDataSourceByDataSourceID)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "azure_data_source_id", "project_id", "status", "status_detail", "scan_at", "created_at", "updated_at",
					"name", "description", "max_score", "subscription_id", "error_notified_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "status", "status_detail", now, now, now, "azure:datasource1", "description", 1.0, "0123456789", now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{AzureDataSourceID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListRelAzureDataSourceByDataSourceID)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListRelAzureDataSourceByDataSourceID(ctx, c.args.AzureDataSourceID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListAzureScanErrorForNotify(t *testing.T) {
	cases := []struct {
		name        string
		want        []*AzureScanError
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			want: []*AzureScanError{
				{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, DataSource: "azure:data_source1", StatusDetail: "status_detail"},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListAzureScanError)).WillReturnRows(sqlmock.NewRows([]string{
					"azure_id", "azure_data_source_id", "project_id", "data_source", "status_detail"}).
					AddRow(uint32(1), uint32(1), uint32(1), "azure:data_source1", "status_detail"))
			},
		},
		{
			name:    "NG DB error",
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListAzureScanError)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListAzureScanErrorForNotify(ctx)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateAzureErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	type args struct {
		ErrorNotifiedAt   time.Time
		AzureID           uint32
		AzureDataSourceID uint32
		ProjectID         uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{ErrorNotifiedAt: now, AzureID: 1, AzureDataSourceID: 1, ProjectID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					regexp.QuoteMeta(updateAzureErrorNotifiedAt)).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ErrorNotifiedAt: now, AzureID: 1, AzureDataSourceID: 1, ProjectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					regexp.QuoteMeta(updateAzureErrorNotifiedAt)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			err = db.UpdateAzureErrorNotifiedAt(ctx, c.args.ErrorNotifiedAt, c.args.AzureID, c.args.AzureDataSourceID, c.args.ProjectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
