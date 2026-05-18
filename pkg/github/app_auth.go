package github

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ca-risken/datasource-api/proto/code"
	jwt "github.com/golang-jwt/jwt/v5"
	ghub "github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

const (
	githubAppJWTBackdate = time.Minute
	githubAppJWTLifetime = 9 * time.Minute
)

var defaultGitHubAppAllowedBaseURLHosts = map[string]struct{}{
	"api.github.com": {},
}

var newGitHubAppHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
}

// AppAuthConfig is the server-side GitHub App credential set.
type AppAuthConfig struct {
	AppID               string
	PrivateKey          string
	AllowedBaseURLHosts []string
}

type appAuth struct {
	appID               int64
	privateKey          *rsa.PrivateKey
	allowedBaseURLHosts map[string]struct{}
}

func newAppAuth(conf *AppAuthConfig) (*appAuth, error) {
	if conf == nil || (conf.AppID == "" && conf.PrivateKey == "") {
		return nil, nil
	}
	if conf.AppID == "" || conf.PrivateKey == "" {
		return nil, errors.New("github app id and private key are required together")
	}
	appID, err := strconv.ParseInt(conf.AppID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse github app id: %w", err)
	}
	privateKey, err := parseGitHubAppPrivateKey(conf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse github app private key: %w", err)
	}
	return &appAuth{
		appID:               appID,
		privateKey:          privateKey,
		allowedBaseURLHosts: allowedGitHubAppBaseURLHosts(conf.AllowedBaseURLHosts),
	}, nil
}

func allowedGitHubAppBaseURLHosts(configuredHosts []string) map[string]struct{} {
	allowedHosts := make(map[string]struct{}, len(defaultGitHubAppAllowedBaseURLHosts)+len(configuredHosts))
	for host := range defaultGitHubAppAllowedBaseURLHosts {
		allowedHosts[host] = struct{}{}
	}
	for _, host := range configuredHosts {
		normalized := strings.ToLower(strings.TrimSpace(host))
		if normalized == "" {
			continue
		}
		allowedHosts[normalized] = struct{}{}
	}
	return allowedHosts
}

func parseGitHubAppPrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	normalized := strings.ReplaceAll(privateKey, `\n`, "\n")
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(normalized))
}

func (g *riskenGitHubClient) SupportsGitHubApp() bool {
	return g.appAuth != nil
}

func (g *riskenGitHubClient) createGitHubAppJWT(now time.Time) (string, error) {
	if g.appAuth == nil {
		return "", errors.New("github app auth is not configured")
	}
	claims := jwt.RegisteredClaims{
		Issuer:    strconv.FormatInt(g.appAuth.appID, 10),
		IssuedAt:  jwt.NewNumericDate(now.Add(-githubAppJWTBackdate)),
		ExpiresAt: jwt.NewNumericDate(now.Add(githubAppJWTLifetime)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(g.appAuth.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign github app jwt: %w", err)
	}
	return signed, nil
}

func (g *riskenGitHubClient) newGitHubAppClient(ctx context.Context, baseURL string) (*ghub.Client, error) {
	base, err := g.appAuth.validateBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	jwtToken, err := g.createGitHubAppJWT(time.Now())
	if err != nil {
		return nil, err
	}
	httpClient := newGitHubAppHTTPClient(ctx, &oauth2.Token{AccessToken: jwtToken, TokenType: "Bearer"})
	client := ghub.NewClient(httpClient)
	if base != nil {
		client.BaseURL = base
	}
	return client, nil
}

func (a *appAuth) validateBaseURL(baseURL string) (*url.URL, error) {
	if baseURL == "" {
		return nil, nil
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf("github app base_url must use https: %s", baseURL)
	}
	if _, ok := a.allowedBaseURLHosts[strings.ToLower(u.Hostname())]; !ok {
		return nil, fmt.Errorf("github app base_url host is not allowed: %s", u.Hostname())
	}
	return u, nil
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
	client, err := g.newGitHubAppClient(ctx, config.BaseUrl)
	if err != nil {
		return "", fmt.Errorf("create github app client: %w", err)
	}
	opts := installationTokenOptions(repoName)
	token, _, err := client.Apps.CreateInstallationToken(ctx, int64(config.InstallationId), opts)
	if err != nil {
		return "", fmt.Errorf("create installation token: %w", err)
	}
	if token.GetToken() == "" {
		return "", errors.New("installation token is empty")
	}
	return token.GetToken(), nil
}

func installationTokenOptions(repoName string) *ghub.InstallationTokenOptions {
	if repoName == "" {
		return nil
	}
	parts := strings.Split(repoName, "/")
	if len(parts) == 2 {
		repoName = parts[1]
	}
	return &ghub.InstallationTokenOptions{
		Repositories: []string{repoName},
	}
}
