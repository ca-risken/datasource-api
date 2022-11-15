package code

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Validate ListDataSourceRequest
func (l *ListDataSourceRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Name, validation.Length(0, 64)),
	)
}

// Validate ListGitHubSettingRequest
func (l *ListGitHubSettingRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
	)
}

// Validate GetGitHubSettingRequest
func (l *GetGitHubSettingRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.GithubSettingId, validation.Required),
	)
}

// Validate PutGitHubSettingRequest
func (p *PutGitHubSettingRequest) Validate() error {
	if p.GithubSetting == nil {
		return errors.New("required GitHub")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required, validation.In(p.GithubSetting.ProjectId)),
	); err != nil {
		return err
	}
	return p.GithubSetting.Validate()
}

// Validate DeleteGitHubSettingRequest
func (d *DeleteGitHubSettingRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.GithubSettingId, validation.Required),
	)
}

// Validate PutGitleaksSettingRequest
func (p *PutGitleaksSettingRequest) Validate() error {
	if p.GitleaksSetting == nil {
		return errors.New("required Gitleaks")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required, validation.In(p.GitleaksSetting.ProjectId)),
	); err != nil {
		return err
	}
	return p.GitleaksSetting.Validate()
}

// Validate DeleteGitleaksSettingRequest
func (d *DeleteGitleaksSettingRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.GithubSettingId, validation.Required),
	)
}

// Validate GetGitleaksCacheRequest
func (g *GetGitleaksCacheRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GithubSettingId, validation.Required),
		validation.Field(&g.RepositoryFullName, validation.Required, validation.Length(0, 255)),
	)
}

// Validate PutGitleaksCacheRequest
func (p *PutGitleaksCacheRequest) Validate() error {
	if p.GitleaksCache == nil {
		return errors.New("required gitleaks_cache")
	}
	return p.GitleaksCache.Validate()
}

// Validate PutDependencySettingRequest
func (p *PutDependencySettingRequest) Validate() error {
	if p.DependencySetting == nil {
		return errors.New("required DependencySetting")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required, validation.In(p.DependencySetting.ProjectId)),
	); err != nil {
		return err
	}
	return p.DependencySetting.Validate()
}

// Validate DeleteDependencySettingRequest
func (d *DeleteDependencySettingRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.GithubSettingId, validation.Required),
	)
}

// Validate ListEnterpriseOrgRequest
func (l *ListGitHubEnterpriseOrgRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.GithubSettingId, validation.Required),
		validation.Field(&l.ProjectId, validation.Required),
	)
}

// Validate PutEnterpriseOrgRequest
func (p *PutGitHubEnterpriseOrgRequest) Validate() error {
	if p.GithubEnterpriseOrg == nil {
		return errors.New("required EnterpriseOrg")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required, validation.In(p.GithubEnterpriseOrg.ProjectId)),
	); err != nil {
		return err
	}
	return p.GithubEnterpriseOrg.Validate()
}

// Validate DeleteEnterpriseOrgRequest
func (d *DeleteGitHubEnterpriseOrgRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.GithubSettingId, validation.Required),
		validation.Field(&d.Organization, validation.Required),
	)
}

// Validate InvokeScanRequest
func (i *InvokeScanGitleaksRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.ProjectId, validation.Required),
		validation.Field(&i.GithubSettingId, validation.Required),
	)
}

// Validate InvokeScanRequest
func (i *InvokeScanDependencyRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.ProjectId, validation.Required),
		validation.Field(&i.GithubSettingId, validation.Required),
	)
}

/**
 * Entity
**/

// Validate GitHubSettingForUpsert
func (g *GitHubSettingForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.Name, validation.Length(0, 64)),
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.BaseUrl, validation.Length(0, 128), is.URL),
		validation.Field(&g.TargetResource, validation.Required, validation.Length(0, 128)),
		validation.Field(&g.GithubUser, validation.Length(0, 64)),
		validation.Field(&g.PersonalAccessToken, validation.Length(0, 255)),
	)
}

// Validate GitleaksSettingForUpsert
func (g *GitleaksSettingForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GithubSettingId, validation.Required),
		validation.Field(&g.CodeDataSourceId, validation.Required),
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.RepositoryPattern, validation.Length(0, 128), validation.By(compilableRegexp(g.RepositoryPattern))),
		validation.Field(&g.StatusDetail, validation.Length(0, 255)),
		validation.Field(&g.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GitleaksCacheForUpsert
func (g *GitleaksCacheForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GithubSettingId, validation.Required),
		validation.Field(&g.RepositoryFullName, validation.Required, validation.Length(0, 255)),
		validation.Field(&g.ScanAt, validation.Required, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate DependencySettingForUpsert
func (g *DependencySettingForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GithubSettingId, validation.Required),
		validation.Field(&g.CodeDataSourceId, validation.Required),
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.StatusDetail, validation.Length(0, 255)),
		validation.Field(&g.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate EnterpriseOrgForUpsert
func (g *GitHubEnterpriseOrgForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GithubSettingId, validation.Required),
		validation.Field(&g.Organization, validation.Required, validation.Length(0, 128)),
		validation.Field(&g.ProjectId, validation.Required),
	)
}

// Check the `ptn`(string) that is compilable regexp pattern
func compilableRegexp(ptn string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if s != ptn {
			return fmt.Errorf("Unexpected string, got: %+v", ptn)
		}
		if _, err := regexp.Compile(ptn); err != nil {
			return fmt.Errorf("Could not regexp complie, pattern=%s, err=%+v", ptn, err)
		}
		return nil
	}
}
