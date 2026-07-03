package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

func (g *riskenGitHubClient) VerifyInstallation(ctx context.Context, config *code.GitHubSetting) (uint64, error) {
	if g.appAuth == nil {
		return 0, errors.New("github app auth is not configured")
	}
	if config == nil {
		return 0, errors.New("github setting is required")
	}
	client, err := g.newGitHubAppClient(ctx, config.BaseUrl)
	if err != nil {
		return 0, fmt.Errorf("create github app client: %w", err)
	}
	installation, _, err := findInstallation(ctx, client.Apps, config)
	if err != nil {
		return 0, fmt.Errorf("find installation: %w", err)
	}
	resolvedInstallationID := installation.GetID()
	if resolvedInstallationID <= 0 {
		return 0, errors.New("installation_id is required")
	}
	if _, err := g.appAuth.ResolveInstallationToken(ctx, &githubappauth.InstallationTokenConfig{
		BaseURL:        config.BaseUrl,
		InstallationID: uint64(resolvedInstallationID),
	}, ""); err != nil {
		return 0, fmt.Errorf("resolve installation token: %w", err)
	}
	return uint64(resolvedInstallationID), nil
}

func (g *riskenGitHubClient) GetGitHubAppInstallationStatus(ctx context.Context, config *code.GitHubSetting) (*code.GitHubAppInstallationStatus, error) {
	if g.appAuth == nil {
		return nil, errors.New("github app auth is not configured")
	}
	if config == nil {
		return nil, errors.New("github setting is required")
	}
	client, err := g.newGitHubAppClient(ctx, config.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("create github app client: %w", err)
	}
	installation, _, err := findInstallation(ctx, client.Apps, config)
	if err != nil {
		return nil, fmt.Errorf("find installation: %w", err)
	}
	installationID := installation.GetID()
	if installationID <= 0 {
		return nil, errors.New("installation_id is required")
	}
	config.InstallationId = uint64(installationID)
	repositories, err := g.ListRepository(ctx, config, "")
	if err != nil {
		return nil, fmt.Errorf("list github app repositories: %w", err)
	}
	return &code.GitHubAppInstallationStatus{
		TargetResource:      config.TargetResource,
		Installed:           true,
		InstallationId:      uint64(installationID),
		RepositorySelection: installation.GetRepositorySelection(),
		RepositoryCount:     uint32(len(repositories)),
		Reason:              "",
	}, nil
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
		return "", errors.New("exchange github app oauth code failed")
	}
	user, err := g.userOAuth.GetAuthenticatedUser(ctx, token)
	if err != nil {
		if err.Error() == "authenticated github user login is empty" {
			return "", errors.New("authenticated github user login is empty")
		}
		return "", errors.New("get authenticated github user failed")
	}
	login := user.GetLogin()
	if login == "" {
		return "", errors.New("authenticated github user login is empty")
	}
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

	repositories, err := buildGitHubRepositoriesFromSetting(config.GetGithubAppSettingRepository())
	if err != nil {
		return login, err
	}

	isInstallationAdmin, err := g.hasGitHubUserInstallationAdmin(ctx, token, config, login)
	if err != nil {
		g.logger.Warnf(ctx, "Failed to verify github app installation admin; fallback to repository admin check: target_resource=%s, github_user=%s, err=%+v", config.TargetResource, login, err)
	}
	if isInstallationAdmin {
		return login, nil
	}

	if err := g.verifyGitHubUserRepositoryAdmin(ctx, token, repositories, login); err != nil {
		return login, err
	}
	return login, nil
}

func buildGitHubRepositoriesFromSetting(repositories []*code.GitHubAppSettingRepository) ([]*ghub.Repository, error) {
	if len(repositories) == 0 {
		return nil, errors.New("github app repository is required")
	}
	ghubRepositories := make([]*ghub.Repository, 0, len(repositories))
	for _, repo := range repositories {
		if repo == nil {
			return nil, errors.New("github app repository is required")
		}
		fullName := repo.GetGithubRepositoryFullName()
		if fullName == "" {
			return nil, errors.New("github repository full name is required")
		}
		ghubRepositories = append(ghubRepositories, &ghub.Repository{FullName: ghub.String(fullName)})
	}
	return ghubRepositories, nil
}

func (g *riskenGitHubClient) hasGitHubUserInstallationAdmin(ctx context.Context, token *oauth2.Token, config *code.GitHubSetting, login string) (bool, error) {
	switch config.Type {
	case code.Type_ORGANIZATION:
		client, err := g.userOAuth.NewUserClient(ctx, token)
		if err != nil {
			return false, err
		}
		membership, _, err := client.Organizations.GetOrgMembership(ctx, login, config.TargetResource)
		if err != nil {
			var ghErr *ghub.ErrorResponse
			if errors.As(err, &ghErr) && ghErr.Response != nil && ghErr.Response.StatusCode == http.StatusNotFound {
				return false, nil
			}
			return false, fmt.Errorf("get organization membership: organization=%s: %w", config.TargetResource, err)
		}
		return membership.GetState() == "active" && membership.GetRole() == "admin", nil
	case code.Type_USER:
		return strings.EqualFold(login, config.TargetResource), nil
	default:
		return false, fmt.Errorf("unknown github type: type=%s", config.Type.String())
	}
}

func (g *riskenGitHubClient) verifyGitHubUserInstallationAdmin(ctx context.Context, token *oauth2.Token, config *code.GitHubSetting, login string) error {
	isInstallationAdmin, err := g.hasGitHubUserInstallationAdmin(ctx, token, config, login)
	if err != nil {
		return err
	}
	if isInstallationAdmin {
		return nil
	}
	switch config.Type {
	case code.Type_ORGANIZATION:
		return fmt.Errorf("authenticated github user is not organization admin: organization=%s", config.TargetResource)
	case code.Type_USER:
		return fmt.Errorf("authenticated github user does not match target user: target_user=%s", config.TargetResource)
	default:
		return fmt.Errorf("unknown github type: type=%s", config.Type.String())
	}
}

func (g *riskenGitHubClient) verifyGitHubUserRepositoryAdmin(ctx context.Context, token *oauth2.Token, repositories []*ghub.Repository, login string) error {
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
