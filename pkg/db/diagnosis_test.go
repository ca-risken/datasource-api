package db

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListDiagnosisScanErrorForNotify(t *testing.T) {
	db, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	cases := []struct {
		name        string
		want        []*DiagnosisScanError
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			want: []*DiagnosisScanError{
				{ScanID: 1, ProjectID: 1, DataSource: "diagnosis:wpscan", StatusDetail: "detail"},
			},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListDiagnosisScanErrorForNotify)).WillReturnRows(
					sqlmock.NewRows([]string{"scan_id", "project_id", "data_source", "status_detail"}).
						AddRow(1, 1, "diagnosis:wpscan", "detail"))
			},
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListDiagnosisScanErrorForNotify)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			c.mockClosure(mock)
			got, err := db.ListDiagnosisScanErrorForNotify(ctx)
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

func TestUpdateDiagnosisWpscanErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	db, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	type args struct {
		errNotifiedAt interface{}
		scanID        uint32
		projectID     uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{errNotifiedAt: now, projectID: 1, scanID: 1},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateDiagnosisWpscanErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, projectID: 1, scanID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateDiagnosisWpscanErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			c.mockClosure(mock)
			err = db.UpdateDiagnosisWpscanErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.scanID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateDiagnosisPortscanErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	db, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	type args struct {
		errNotifiedAt interface{}
		scanID        uint32
		projectID     uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{errNotifiedAt: now, projectID: 1, scanID: 1},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateDiagnosisPortscanErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, projectID: 1, scanID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateDiagnosisPortscanErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			c.mockClosure(mock)
			err = db.UpdateDiagnosisPortscanErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.scanID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateDiagnosisAppScanErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	db, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	type args struct {
		errNotifiedAt interface{}
		scanID        uint32
		projectID     uint32
	}
	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			args: args{errNotifiedAt: now, projectID: 1, scanID: 1},
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateDiagnosisAppScanErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, projectID: 1, scanID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateDiagnosisAppScanErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			c.mockClosure(mock)
			err = db.UpdateDiagnosisAppScanErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.scanID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
