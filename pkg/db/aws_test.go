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

func TestListScanError(t *testing.T) {
	cases := []struct {
		name        string
		want        []*AWSScanError
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK",
			want: []*AWSScanError{
				{AWSID: 1, AWSDataSourceID: 1, DataSource: "aws:portscan", ProjectID: 1, StatusDetail: "error detail"},
				{AWSID: 1, AWSDataSourceID: 2, DataSource: "aws:access-analyzer", ProjectID: 1, StatusDetail: "error detail"},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListAWSScanError)).WillReturnRows(sqlmock.NewRows([]string{
					"aws_id", "aws_data_source_id", "data_source", "project_id", "status_detail"}).
					AddRow(uint32(1), uint32(1), "aws:portscan", uint32(1), "error detail").
					AddRow(uint32(1), uint32(2), "aws:access-analyzer", uint32(1), "error detail"))
			},
		},
		{
			name:    "NG DB error",
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectListAWSScanError)).WillReturnError(errors.New("DB error"))
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
			got, err := db.ListAWSScanErrorForNotify(ctx)
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

func TestUpdateErrorNotifiedAt(t *testing.T) {
	now := time.Now()
	type args struct {
		errNotifiedAt   time.Time
		awsID           uint32
		awsDataSourceID uint32
		projectID       uint32
	}

	cases := []struct {
		name        string
		args        args
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "OK",
			args:    args{errNotifiedAt: now, awsID: 1, awsDataSourceID: 1, projectID: 1},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateAWSErrorNotifiedAt)).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "NG DB error",
			args:    args{errNotifiedAt: now, awsID: 1, awsDataSourceID: 1, projectID: 1},
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(updateAWSErrorNotifiedAt)).WillReturnError(errors.New("DB error"))
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
			err = db.UpdateAWSErrorNotifiedAt(ctx, c.args.errNotifiedAt, c.args.awsID, c.args.awsDataSourceID, c.args.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
