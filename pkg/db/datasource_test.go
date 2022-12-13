package db

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCleanWithNoProject(t *testing.T) {
	client, mock, err := newDBMock()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		mockSQL []string
		wantErr bool
		mockErr error
	}{
		{
			name: "OK",
			mockSQL: []string{
				"delete tbl from aws tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from aws_rel_data_source tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from gcp tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from gcp_data_source tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from wpscan_setting tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from portscan_setting tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from portscan_target tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from application_scan tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from application_scan_basic_setting tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from osint tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from rel_osint_data_source tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from osint_detect_word tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from code_github_setting tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from code_gitleaks_setting tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from code_dependency_setting tbl where not exists(select * from project p where p.project_id = tbl.project_id)",
				"delete tbl from code_gitleaks_cache tbl where not exists(select * from code_github_setting github where github.code_github_setting_id = tbl.code_github_setting_id)",
			},
			wantErr: false,
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			for _, sql := range c.mockSQL {
				mock.ExpectExec(regexp.QuoteMeta(sql)).WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
			}
			if c.mockErr != nil {
				mock.ExpectExec(regexp.QuoteMeta(`delete tbl from`)).WillReturnError(c.mockErr)
			}

			err := client.CleanWithNoProject(ctx)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
		})
	}
}
