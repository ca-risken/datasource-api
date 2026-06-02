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
	if err := g.verifyGitHubUserPermission(ctx, token, config, login); err != nil {
		return login, err
	}
	return login, nil
}

func (g *riskenGitHubClient) verifyGitHubUserPermission(ctx context.Context, token *oauth2.Token, config *code.GitHubSetting, login string) error {
	switch config.Type {
	case code.Type_USER:
		if !strings.EqualFold(login, config.TargetResource) {
			return errors.New("authenticated github user does not match target user")
		}
		return nil
	case code.Type_ORGANIZATION:
		return g.verifyOrganizationOwner(ctx, token, config.TargetResource, login)
	default:
		return fmt.Errorf("unknown github type: type=%s", config.Type.String())
	}
}

func (g *riskenGitHubClient) verifyOrganizationOwner(ctx context.Context, token *oauth2.Token, org, login string) error {
	client, err := g.userOAuth.NewUserClient(ctx, token)
	if err != nil {
		return err
	}
	membership, _, err := client.Organizations.GetOrgMembership(ctx, login, org)
	if err != nil {
		return fmt.Errorf("get organization membership: %w", err)
	}
	if membership == nil || membership.GetRole() != "admin" {
		return errors.New("authenticated github user is not organization owner")
	}
	return nil
}
