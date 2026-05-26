package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/common/pkg/githubappauth"
	"github.com/ca-risken/datasource-api/proto/code"
	ghub "github.com/google/go-github/v44/github"
)

// AppAuthConfig is the server-side GitHub App credential set.
type AppAuthConfig = githubappauth.Config

type appAuth = githubappauth.Client

func newAppAuth(conf *AppAuthConfig) (*appAuth, error) {
	client, err := githubappauth.NewClient(conf)
	if err != nil {
		return nil, err
	}
	if !client.Enabled() {
		return nil, nil
	}
	return client, nil
}

func (g *riskenGitHubClient) SupportsGitHubApp() bool {
	return g.appAuth != nil && g.appAuth.Enabled()
}

func (g *riskenGitHubClient) newGitHubAppClient(ctx context.Context, baseURL string) (*ghub.Client, error) {
	if g.appAuth == nil {
		return nil, errors.New("github app auth is not configured")
	}
	return g.appAuth.NewGitHubClient(ctx, baseURL)
}

func (g *riskenGitHubClient) ResolveInstallationToken(ctx context.Context, config *code.GitHubSetting, repoName string) (string, error) {
	if g.appAuth == nil {
		return "", errors.New("github app auth is not configured")
	}
	if config == nil {
		return "", errors.New("github setting is required")
	}
	if config.InstallationId == 0 {
		return "", errors.New("installation_id is required")
	}
	return g.appAuth.ResolveInstallationToken(ctx, &githubappauth.InstallationTokenConfig{
		BaseURL:        config.BaseUrl,
		InstallationID: config.InstallationId,
	}, repoName)
}

func (g *riskenGitHubClient) VerifyInstallation(ctx context.Context, config *code.GitHubSetting) error {
	if g.appAuth == nil {
		return errors.New("github app auth is not configured")
	}
	if config == nil {
		return errors.New("github setting is required")
	}
	if config.InstallationId == 0 {
		return errors.New("installation_id is required")
	}
	client, err := g.newGitHubAppClient(ctx, config.BaseUrl)
	if err != nil {
		return fmt.Errorf("create github app client: %w", err)
	}
	installation, _, err := findInstallation(ctx, client.Apps, config)
	if err != nil {
		return fmt.Errorf("find installation: %w", err)
	}
	if installation.GetID() != int64(config.InstallationId) {
		g.logger.Warnf(ctx, "github app installation_id mismatch: expected=%d, actual=%d", config.InstallationId, installation.GetID())
		return errors.New("installation_id does not match target resource")
	}
	if _, err := g.ResolveInstallationToken(ctx, config, ""); err != nil {
		return fmt.Errorf("resolve installation token: %w", err)
	}
	return nil
}

func findInstallation(ctx context.Context, appSvc GitHubAppService, config *code.GitHubSetting) (*ghub.Installation, *ghub.Response, error) {
	switch config.Type {
	case code.Type_ORGANIZATION:
		return appSvc.FindOrganizationInstallation(ctx, config.TargetResource)
	case code.Type_USER:
		return appSvc.FindUserInstallation(ctx, config.TargetResource)
	default:
		return nil, nil, fmt.Errorf("unknown github type: type=%s", config.Type.String())
	}
}
