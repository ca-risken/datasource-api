package code

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// validateRepositoryName validates repository name format (owner/repo)
func validateRepositoryName(value any) error {
	s, ok := value.(string)
	if !ok {
		return nil // Skip validation for non-string values
	}
	if s == "" {
		return nil // Empty is allowed (optional field)
	}
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return errors.New("repository name must be in format 'owner/repo'")
	}
	if parts[0] == "" || parts[1] == "" {
		return errors.New("repository name must have both owner and repo name")
	}
	return nil
}

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

// Validate ListGitleaksCacheRequest
func (g *ListGitleaksCacheRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.GithubSettingId, validation.Required),
	)
}

// Validate GetGitleaksCacheRequest
func (g *GetGitleaksCacheRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.GithubSettingId, validation.Required),
		validation.Field(&g.RepositoryFullName, validation.Required, validation.Length(0, 255)),
	)
}

// Validate PutGitleaksCacheRequest
func (p *PutGitleaksCacheRequest) Validate() error {
	if p.GitleaksCache == nil {
		return errors.New("required gitleaks_cache")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required),
	); err != nil {
		return err
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

// Validate PutCodeScanSettingRequest
func (p *PutCodeScanSettingRequest) Validate() error {
	if p.CodeScanSetting == nil {
		return errors.New("required CodeScanSetting")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required, validation.In(p.CodeScanSetting.ProjectId)),
	); err != nil {
		return err
	}
	return p.CodeScanSetting.Validate()
}

// Validate DeleteCodeScanSettingRequest
func (d *DeleteCodeScanSettingRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.GithubSettingId, validation.Required),
	)
}

// Validate InvokeScanRequest
func (i *InvokeScanGitleaksRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.ProjectId, validation.Required),
		validation.Field(&i.GithubSettingId, validation.Required),
		validation.Field(&i.RepositoryName, validation.Length(0, 255), validation.By(validateRepositoryName)),
	)
}

// Validate InvokeScanRequest
func (i *InvokeScanDependencyRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.ProjectId, validation.Required),
		validation.Field(&i.GithubSettingId, validation.Required),
		validation.Field(&i.RepositoryName, validation.Length(0, 255), validation.By(validateRepositoryName)),
	)
}

// Validate InvokeScanCodeScanRequest
func (i *InvokeScanCodeScanRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.ProjectId, validation.Required),
		validation.Field(&i.GithubSettingId, validation.Required),
		validation.Field(&i.RepositoryName, validation.Length(0, 255), validation.By(validateRepositoryName)),
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
		validation.Field(&g.RepositoryPattern, validation.Length(0, 128)),
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

// Validate CodeScanSettingForUpsert
func (c *CodeScanSettingForUpsert) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.GithubSettingId, validation.Required),
		validation.Field(&c.CodeDataSourceId, validation.Required),
		validation.Field(&c.ProjectId, validation.Required),
		validation.Field(&c.RepositoryPattern, validation.Length(0, 128)),
		validation.Field(&c.StatusDetail, validation.Length(0, 255)),
		validation.Field(&c.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate CodeScanRepositoryStatusForUpsert
func (c *CodeScanRepositoryStatusForUpsert) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.GithubSettingId, validation.Required),
		validation.Field(&c.RepositoryFullName, validation.Required, validation.Length(0, 255)),
		validation.Field(&c.Status, validation.Required),
		validation.Field(&c.StatusDetail, validation.Length(0, 255)),
		validation.Field(&c.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate PutCodeScanRepositoryStatusRequest
func (p *PutCodeScanRepositoryStatusRequest) Validate() error {
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required),
		validation.Field(&p.GithubSettingId, validation.Required),
		validation.Field(&p.RepositoryFullName, validation.Required),
		validation.Field(&p.Status, validation.Required),
	); err != nil {
		return err
	}
	return nil
}

// Validate GitleaksRepositoryStatusForUpsert
func (g *GitleaksRepositoryStatusForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.GithubSettingId, validation.Required),
		validation.Field(&g.RepositoryFullName, validation.Required, validation.Length(0, 255)),
		validation.Field(&g.Status, validation.Required),
		validation.Field(&g.StatusDetail, validation.Length(0, 255)),
		validation.Field(&g.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate DependencyRepositoryStatusForUpsert
func (d *DependencyRepositoryStatusForUpsert) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.GithubSettingId, validation.Required),
		validation.Field(&d.RepositoryFullName, validation.Required, validation.Length(0, 255)),
		validation.Field(&d.Status, validation.Required),
		validation.Field(&d.StatusDetail, validation.Length(0, 255)),
		validation.Field(&d.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}
