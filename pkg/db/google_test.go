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

func TestListGCPScanErrorForNotify(t *testing.T) {
	cases := []struct {
		name        string
		want        []*GCPScanError
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			want: []*GCPScanError{
				{GCPID: 1, GoogleDataSourceID: 1, DataSource: "google:portscan", ProjectID: 1, StatusDetail: "error detail"},
				{GCPID: 1, GoogleDataSourceID: 2, DataSource: "google:scc", ProjectID: 1, StatusDetail: "error detail"},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListGCPScanError)).WillReturnRows(sqlmock.NewRows([]string{
					"gcp_id", "google_data_source_id", "data_source", "project_id", "status_detail"}).
					AddRow(uint32(1), uint32(1), "google:portscan", uint32(1), "error detail").
					AddRow(uint32(1), uint32(2), "google:scc", uint32(1), "error detail"))
			},
		},
		{
			name:    "NG DB error",
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListGCPScanError)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatal(err)
			}
			c.mockClosure(mock)
			got, err := db.ListGCPScanErrorForNotify(ctx)
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

func TestUpdateGCPErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	type args struct {
		errNotifiedAt      time.Time
		gcpID              uint32
		googleDataSourceID uint32
		projectID          uint32
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{errNotifiedAt: now, gcpID: 1, googleDataSourceID: 1, projectID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateGCPErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, gcpID: 1, googleDataSourceID: 1, projectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateGCPErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock, err := newDBMock()
			if err != nil {
				t.Fatal(err)
			}
			c.mockClosure(mock)
			err = db.UpdateGCPErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.gcpID, c.args.googleDataSourceID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
