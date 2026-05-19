package github

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/code"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

func generateRSAPrivateKeyPEM(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa private key: %v", err)
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	return privateKey, string(pem.EncodeToMemory(block))
}

func TestNewGithubClientWithAppAuth(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	cases := []struct {
		name      string
		conf      *AppAuthConfig
		wantApp   bool
		wantError bool
	}{
		{
			name: "OK no app auth",
		},
		{
			name: "OK empty app auth",
			conf: &AppAuthConfig{},
		},
		{
			name:    "OK app auth",
			conf:    &AppAuthConfig{AppID: "12345", PrivateKey: privateKeyPEM},
			wantApp: true,
		},
		{
			name:      "NG missing private key",
			conf:      &AppAuthConfig{AppID: "12345"},
			wantError: true,
		},
		{
			name:      "NG invalid app id",
			conf:      &AppAuthConfig{AppID: "invalid", PrivateKey: privateKeyPEM},
			wantError: true,
		},
		{
			name:      "NG invalid private key",
			conf:      &AppAuthConfig{AppID: "12345", PrivateKey: "invalid"},
			wantError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, err := NewGithubClientWithAppAuth("default-token", c.conf, logging.NewLogger())
			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got := client.SupportsGitHubApp(); got != c.wantApp {
				t.Fatalf("Unexpected GitHub App support: want=%t, got=%t", c.wantApp, got)
			}
		})
	}
}

func TestParseGitHubAppPrivateKey(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	escapedPEM := strings.ReplaceAll(privateKeyPEM, "\n", `\n`)
	cases := []struct {
		name      string
		input     string
		wantError bool
	}{
		{name: "OK raw pem", input: privateKeyPEM},
		{name: "OK escaped pem", input: escapedPEM},
		{name: "NG invalid pem", input: "invalid", wantError: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseGitHubAppPrivateKey(c.input)
			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("Expected private key but got nil")
			}
		})
	}
}

func TestCreateGitHubAppJWT(t *testing.T) {
	privateKey, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	client, err := NewGithubClientWithAppAuth("default-token", &AppAuthConfig{AppID: "12345", PrivateKey: privateKeyPEM}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	now := time.Now().Truncate(time.Second)
	tokenString, err := client.createGitHubAppJWT(now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			t.Fatalf("Unexpected signing method: %v", token.Method.Alg())
		}
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("Expected valid token")
	}
	if claims.Issuer != "12345" {
		t.Fatalf("Unexpected issuer: %s", claims.Issuer)
	}
	if claims.IssuedAt.Unix() != now.Add(-githubAppJWTBackdate).Unix() {
		t.Fatalf("Unexpected issued_at: %v", claims.IssuedAt)
	}
	if claims.ExpiresAt.Unix() != now.Add(githubAppJWTLifetime).Unix() {
		t.Fatalf("Unexpected expires_at: %v", claims.ExpiresAt)
	}
}

func TestResolveInstallationToken(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	var gotAuthorization string
	var gotRepositories []string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/app/installations/12345/access_tokens" {
			t.Fatalf("Unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("Unexpected method: %s", r.Method)
		}
		gotAuthorization = r.Header.Get("Authorization")
		var body struct {
			Repositories []string `json:"repositories"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		gotRepositories = body.Repositories
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"token":"installation-token"}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}
	origNewGitHubAppHTTPClient := newGitHubAppHTTPClient
	newGitHubAppHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
		client := server.Client()
		client.Transport = &oauth2.Transport{
			Source: oauth2.StaticTokenSource(token),
			Base:   client.Transport,
		}
		return client
	}
	defer func() {
		newGitHubAppHTTPClient = origNewGitHubAppHTTPClient
	}()

	client, err := NewGithubClientWithAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	got, err := client.ResolveInstallationToken(context.Background(), &code.GitHubSetting{
		BaseUrl:        server.URL + "/",
		InstallationId: 12345,
	}, "owner/repo")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "installation-token" {
		t.Fatalf("Unexpected token: %s", got)
	}
	if !strings.HasPrefix(gotAuthorization, "Bearer ") {
		t.Fatalf("Unexpected authorization header: %s", gotAuthorization)
	}
	if len(gotRepositories) != 1 || gotRepositories[0] != "repo" {
		t.Fatalf("Unexpected repositories: %+v", gotRepositories)
	}
}

func TestResolveInstallationTokenError(t *testing.T) {
	client := NewGithubClient("default-token", logging.NewLogger())
	if _, err := client.ResolveInstallationToken(context.Background(), &code.GitHubSetting{InstallationId: 12345}, ""); err == nil {
		t.Fatal("Expected error but got none")
	}
}

func TestResolveInstallationTokenRejectsUntrustedBaseURL(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	client, err := NewGithubClientWithAppAuth("default-token", &AppAuthConfig{AppID: "12345", PrivateKey: privateKeyPEM}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	_, err = client.ResolveInstallationToken(context.Background(), &code.GitHubSetting{
		BaseUrl:        "https://attacker.example/",
		InstallationId: 12345,
	}, "")
	if err == nil {
		t.Fatal("Expected error but got none")
	}
	if !strings.Contains(err.Error(), "base_url host is not allowed") {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestVerifyInstallation(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	var gotFindInstallation bool
	var gotCreateToken bool
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/orgs/target/installation":
			gotFindInstallation = true
			if _, err := w.Write([]byte(`{"id":12345}`)); err != nil {
				t.Fatalf("write installation response: %v", err)
			}
		case r.Method == http.MethodPost && r.URL.Path == "/app/installations/12345/access_tokens":
			gotCreateToken = true
			if _, err := w.Write([]byte(`{"token":"installation-token"}`)); err != nil {
				t.Fatalf("write token response: %v", err)
			}
		default:
			t.Fatalf("Unexpected request: method=%s path=%s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}
	origNewGitHubAppHTTPClient := newGitHubAppHTTPClient
	newGitHubAppHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
		client := server.Client()
		client.Transport = &oauth2.Transport{
			Source: oauth2.StaticTokenSource(token),
			Base:   client.Transport,
		}
		return client
	}
	defer func() {
		newGitHubAppHTTPClient = origNewGitHubAppHTTPClient
	}()

	client, err := NewGithubClientWithAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = client.VerifyInstallation(context.Background(), &code.GitHubSetting{
		Type:           code.Type_ORGANIZATION,
		TargetResource: "target",
		BaseUrl:        server.URL,
		InstallationId: 12345,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !gotFindInstallation || !gotCreateToken {
		t.Fatalf("Expected find installation and create token calls, gotFindInstallation=%t, gotCreateToken=%t", gotFindInstallation, gotCreateToken)
	}
}

func TestVerifyInstallationReturnsAbstractMismatchError(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/orgs/target/installation":
			if _, err := w.Write([]byte(`{"id":99999}`)); err != nil {
				t.Fatalf("write installation response: %v", err)
			}
		default:
			t.Fatalf("Unexpected request: method=%s path=%s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}
	origNewGitHubAppHTTPClient := newGitHubAppHTTPClient
	newGitHubAppHTTPClient = func(ctx context.Context, token *oauth2.Token) *http.Client {
		client := server.Client()
		client.Transport = &oauth2.Transport{
			Source: oauth2.StaticTokenSource(token),
			Base:   client.Transport,
		}
		return client
	}
	defer func() {
		newGitHubAppHTTPClient = origNewGitHubAppHTTPClient
	}()

	client, err := NewGithubClientWithAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = client.VerifyInstallation(context.Background(), &code.GitHubSetting{
		Type:           code.Type_ORGANIZATION,
		TargetResource: "target",
		BaseUrl:        server.URL,
		InstallationId: 12345,
	})
	if err == nil {
		t.Fatal("Expected error but got none")
	}
	if err.Error() != "installation_id does not match target resource" {
		t.Fatalf("Unexpected error: %v", err)
	}
	if strings.Contains(err.Error(), "99999") || strings.Contains(err.Error(), "12345") {
		t.Fatalf("Error exposes installation_id: %v", err)
	}
}

func TestResolveInstallationTokenAllowsConfiguredBaseURLHost(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	client, err := NewGithubClientWithAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{"ghe.example.com"},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if _, err := client.appAuth.validateBaseURL("https://ghe.example.com/api/v3/"); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestValidateGitHubAppBaseURLAddsTrailingSlash(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	client, err := NewGithubClientWithAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{"ghe.example.com"},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	got, err := client.appAuth.validateBaseURL("https://ghe.example.com/api/v3")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got.String() != "https://ghe.example.com/api/v3/" {
		t.Fatalf("Unexpected base URL: %s", got.String())
	}
}

func TestInstallationTokenOptions(t *testing.T) {
	cases := []struct {
		name      string
		repoName  string
		wantRepo  string
		wantNil   bool
		wantError bool
	}{
		{name: "OK empty", wantNil: true},
		{name: "OK repo", repoName: "repo", wantRepo: "repo"},
		{name: "OK owner repo", repoName: "owner/repo", wantRepo: "repo"},
		{name: "OK trims spaces", repoName: " owner/repo ", wantRepo: "repo"},
		{name: "NG too many segments", repoName: "a/b/c", wantError: true},
		{name: "NG empty owner", repoName: "/repo", wantError: true},
		{name: "NG empty repo", repoName: "owner/", wantError: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := installationTokenOptions(c.repoName)
			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if c.wantNil {
				if got != nil {
					t.Fatalf("Expected nil options but got %+v", got)
				}
				return
			}
			if got == nil || len(got.Repositories) != 1 || got.Repositories[0] != c.wantRepo {
				t.Fatalf("Unexpected options: %+v", got)
			}
		})
	}
}
