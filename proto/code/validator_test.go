package code

import (
	"testing"
	"time"
)

const (
	stringLength65           = "12345678901234567890123456789012345678901234567890123456789012345"
	stringLength129          = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789"
	stringLength256          = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456"
	unixtime19691231T235959  = -1
	unixtime100000101T000000 = 253402268400
)

func TestValidate_ListDataSourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListDataSourceRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListDataSourceRequest{CodeDataSourceId: 1, Name: "name"},
		},
		{
			name:  "OK empty",
			input: &ListDataSourceRequest{},
		},
		{
			name:    "NG length(name)",
			input:   &ListDataSourceRequest{CodeDataSourceId: 1, Name: stringLength65},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_ListGitHubSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListGitHubSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListGitHubSettingRequest{ProjectId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListGitHubSettingRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetGitHubSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetGitHubSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &GetGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetGitHubSettingRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &GetGitHubSettingRequest{GithubSettingId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutGitHubSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutGitHubSettingRequest
		wantErr bool
	}{
		{
			name: "OK",
			input: &PutGitHubSettingRequest{ProjectId: 1, GithubSetting: &GitHubSettingForUpsert{
				ProjectId: 1, Type: Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			}},
		},
		{
			name:    "NG No github_setting",
			input:   &PutGitHubSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "NG Invalid project_id",
			input: &PutGitHubSettingRequest{ProjectId: 999, GithubSetting: &GitHubSettingForUpsert{
				ProjectId: 1, Type: Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteGitHubSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteGitHubSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DeleteGitHubSettingRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteGitHubSettingRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &DeleteGitHubSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutGitleaksSettingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *PutGitleaksSettingRequest
		wantErr bool
	}{
		{
			name: "OK",
			input: &PutGitleaksSettingRequest{ProjectId: 1, GitleaksSetting: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			}},
		},
		{
			name:    "NG No gitleaks_setting",
			input:   &PutGitleaksSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "NG Invalid project_id",
			input: &PutGitleaksSettingRequest{GitleaksSetting: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "repo", Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteGitleaksSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteGitleaksSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DeleteGitleaksSettingRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteGitleaksSettingRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &DeleteGitleaksSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetGitleaksCacheRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetGitleaksCacheRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &GetGitleaksCacheRequest{ProjectId: 1, GithubSettingId: 1, RepositoryFullName: "repo"},
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetGitleaksCacheRequest{GithubSettingId: 1, RepositoryFullName: "repo"},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &GetGitleaksCacheRequest{ProjectId: 1, RepositoryFullName: "repo"},
			wantErr: true,
		},
		{
			name:    "NG Required(repository_full_name)",
			input:   &GetGitleaksCacheRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Length(repository_full_name)",
			input:   &GetGitleaksCacheRequest{GithubSettingId: 1, RepositoryFullName: stringLength256},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutGitleaksCacheRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *PutGitleaksCacheRequest
		wantErr bool
	}{
		{
			name: "OK",
			input: &PutGitleaksCacheRequest{
				ProjectId: 1,
				GitleaksCache: &GitleaksCacheForUpsert{
					GithubSettingId: 1, RepositoryFullName: "repo", ScanAt: now.Unix(),
				},
			},
		},
		{
			name: "NG Required project_id",
			input: &PutGitleaksCacheRequest{
				GitleaksCache: &GitleaksCacheForUpsert{
					GithubSettingId: 1, RepositoryFullName: "repo", ScanAt: now.Unix(),
				},
			},
			wantErr: true,
		},
		{
			name:    "NG Required gitleaks_chache",
			input:   &PutGitleaksCacheRequest{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutDependencySettingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *PutDependencySettingRequest
		wantErr bool
	}{
		{
			name: "OK",
			input: &PutDependencySettingRequest{ProjectId: 1, DependencySetting: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			}},
		},
		{
			name:    "NG No github_setting_id",
			input:   &PutDependencySettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "NG Invalid project_id",
			input: &PutDependencySettingRequest{DependencySetting: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteDependencySettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteDependencySettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DeleteDependencySettingRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteDependencySettingRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &DeleteDependencySettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_ListGitHubEnterpriseOrgRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListGitHubEnterpriseOrgRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListGitHubEnterpriseOrgRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &ListGitHubEnterpriseOrgRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutGitHubEnterpriseOrgRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutGitHubEnterpriseOrgRequest
		wantErr bool
	}{
		{
			name: "OK",
			input: &PutGitHubEnterpriseOrgRequest{ProjectId: 1, GithubEnterpriseOrg: &GitHubEnterpriseOrgForUpsert{
				GithubSettingId: 1, ProjectId: 1, Organization: "org",
			}},
		},
		{
			name:    "NG No github_enterprise_org",
			input:   &PutGitHubEnterpriseOrgRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "NG Invalid project_id",
			input: &PutGitHubEnterpriseOrgRequest{ProjectId: 999, GithubEnterpriseOrg: &GitHubEnterpriseOrgForUpsert{
				GithubSettingId: 1, ProjectId: 1, Organization: "org",
			}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DeleteGitHubEnterpriseOrgRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteGitHubEnterpriseOrgRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DeleteGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1, Organization: "org"},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteGitHubEnterpriseOrgRequest{GithubSettingId: 1, Organization: "org"},
			wantErr: true,
		},
		{
			name:    "NG Required(gitleaks_id)",
			input:   &DeleteGitHubEnterpriseOrgRequest{ProjectId: 1, Organization: "org"},
			wantErr: true,
		},
		{
			name:    "NG Required(organization)",
			input:   &DeleteGitHubEnterpriseOrgRequest{ProjectId: 1, GithubSettingId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_InvokeScanGitleaksRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *InvokeScanGitleaksRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &InvokeScanGitleaksRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &InvokeScanGitleaksRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &InvokeScanGitleaksRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_InvokeScanDependencyRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *InvokeScanDependencyRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &InvokeScanDependencyRequest{ProjectId: 1, GithubSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &InvokeScanDependencyRequest{GithubSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &InvokeScanDependencyRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GitHubSettingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *GitHubSettingForUpsert
		wantErr bool
	}{
		{
			name: "OK",
			input: &GitHubSettingForUpsert{
				Name: "name", ProjectId: 1, Type: Type_ENTERPRISE, BaseUrl: "https://api.github.com/", TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			},
		},
		{
			name: "OK minimize",
			input: &GitHubSettingForUpsert{
				ProjectId: 1, Type: Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			},
		},
		{
			name: "NG Length(name)",
			input: &GitHubSettingForUpsert{
				Name: stringLength65, ProjectId: 1, Type: Type_ENTERPRISE, BaseUrl: "https://api.github.com/", TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			},
			wantErr: true,
		},
		{
			name: "NG Required(project_id)",
			input: &GitHubSettingForUpsert{
				Name: "name", Type: Type_ENTERPRISE, BaseUrl: "https://api.github.com/", TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			},
			wantErr: true,
		},
		{
			name: "NG Length(base_url)",
			input: &GitHubSettingForUpsert{
				Name: "name", ProjectId: 1, Type: Type_ENTERPRISE, BaseUrl: stringLength129, TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			},
			wantErr: true,
		},
		{
			name: "NG Not URL(base_url)",
			input: &GitHubSettingForUpsert{
				ProjectId: 1, Type: Type_ORGANIZATION, BaseUrl: "not URL pattern", TargetResource: "target", GithubUser: "user", PersonalAccessToken: "xxx",
			},
			wantErr: true,
		},
		{
			name: "NG Required(targetResource)",
			input: &GitHubSettingForUpsert{
				Name: "name", ProjectId: 1, Type: Type_ENTERPRISE, TargetResource: "", GithubUser: "user", PersonalAccessToken: "xxx",
			},
			wantErr: true,
		},
		{
			name: "NG Length(targetResource)",
			input: &GitHubSettingForUpsert{
				Name: "name", ProjectId: 1, Type: Type_ENTERPRISE, TargetResource: stringLength129, GithubUser: "user", PersonalAccessToken: "xxx",
			},
			wantErr: true,
		},
		{
			name: "NG Length(github_user)",
			input: &GitHubSettingForUpsert{
				Name: "name", ProjectId: 1, Type: Type_ENTERPRISE, TargetResource: "target", GithubUser: stringLength65,
			},
			wantErr: true,
		},
		{
			name: "NG Length(personal_access_token)",
			input: &GitHubSettingForUpsert{
				Name: "name", ProjectId: 1, Type: Type_ENTERPRISE, TargetResource: "target", GithubUser: "user", PersonalAccessToken: stringLength256,
			},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GitleaksSettingForUpsert(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *GitleaksSettingForUpsert
		wantErr bool
	}{
		{
			name: "OK",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "some-repo", Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
		},
		{
			name: "OK minimize",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1,
			},
		},
		{
			name: "NG Required(github_setting_id)",
			input: &GitleaksSettingForUpsert{
				CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "some-repo", Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Required(code_data_source_id)",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, ProjectId: 1, RepositoryPattern: "some-repo", Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Required(project_id)",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, RepositoryPattern: "some-repo", Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Length(RepositoryPattern)",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: stringLength129, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Uncompilable(RepositoryPattern)",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "*xxx", Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Length(status_detail)",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "some-repo", Status: Status_OK, StatusDetail: stringLength256, ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Min(scan_at)",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "some-repo", Status: Status_OK, StatusDetail: "detail", ScanAt: unixtime19691231T235959,
			},
			wantErr: true,
		},
		{
			name: "NG Max(scan_at)",
			input: &GitleaksSettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, RepositoryPattern: "some-repo", Status: Status_OK, StatusDetail: "detail", ScanAt: unixtime100000101T000000,
			},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GitleaksCacheForUpsert(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *GitleaksCacheForUpsert
		wantErr bool
	}{
		{
			name: "OK",
			input: &GitleaksCacheForUpsert{
				GithubSettingId: 1, RepositoryFullName: "repo", ScanAt: now.Unix(),
			},
		},
		{
			name: "NG Required(github_setting_id)",
			input: &GitleaksCacheForUpsert{
				RepositoryFullName: "repo", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Required(code_data_source_id)",
			input: &GitleaksCacheForUpsert{
				GithubSettingId: 1, ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Required(scan_at)",
			input: &GitleaksCacheForUpsert{
				GithubSettingId: 1, RepositoryFullName: "repo",
			},
			wantErr: true,
		},
		{
			name: "NG Min(scan_at)",
			input: &GitleaksCacheForUpsert{
				GithubSettingId: 1, RepositoryFullName: "repo", ScanAt: unixtime19691231T235959,
			},
			wantErr: true,
		},
		{
			name: "NG Max(scan_at)",
			input: &GitleaksCacheForUpsert{
				GithubSettingId: 1, RepositoryFullName: "repo", ScanAt: unixtime100000101T000000,
			},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_DependencySettingForUpsert(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *DependencySettingForUpsert
		wantErr bool
	}{
		{
			name: "OK",
			input: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
		},
		{
			name: "OK minimize",
			input: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1,
			},
		},
		{
			name: "NG Required(github_setting_id)",
			input: &DependencySettingForUpsert{
				CodeDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Required(code_data_source_id)",
			input: &DependencySettingForUpsert{
				GithubSettingId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Required(project_id)",
			input: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Length(status_detail)",
			input: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: stringLength256, ScanAt: now.Unix(),
			},
			wantErr: true,
		},
		{
			name: "NG Min(scan_at)",
			input: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: unixtime19691231T235959,
			},
			wantErr: true,
		},
		{
			name: "NG Max(scan_at)",
			input: &DependencySettingForUpsert{
				GithubSettingId: 1, CodeDataSourceId: 1, ProjectId: 1, Status: Status_OK, StatusDetail: "detail", ScanAt: unixtime100000101T000000,
			},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GitHubEnterpriseOrgForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *GitHubEnterpriseOrgForUpsert
		wantErr bool
	}{
		{
			name:  "OK",
			input: &GitHubEnterpriseOrgForUpsert{GithubSettingId: 1, Organization: "org", ProjectId: 1},
		},
		{
			name:    "NG Required(github_setting_id)",
			input:   &GitHubEnterpriseOrgForUpsert{Organization: "org", ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(organization)",
			input:   &GitHubEnterpriseOrgForUpsert{GithubSettingId: 1, ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Length(login)",
			input:   &GitHubEnterpriseOrgForUpsert{GithubSettingId: 1, Organization: stringLength129, ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GitHubEnterpriseOrgForUpsert{GithubSettingId: 1, Organization: "login"},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
