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
	"github.com/ca-risken/datasource-api/proto/code"
)

func TestListGitHubSetting(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *[]model.CodeGitHubSetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK no param",
			args: args{ProjectID: 0, CodeGitHubSettingID: 0},
			want: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "github_setting1", ProjectID: 1, Type: "USER", TargetResource: "target", CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, Name: "github_setting2", ProjectID: 1, Type: "USER", TargetResource: "target", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_github_setting where 1=1")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "name", "project_id", "type", "target_resource", "created_at", "updated_at"}).
					AddRow(uint32(1), "github_setting1", uint32(1), "USER", "target", now, now).
					AddRow(uint32(2), "github_setting2", uint32(1), "USER", "target", now, now))
			},
		},
		{
			name: "OK (project_id)",
			args: args{ProjectID: 1, CodeGitHubSettingID: 0},
			want: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "github_setting1", ProjectID: 1, Type: "USER", TargetResource: "target", CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, Name: "github_setting2", ProjectID: 1, Type: "USER", TargetResource: "target", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_github_setting where 1=1 and project_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "name", "project_id", "type", "target_resource", "created_at", "updated_at"}).
					AddRow(uint32(1), "github_setting1", uint32(1), "USER", "target", now, now).
					AddRow(uint32(2), "github_setting2", uint32(1), "USER", "target", now, now))
			},
		},
		{
			name: "OK (code_github_setting_id)",
			args: args{ProjectID: 0, CodeGitHubSettingID: 1},
			want: &[]model.CodeGitHubSetting{
				{CodeGitHubSettingID: 1, Name: "github_setting1", ProjectID: 1, Type: "USER", TargetResource: "target", CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, Name: "github_setting2", ProjectID: 1, Type: "USER", TargetResource: "target", CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_github_setting where 1=1 and code_github_setting_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "name", "project_id", "type", "target_resource", "created_at", "updated_at"}).
					AddRow(uint32(1), "github_setting1", uint32(1), "USER", "target", now, now).
					AddRow(uint32(2), "github_setting2", uint32(1), "USER", "target", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_github_setting where 1=1 and project_id = ? and code_github_setting_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListGitHubSetting(ctx, c.args.ProjectID, c.args.CodeGitHubSettingID)
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

func TestGetGitHubSetting(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *model.CodeGitHubSetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    &model.CodeGitHubSetting{CodeGitHubSettingID: 1, Name: "github_setting1", ProjectID: 1, Type: "USER", TargetResource: "target", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_github_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "name", "project_id", "type", "target_resource", "created_at", "updated_at"}).
					AddRow(uint32(1), "github_setting1", uint32(1), "USER", "target", now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_github_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetGitHubSetting(ctx, c.args.ProjectID, c.args.CodeGitHubSettingID)
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

func TestUpsertGitHubSetting(t *testing.T) {
	now := time.Now()
	type args struct {
		data *code.GitHubSettingForUpsert
	}

	cases := []struct {
		name        string
		args        args
		want        *model.CodeGitHubSetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK with token",
			args: args{
				data: &code.GitHubSettingForUpsert{GithubSettingId: 1, Name: "name", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token"},
			},
			want:    &model.CodeGitHubSetting{CodeGitHubSettingID: 1, Name: "github_setting1", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(upsertGitHubWithToken)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetCodeGitHubSettingByUniqueIndex)).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "name", "project_id", "type", "target_resource", "github_user", "personal_access_token", "created_at", "updated_at"}).
					AddRow(uint32(1), "github_setting1", uint32(1), "ORGANIZATION", "target", "user", "token", now, now))
			},
		},
		{
			name: "OK without token",
			args: args{
				data: &code.GitHubSettingForUpsert{GithubSettingId: 1, Name: "name", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user"},
			},
			want:    &model.CodeGitHubSetting{CodeGitHubSettingID: 1, Name: "github_setting1", ProjectID: 1, Type: "ORGANIZATION", TargetResource: "target", GitHubUser: "user", PersonalAccessToken: "token", CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(upsertGitHubSettingWithoutToken)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetCodeGitHubSettingByUniqueIndex)).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "name", "project_id", "type", "target_resource", "github_user", "personal_access_token", "created_at", "updated_at"}).
					AddRow(uint32(1), "github_setting1", uint32(1), "ORGANIZATION", "target", "user", "token", now, now))
			},
		},
		{
			name: "NG DB error",
			args: args{
				data: &code.GitHubSettingForUpsert{GithubSettingId: 1, Name: "name", ProjectId: 1, Type: code.Type_ORGANIZATION, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "token"},
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(upsertGitHubWithToken)).WillReturnError(errors.New("DB error"))
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
			got, err := db.UpsertGitHubSetting(ctx, c.args.data)
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

func TestDeleteGitHubSetting(t *testing.T) {
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `code_github_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `code_github_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
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
			err = db.DeleteGitHubSetting(ctx, c.args.ProjectID, c.args.CodeGitHubSettingID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListGitleaksSetting(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *[]model.CodeGitleaksSetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{ProjectID: 1},
			want: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: true, ScanInternal: true, ScanPrivate: true, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: false, ScanInternal: true, ScanPrivate: false, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_gitleaks_setting where project_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "scan_public", "scan_internal", "scan_private", "status", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), true, true, true, "OK", now, now, now).
					AddRow(uint32(2), uint32(1), uint32(1), false, true, false, "OK", now, now, now))
			},
		},
		{
			name: "OK project_id 0 value",
			args: args{ProjectID: 0},
			want: &[]model.CodeGitleaksSetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: true, ScanInternal: true, ScanPrivate: true, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: false, ScanInternal: true, ScanPrivate: false, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_gitleaks_setting")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "scan_public", "scan_internal", "scan_private", "status", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), true, true, true, "OK", now, now, now).
					AddRow(uint32(2), uint32(1), uint32(1), false, true, false, "OK", now, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_gitleaks_setting where project_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListGitleaksSetting(ctx, c.args.ProjectID)
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

func TestGetGitleaksSetting(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *model.CodeGitleaksSetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    &model.CodeGitleaksSetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: true, ScanInternal: true, ScanPrivate: true, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_gitleaks_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "scan_public", "scan_internal", "scan_private", "status", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), true, true, true, "OK", now, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_gitleaks_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetGitleaksSetting(ctx, c.args.ProjectID, c.args.CodeGitHubSettingID)
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

func TestUpsertGitleaksSetting(t *testing.T) {
	now := time.Now()
	type args struct {
		data *code.GitleaksSettingForUpsert
	}

	cases := []struct {
		name        string
		args        args
		want        *model.CodeGitleaksSetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{
				data: &code.GitleaksSettingForUpsert{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: true, ScanInternal: true, ScanPrivate: true, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want:    &model.CodeGitleaksSetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, ScanPublic: true, ScanInternal: true, ScanPrivate: true, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(upsertGitleaksWithToken)).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_gitleaks_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "scan_public", "scan_internal", "scan_private", "status", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), true, true, true, "OK", now, now, now))
			},
		},
		{
			name: "NG DB error",
			args: args{
				data: &code.GitleaksSettingForUpsert{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", ScanPublic: true, ScanInternal: true, ScanPrivate: true, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(upsertGitleaksWithToken)).WillReturnError(errors.New("DB error"))
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
			got, err := db.UpsertGitleaksSetting(ctx, c.args.data)
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

func TestDeleteGitleaksSetting(t *testing.T) {
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `code_gitleaks_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `code_gitleaks_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
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
			err = db.DeleteGitleaksSetting(ctx, c.args.ProjectID, c.args.CodeGitHubSettingID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

var codeGitleaksCacheTableColumn = []string{"code_github_setting_id", "repository_full_name", "scan_at", "created_at", "updated_at"}

func TestListGitleaksCache(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}
	cases := []struct {
		name        string
		args        *args
		want        *[]model.CodeGitleaksCache
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: &args{ProjectID: 1, CodeGitHubSettingID: 1},
			want: &[]model.CodeGitleaksCache{
				{CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, RepositoryFullName: "owner/repo2", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(selectListGitleaksCache)).
					WillReturnRows(sqlmock.NewRows(codeGitleaksCacheTableColumn).
						AddRow(uint32(1), "owner/repo", now, now, now).
						AddRow(uint32(2), "owner/repo2", now, now, now),
					)
			},
		},
		{
			name:    "NG DB error",
			args:    &args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(selectListGitleaksCache)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%+v' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.ListGitleaksCache(context.TODO(), c.args.ProjectID, c.args.CodeGitHubSettingID)
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

func TestGetGitleaksCache(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
		RepositoryFullName  string
		Immediately         bool
	}
	cases := []struct {
		name        string
		args        *args
		want        *model.CodeGitleaksCache
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK(immediately: true)",
			args:    &args{ProjectID: 1, CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", Immediately: true},
			want:    &model.CodeGitleaksCache{CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetGitleaksCache)).
					WillReturnRows(sqlmock.NewRows(codeGitleaksCacheTableColumn).
						AddRow(uint32(1), "owner/repo", now, now, now),
					)
			},
		},
		{
			name:    "OK(immediately: false)",
			args:    &args{ProjectID: 1, CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", Immediately: false},
			want:    &model.CodeGitleaksCache{CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetGitleaksCache)).
					WillReturnRows(sqlmock.NewRows(codeGitleaksCacheTableColumn).
						AddRow(uint32(1), "owner/repo", now, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    &args{ProjectID: 1, CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo"},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetGitleaksCache)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%+v' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.GetGitleaksCache(context.TODO(), c.args.ProjectID, c.args.CodeGitHubSettingID, c.args.RepositoryFullName, false)
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

func TestUpsertGitleaksCache(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID uint32
		Data      *code.GitleaksCacheForUpsert
	}

	cases := []struct {
		name        string
		args        args
		want        *model.CodeGitleaksCache
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{
				ProjectID: 1,
				Data:      &code.GitleaksCacheForUpsert{GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix()},
			},
			want:    &model.CodeGitleaksCache{CodeGitHubSettingID: 1, RepositoryFullName: "owner/repo", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					regexp.QuoteMeta(upsertGitleaksCache)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(
					regexp.QuoteMeta(selectGetGitleaksCache)).
					WillReturnRows(sqlmock.NewRows(codeGitleaksCacheTableColumn).
						AddRow(uint32(1), "owner/repo", now, now, now))
			},
		},
		{
			name: "NG DB error",
			args: args{
				ProjectID: 1,
				Data:      &code.GitleaksCacheForUpsert{GithubSettingId: 1, RepositoryFullName: "owner/repo", ScanAt: now.Unix()},
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(
					regexp.QuoteMeta(upsertGitleaksCache)).
					WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			c.mockClosure(mock)
			got, err := db.UpsertGitleaksCache(context.TODO(), c.args.ProjectID, c.args.Data)
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

func TestDeleteGitleaksCache(t *testing.T) {
	type args struct {
		CodeGitHubSettingID uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{CodeGitHubSettingID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM .*").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{CodeGitHubSettingID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM .*").
					WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
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
			err = db.DeleteGitleaksCache(ctx, c.args.CodeGitHubSettingID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListDependencySetting(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *[]model.CodeDependencySetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{ProjectID: 1},
			want: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_dependency_setting where project_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "status", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "OK", now, now, now).
					AddRow(uint32(2), uint32(1), uint32(1), "OK", now, now, now))
			},
		},
		{
			name: "OK project_id 0 value",
			args: args{ProjectID: 0},
			want: &[]model.CodeDependencySetting{
				{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
				{CodeGitHubSettingID: 2, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_dependency_setting")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "status", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "OK", now, now, now).
					AddRow(uint32(2), uint32(1), uint32(1), "OK", now, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select * from code_dependency_setting where project_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListDependencySetting(ctx, c.args.ProjectID)
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

func TestGetDependencySetting(t *testing.T) {
	now := time.Now()
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}

	cases := []struct {
		name        string
		args        args
		want        *model.CodeDependencySetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    &model.CodeDependencySetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", ScanAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "status", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "OK", now, now, now))
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnError(errors.New("DB error"))
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
			got, err := db.GetDependencySetting(ctx, c.args.ProjectID, c.args.CodeGitHubSettingID)
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

func TestUpsertDependencySetting(t *testing.T) {
	now := time.Now()
	type args struct {
		data *code.DependencySettingForUpsert
	}

	cases := []struct {
		name        string
		args        args
		want        *model.CodeDependencySetting
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK update",
			args: args{
				data: &code.DependencySettingForUpsert{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want:    &model.CodeDependencySetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "detail", ScanAt: now, ErrorNotifiedAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ? ORDER BY `code_dependency_setting`.`code_github_setting_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "status", "status_detail", "scan_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "OK", "detail", now, now, now))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `code_dependency_setting` SET `code_data_source_id`=?,`code_github_setting_id`=?,`project_id`=?,`scan_at`=?,`status`=?,`status_detail`=?,`updated_at`=? WHERE (project_id = ? AND code_github_setting_id = ?) AND `code_github_setting_id` = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ? ORDER BY `code_dependency_setting`.`code_github_setting_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "status", "status_detail", "scan_at", "error_notified_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "OK", "detail", now, now, now, now))
			},
		},
		{
			name: "OK insert",
			args: args{
				data: &code.DependencySettingForUpsert{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want:    &model.CodeDependencySetting{CodeGitHubSettingID: 1, CodeDataSourceID: 1, ProjectID: 1, Status: "OK", StatusDetail: "detail", ScanAt: now, ErrorNotifiedAt: now, CreatedAt: now, UpdatedAt: now},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ? ORDER BY `code_dependency_setting`.`code_github_setting_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "status", "scan_at", "created_at", "updated_at"}))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `code_dependency_setting` (`code_data_source_id`,`project_id`,`status`,`status_detail`,`scan_at`,`created_at`,`updated_at`,`code_github_setting_id`) VALUES (?,?,?,?,?,?,?,?)")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ? ORDER BY `code_dependency_setting`.`code_github_setting_id` LIMIT 1")).WillReturnRows(sqlmock.NewRows([]string{
					"code_github_setting_id", "code_data_source_id", "project_id", "status", "status_detail", "scan_at", "error_notified_at", "created_at", "updated_at"}).
					AddRow(uint32(1), uint32(1), uint32(1), "OK", "detail", now, now, now, now))
			},
		},
		{
			name: "NG DB error",
			args: args{
				data: &code.DependencySettingForUpsert{GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: code.Status_OK, StatusDetail: "detail", ScanAt: now.Unix()},
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ? ORDER BY `code_dependency_setting`.`code_github_setting_id` LIMIT 1")).WillReturnError(errors.New("DB error"))
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
			got, err := db.UpsertDependencySetting(ctx, c.args.data)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			// 自動生成されるタイムスタンプをwantで指定できないのでそれ以外の値を比較
			if c.want != nil && !((got.CodeGitHubSettingID == c.want.CodeGitHubSettingID) && (got.CodeDataSourceID == c.want.CodeDataSourceID) && (got.ProjectID == c.want.ProjectID) && (got.Status == c.want.Status) && (got.StatusDetail == c.want.StatusDetail)) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteDependencySetting(t *testing.T) {
	type args struct {
		ProjectID           uint32
		CodeGitHubSettingID uint32
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "NG DB error",
			args:    args{ProjectID: 1, CodeGitHubSettingID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `code_dependency_setting` WHERE project_id = ? AND code_github_setting_id = ?")).WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
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
			err = db.DeleteDependencySetting(ctx, c.args.ProjectID, c.args.CodeGitHubSettingID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListCodeGitleaksScanErrorForNotify(t *testing.T) {
	db, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	cases := []struct {
		name        string
		want        []*GitHubScanError
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			want: []*GitHubScanError{
				{CodeGithubSettingID: 1, ProjectID: 1, DataSource: "code:gitleaks", StatusDetail: "detail"},
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListCodeGitHubScanError)).WillReturnRows(
					sqlmock.NewRows([]string{"code_github_setting_id", "project_id", "data_source", "status_detail"}).
						AddRow(1, 1, "code:gitleaks", "detail"))
			},
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListCodeGitHubScanError)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			c.mockClosure(mock)
			got, err := db.ListCodeGitHubScanErrorForNotify(ctx)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateCodeGitleaksErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	db, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	type args struct {
		errNotifiedAt       interface{}
		codeGithubSettingID uint32
		projectID           uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{errNotifiedAt: now, projectID: 1, codeGithubSettingID: 1},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateCodeGitleaksErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, projectID: 1, codeGithubSettingID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateCodeGitleaksErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			c.mockClosure(mock)
			err = db.UpdateCodeGitleaksErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.codeGithubSettingID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateCodeDependencyErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	db, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	type args struct {
		errNotifiedAt       interface{}
		codeGithubSettingID uint32
		projectID           uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{errNotifiedAt: now, projectID: 1, codeGithubSettingID: 1},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateCodeDependencyErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, projectID: 1, codeGithubSettingID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateCodeDependencyErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			c.mockClosure(mock)
			err = db.UpdateCodeDependencyErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.codeGithubSettingID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
