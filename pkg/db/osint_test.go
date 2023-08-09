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

func TestListOsintScanErrorForNotify(t *testing.T) {
	cases := []struct {
		name        string
		want        []*OsintScanError
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			want: []*OsintScanError{
				{RelOsintDataSourceID: 1, DataSource: "osint:subdomain", ProjectID: 1, StatusDetail: "error detail"},
				{RelOsintDataSourceID: 2, DataSource: "osint:website", ProjectID: 1, StatusDetail: "error detail"},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListOsintScanErrorForNotify)).WillReturnRows(sqlmock.NewRows([]string{
					"rel_osint_data_source_id", "data_source", "project_id", "status_detail"}).
					AddRow(uint32(1), "osint:subdomain", uint32(1), "error detail").
					AddRow(uint32(2), "osint:website", uint32(1), "error detail"))
			},
		},
		{
			name:    "NG DB error",
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListOsintScanErrorForNotify)).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListOsintScanErrorForNotify(ctx)
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

func TestUpdateOsintErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	type args struct {
		errNotifiedAt        time.Time
		relOsintDataSourceID uint32
		projectID            uint32
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{errNotifiedAt: now, relOsintDataSourceID: 1, projectID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateOsintErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, relOsintDataSourceID: 1, projectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateOsintErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
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
			err = db.UpdateOsintErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.relOsintDataSourceID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
