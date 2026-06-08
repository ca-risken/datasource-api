package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ca-risken/common/pkg/githubappauth"
	"github.com/ca-risken/datasource-api/proto/code"
	ghub "github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

// AppAuthConfig is the server-side GitHub App credential set.
type AppAuthConfig = githubappauth.Config

type appAuth = githubappauth.Client

type OAuthConfig = githubappauth.OAuthConfig

type userOAuth = githubappauth.OAuthClient

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

func newUserOAuth(conf *OAuthConfig) (*userOAuth, error) {
	client, err := githubappauth.NewOAuthClient(conf)
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

func (g *riskenGitHubClient) VerifyUserToServer(ctx context.Context, config *code.GitHubSetting, oauthCode string) (string, error) {
	if g.userOAuth == nil {
		return "", errors.New("github app oauth is not configured")
	}
	if config == nil {
		return "", errors.New("github setting is required")
	}
	if config.AuthMode != code.GitHubAuthModeGitHubApp {
		return "", errors.New("github setting is not github app auth mode")
	}
	token, err := g.userOAuth.ExchangeCode(ctx, oauthCode)
	if err != nil {
		return "", err
	}
	user, err := g.userOAuth.GetAuthenticatedUser(ctx, token)
	if err != nil {
		return "", err
	}
	login := user.GetLogin()
	client, err := g.newGitHubAppClient(ctx, config.BaseUrl)
	if err != nil {
		return login, fmt.Errorf("create github app client: %w", err)
	}
	installation, _, err := findInstallation(ctx, client.Apps, config)
	if err != nil {
		return login, fmt.Errorf("find installation: %w", err)
	}
	if installation.GetID() != int64(config.InstallationId) {
		g.logger.Warnf(ctx, "github app installation_id mismatch: expected=%d, actual=%d", config.InstallationId, installation.GetID())
		return login, errors.New("installation_id does not match target resource")
	}
	if installation.GetRepositorySelection() == "all" {
		if err := g.verifyGitHubUserInstallationAdmin(ctx, token, config, login); err != nil {
			return login, err
		}
		return login, nil
	}

	repositories, err := g.ListRepository(ctx, config, "")
	if err != nil {
		return login, fmt.Errorf("list github app repositories: %w", err)
	}
	if err := g.verifyGitHubUserRepositoryAdmin(ctx, token, repositories, login); err != nil {
		return login, err
	}
	return login, nil
}

func (g *riskenGitHubClient) verifyGitHubUserInstallationAdmin(ctx context.Context, token *oauth2.Token, config *code.GitHubSetting, login string) error {
	switch config.Type {
	case code.Type_ORGANIZATION:
		client, err := g.userOAuth.NewUserClient(ctx, token)
		if err != nil {
			return err
		}
		membership, _, err := client.Organizations.GetOrgMembership(ctx, login, config.TargetResource)
		if err != nil {
			return fmt.Errorf("get organization membership: organization=%s: %w", config.TargetResource, err)
		}
		if membership.GetState() != "active" || membership.GetRole() != "admin" {
			return fmt.Errorf("authenticated github user is not organization admin: organization=%s", config.TargetResource)
		}
		return nil
	case code.Type_USER:
		if !strings.EqualFold(login, config.TargetResource) {
			return fmt.Errorf("authenticated github user does not match target user: target_user=%s", config.TargetResource)
		}
		return nil
	default:
		return fmt.Errorf("unknown github type: type=%s", config.Type.String())
	}
}

func (g *riskenGitHubClient) verifyGitHubUserRepositoryAdmin(ctx context.Context, token *oauth2.Token, repositories []*ghub.Repository, login string) error {
	if len(repositories) == 0 {
		return errors.New("github app repository is required")
	}
	client, err := g.userOAuth.NewUserClient(ctx, token)
	if err != nil {
		return err
	}
	for _, repository := range repositories {
		repositoryFullName := repository.GetFullName()
		owner, repo, err := splitRepositoryFullName(repositoryFullName)
		if err != nil {
			return err
		}
		permission, _, err := client.Repositories.GetPermissionLevel(ctx, owner, repo, login)
		if err != nil {
			return fmt.Errorf("get repository permission: repository_full_name=%s: %w", repositoryFullName, err)
		}
		if permission.GetPermission() != "admin" {
			return fmt.Errorf("authenticated github user is not repository admin: repository_full_name=%s", repositoryFullName)
		}
	}
	return nil
}

func splitRepositoryFullName(repositoryFullName string) (string, string, error) {
	parts := strings.Split(repositoryFullName, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repository name format: %s, expected 'owner/repo'", repositoryFullName)
	}
	return parts[0], parts[1], nil
}
