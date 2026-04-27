package github

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

const RETRY_NUM uint64 = 3

type GithubServiceClient interface {
	ListRepository(ctx context.Context, config *code.GitHubSetting, repoName string) ([]*github.Repository, error)
	Clone(ctx context.Context, token string, cloneURL string, dstDir string) error
}

type GitHubRepoService interface {
	List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	ListByOrg(ctx context.Context, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error)
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
}

type GitHubAppService interface {
	FindOrganizationInstallation(ctx context.Context, org string) (*github.Installation, *github.Response, error)
	FindUserInstallation(ctx context.Context, user string) (*github.Installation, *github.Response, error)
	CreateInstallationToken(ctx context.Context, id int64, opts *github.InstallationTokenOptions) (*github.InstallationToken, *github.Response, error)
	ListRepos(ctx context.Context, opts *github.ListOptions) (*github.ListRepositories, *github.Response, error)
}

type GitHubV3Client struct {
	Repositories GitHubRepoService
	Apps         GitHubAppService
	*github.Client
}

type AppAuthConfig struct {
	AppID      string
	PrivateKey string
}

type gitHubAppAuthenticator struct {
	appID      int64
	privateKey *rsa.PrivateKey
}

type riskenGitHubClient struct {
	defaultToken string
	appAuth      *gitHubAppAuthenticator
	retryer      backoff.BackOff
	logger       logging.Logger
}

func NewGithubClient(defaultToken string, logger logging.Logger) *riskenGitHubClient {
	return NewGithubClientWithAppAuth(defaultToken, nil, logger)
}

func NewGithubClientWithAppAuth(defaultToken string, appAuthCfg *AppAuthConfig, logger logging.Logger) *riskenGitHubClient {
	appAuth, err := newGitHubAppAuthenticator(appAuthCfg)
	if err != nil {
		logger.Warnf(context.Background(), "GitHub App auth disabled: err=%+v", err)
	}
	return &riskenGitHubClient{
		defaultToken: defaultToken,
		appAuth:      appAuth,
		retryer:      backoff.WithMaxRetries(backoff.NewExponentialBackOff(), RETRY_NUM),
		logger:       logger,
	}
}

func newGitHubAppAuthenticator(cfg *AppAuthConfig) (*gitHubAppAuthenticator, error) {
	if cfg == nil {
		return nil, nil
	}
	appID := strings.TrimSpace(cfg.AppID)
	if appID == "" {
		return nil, nil
	}
	privateKey, err := parseGitHubAppPrivateKey(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}
	if privateKey == nil {
		return nil, errors.New("github app private key is required")
	}
	numericAppID, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid github app id: %w", err)
	}
	return &gitHubAppAuthenticator{
		appID:      numericAppID,
		privateKey: privateKey,
	}, nil
}

func parseGitHubAppPrivateKey(rawKey string) (*rsa.PrivateKey, error) {
	keyText := strings.TrimSpace(rawKey)
	if keyText == "" {
		return nil, nil
	}
	keyText = strings.ReplaceAll(keyText, "\\n", "\n")
	block, _ := pem.Decode([]byte(keyText))
	if block == nil {
		return nil, errors.New("decode github app private key PEM")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse github app private key: %w", err)
	}
	key, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("github app private key must be RSA")
	}
	return key, nil
}

func (g *riskenGitHubClient) newV3Client(ctx context.Context, token, baseURL string) (*GitHubV3Client, error) {
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getToken(token, g.defaultToken)},
	))
	client := github.NewClient(httpClient)
	if baseURL != "" { // Default: "https://api.github.com/"
		u, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		client.BaseURL = u
	}
	return &GitHubV3Client{
		Repositories: client.Repositories,
		Apps:         client.Apps,
		Client:       client,
	}, nil
}

func (g *riskenGitHubClient) newAppClient(ctx context.Context, baseURL string) (*GitHubV3Client, error) {
	if g.appAuth == nil {
		return nil, errors.New("github app auth is not configured")
	}
	jwtToken, err := g.appAuth.createJWT(time.Now())
	if err != nil {
		return nil, err
	}
	return g.newV3Client(ctx, jwtToken, baseURL)
}

func getToken(token, defaultToken string) string {
	if token != "" {
		return token
	}
	return defaultToken
}

func (g *riskenGitHubClient) Clone(ctx context.Context, token string, cloneURL string, dstDir string) error {
	operation := func() error {
		_, err := git.PlainClone(dstDir, false, &git.CloneOptions{
			URL: cloneURL,
			Auth: &http.BasicAuth{
				Username: "dummy", // anything except an empty string
				Password: getToken(token, g.defaultToken),
			},
		})
		return err
	}

	if err := backoff.RetryNotify(operation, g.retryer, g.newRetryLogger(ctx, "github clone")); err != nil {
		return fmt.Errorf("failed to clone %s to %s: %w", cloneURL, dstDir, err)
	}

	return nil
}

func (g *riskenGitHubClient) ListRepository(ctx context.Context, config *code.GitHubSetting, repoName string) ([]*github.Repository, error) {
	token, err := g.resolveAccessToken(ctx, config, repoName)
	if err != nil {
		return nil, err
	}
	client, err := g.newV3Client(ctx, token, config.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("create github-v3 client: %w", err)
	}

	// First check if repoName is specified for single repository scan
	if repoName != "" {
		repository, err := g.GetSingleRepository(ctx, client, config, repoName)
		if err != nil {
			return nil, err
		}
		return []*github.Repository{repository}, nil
	}

	if config.PersonalAccessToken == "" && g.appAuth != nil {
		return g.listRepositoryForInstallation(ctx, client.Apps, config.TargetResource)
	}

	// Handle bulk repository scan based on config.Type
	var repos []*github.Repository

	switch config.Type {
	case code.Type_ORGANIZATION:
		repos, err = g.listRepositoryForOrg(ctx, client.Repositories, config)
		if err != nil {
			return repos, err
		}

	case code.Type_USER:
		// Check target user(targetResource) == authenticated user(PAT user)
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			return nil, err
		}
		isAuthUser := user.Login != nil && *user.Login == config.TargetResource
		repos, err = g.listRepositoryForUser(ctx, client.Repositories, config, isAuthUser)
		if err != nil {
			return repos, err
		}

	default:
		return nil, fmt.Errorf("unknown github type: type=%s", config.Type.String())
	}

	return repos, nil
}

func (g *riskenGitHubClient) resolveAccessToken(ctx context.Context, config *code.GitHubSetting, repoName string) (string, error) {
	if config.PersonalAccessToken != "" {
		return config.PersonalAccessToken, nil
	}
	if g.appAuth != nil {
		return g.resolveInstallationToken(ctx, config, repoName)
	}
	return g.defaultToken, nil
}

func (g *riskenGitHubClient) resolveInstallationToken(ctx context.Context, config *code.GitHubSetting, repoName string) (string, error) {
	appClient, err := g.newAppClient(ctx, config.BaseUrl)
	if err != nil {
		return "", fmt.Errorf("create github app client: %w", err)
	}
	if repoName != "" {
		if err := validateRepositoryBelongsToTarget(config, repoName); err != nil {
			return "", err
		}
	}
	installation, err := g.findInstallation(ctx, appClient.Apps, config)
	if err != nil {
		return "", err
	}
	token, _, err := appClient.Apps.CreateInstallationToken(ctx, installation.GetID(), installationTokenOptions(repoName))
	if err != nil {
		return "", fmt.Errorf("create installation token: %w", err)
	}
	if token == nil || token.Token == nil || strings.TrimSpace(token.GetToken()) == "" {
		return "", errors.New("installation token is empty")
	}
	return token.GetToken(), nil
}

func installationTokenOptions(repoName string) *github.InstallationTokenOptions {
	if repoName == "" {
		return nil
	}
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return nil
	}
	return &github.InstallationTokenOptions{
		Repositories: []string{parts[1]},
	}
}

func (g *riskenGitHubClient) findInstallation(ctx context.Context, appSvc GitHubAppService, config *code.GitHubSetting) (*github.Installation, error) {
	switch config.Type {
	case code.Type_ORGANIZATION:
		installation, _, err := appSvc.FindOrganizationInstallation(ctx, config.TargetResource)
		if err != nil {
			return nil, fmt.Errorf("find organization installation for %s: %w", config.TargetResource, err)
		}
		return installation, nil
	case code.Type_USER:
		installation, _, err := appSvc.FindUserInstallation(ctx, config.TargetResource)
		if err != nil {
			return nil, fmt.Errorf("find user installation for %s: %w", config.TargetResource, err)
		}
		return installation, nil
	default:
		return nil, fmt.Errorf("github app auth does not support type: type=%s", config.Type.String())
	}
}

func (g *riskenGitHubClient) listRepositoryForInstallation(ctx context.Context, appSvc GitHubAppService, targetResource string) ([]*github.Repository, error) {
	var allRepo []*github.Repository
	opt := &github.ListOptions{PerPage: 100}
	for {
		repositories, resp, err := appSvc.ListRepos(ctx, opt)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, errors.New("github app list repositories response is nil")
		}
		repos := []*github.Repository{}
		if repositories != nil {
			repos = repositories.Repositories
		}
		if targetResource != "" {
			repos = filterRepositoriesByOwner(repos, targetResource)
		}
		allRepo = append(allRepo, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepo, nil
}

func filterRepositoriesByOwner(repos []*github.Repository, owner string) []*github.Repository {
	if owner == "" {
		return repos
	}
	filtered := make([]*github.Repository, 0, len(repos))
	for _, repo := range repos {
		if repo == nil || repo.Owner == nil || repo.Owner.Login == nil {
			continue
		}
		if strings.EqualFold(repo.GetOwner().GetLogin(), owner) {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

const (
	githubVisibilityAll string = "all"
)

func (g *riskenGitHubClient) listRepositoryForUser(ctx context.Context, repository GitHubRepoService, config *code.GitHubSetting, isAuthUser bool) ([]*github.Repository, error) {
	repos, err := g.listRepositoryForUserWithOption(ctx, repository, config.TargetResource, isAuthUser)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (g *riskenGitHubClient) listRepositoryForUserWithOption(ctx context.Context, repository GitHubRepoService, login string, isAuthUser bool) ([]*github.Repository, error) {
	var allRepo []*github.Repository
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Type:        githubVisibilityAll,
	}

	for {
		var repos []*github.Repository
		var resp *github.Response
		var err error

		if isAuthUser {
			// Use authenticated user endpoint to access private repositories
			repos, resp, err = repository.List(ctx, "", opt)
		} else {
			// Use public user endpoint for other users
			repos, resp, err = repository.List(ctx, login, opt)
		}

		if err != nil {
			return nil, err
		}
		g.logger.Infof(ctx, "Success GitHub API for user repos, %s,login:%s, option:%+v, repo_count: %d, response:%+v", login, opt, len(repos), resp)

		for _, r := range repos {
			// Filter repositories by user owner
			if r.Owner != nil && r.Owner.Login != nil && *r.Owner.Login == login {
				allRepo = append(allRepo, r)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepo, nil
}

func (g *riskenGitHubClient) listRepositoryForOrg(ctx context.Context, repository GitHubRepoService, config *code.GitHubSetting) ([]*github.Repository, error) {
	repos, err := g.listRepositoryForOrgWithOption(ctx, repository, config.TargetResource)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (g *riskenGitHubClient) listRepositoryForOrgWithOption(ctx context.Context, repository GitHubRepoService, login string) ([]*github.Repository, error) {
	var allRepo []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Type:        githubVisibilityAll,
	}
	for {
		repos, resp, err := repository.ListByOrg(ctx, login, opt)
		if err != nil {
			return nil, err
		}
		g.logger.Infof(ctx, "Success GitHub API for user repos, login:%s, option:%+v, repo_count: %d, response:%+v", login, opt, len(repos), resp)
		allRepo = append(allRepo, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepo, nil
}

func (g *riskenGitHubClient) GetSingleRepository(ctx context.Context, client *GitHubV3Client, config *code.GitHubSetting, repoName string) (*github.Repository, error) {
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository name format: %s, expected 'owner/repo'", repoName)
	}
	if err := validateRepositoryBelongsToTarget(config, repoName); err != nil {
		return nil, err
	}
	owner, repo := parts[0], parts[1]

	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository %s: %w", repoName, err)
	}

	return repository, nil
}

func validateRepositoryBelongsToTarget(config *code.GitHubSetting, repoName string) error {
	parts := strings.Split(repoName, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository name format: %s, expected 'owner/repo'", repoName)
	}
	owner := parts[0]
	if !strings.EqualFold(owner, config.TargetResource) {
		return fmt.Errorf("repository %s does not belong to %s %s", repoName, config.Type.String(), config.TargetResource)
	}
	return nil
}

func (t *riskenGitHubClient) newRetryLogger(ctx context.Context, funcName string) func(error, time.Duration) {
	return func(err error, ti time.Duration) {
		t.logger.Warnf(ctx, "[RetryLogger] %s error: duration=%+v, err=%+v", funcName, ti, err)
	}
}

func (a *gitHubAppAuthenticator) createJWT(now time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Add(-1 * time.Minute).Unix(),
		"exp": now.Add(9 * time.Minute).Unix(),
		"iss": strconv.FormatInt(a.appID, 10),
	})
	signedToken, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign github app jwt: %w", err)
	}
	return signedToken, nil
}
