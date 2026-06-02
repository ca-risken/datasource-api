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

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/code"
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

func useTestGitHubAppTransport(t *testing.T, server *httptest.Server) {
	t.Helper()
	origTransport := http.DefaultTransport
	http.DefaultTransport = server.Client().Transport
	t.Cleanup(func() {
		http.DefaultTransport = origTransport
	})
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
	useTestGitHubAppTransport(t, server)

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
	useTestGitHubAppTransport(t, server)

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
		BaseUrl:        server.URL + "/",
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
	useTestGitHubAppTransport(t, server)

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
		BaseUrl:        server.URL + "/",
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

func TestVerifyUserToServer(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	var gotCode string
	var gotUser bool
	var gotRepos bool
	var gotPermission bool
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/login/oauth/access_token":
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse form: %v", err)
			}
			gotCode = r.Form.Get("code")
			if _, err := w.Write([]byte(`{"access_token":"user-token","token_type":"bearer"}`)); err != nil {
				t.Fatalf("write token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/user":
			gotUser = true
			if _, err := w.Write([]byte(`{"login":"octocat"}`)); err != nil {
				t.Fatalf("write user response: %v", err)
			}
		case r.Method == http.MethodPost && r.URL.Path == "/app/installations/12345/access_tokens":
			if _, err := w.Write([]byte(`{"token":"installation-token"}`)); err != nil {
				t.Fatalf("write installation token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/installation/repositories":
			gotRepos = true
			if _, err := w.Write([]byte(`{"total_count":1,"repositories":[{"full_name":"owner/repo","owner":{"login":"owner"}}]}`)); err != nil {
				t.Fatalf("write repositories response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/repos/owner/repo/collaborators/octocat/permission":
			gotPermission = true
			if _, err := w.Write([]byte(`{"permission":"admin"}`)); err != nil {
				t.Fatalf("write permission response: %v", err)
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
	useTestGitHubAppTransport(t, server)

	client, err := NewGithubClientWithGitHubAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	}, &OAuthConfig{
		ClientID:                 "client-id",
		ClientSecret:             "client-secret",
		OAuthBaseURL:             server.URL,
		APIBaseURL:               server.URL,
		AllowedOAuthBaseURLHosts: []string{serverURL.Hostname()},
		AllowedAPIBaseURLHosts:   []string{serverURL.Hostname()},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())
	got, err := client.VerifyUserToServer(ctx, &code.GitHubSetting{
		AuthMode:       code.GitHubAuthModeGitHubApp,
		Type:           code.Type_ORGANIZATION,
		TargetResource: "owner",
		BaseUrl:        server.URL + "/",
		InstallationId: 12345,
	}, "oauth-code")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "octocat" {
		t.Fatalf("Unexpected verified user: %s", got)
	}
	if gotCode != "oauth-code" || !gotUser || !gotRepos || !gotPermission {
		t.Fatalf("Expected oauth/user/repos/permission calls, gotCode=%s, gotUser=%t, gotRepos=%t, gotPermission=%t", gotCode, gotUser, gotRepos, gotPermission)
	}
}

func TestVerifyUserToServerUserTypeUsesRepositoryPermission(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/login/oauth/access_token":
			if _, err := w.Write([]byte(`{"access_token":"user-token","token_type":"bearer"}`)); err != nil {
				t.Fatalf("write token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/user":
			if _, err := w.Write([]byte(`{"login":"octocat"}`)); err != nil {
				t.Fatalf("write user response: %v", err)
			}
		case r.Method == http.MethodPost && r.URL.Path == "/app/installations/12345/access_tokens":
			if _, err := w.Write([]byte(`{"token":"installation-token"}`)); err != nil {
				t.Fatalf("write installation token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/installation/repositories":
			if _, err := w.Write([]byte(`{"total_count":1,"repositories":[{"full_name":"octocat/repo","owner":{"login":"octocat"}}]}`)); err != nil {
				t.Fatalf("write repositories response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/repos/octocat/repo/collaborators/octocat/permission":
			if _, err := w.Write([]byte(`{"permission":"admin"}`)); err != nil {
				t.Fatalf("write permission response: %v", err)
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
	useTestGitHubAppTransport(t, server)

	client, err := NewGithubClientWithGitHubAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	}, &OAuthConfig{
		ClientID:                 "client-id",
		ClientSecret:             "client-secret",
		OAuthBaseURL:             server.URL,
		APIBaseURL:               server.URL,
		AllowedOAuthBaseURLHosts: []string{serverURL.Hostname()},
		AllowedAPIBaseURLHosts:   []string{serverURL.Hostname()},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())
	got, err := client.VerifyUserToServer(ctx, &code.GitHubSetting{
		AuthMode:       code.GitHubAuthModeGitHubApp,
		Type:           code.Type_USER,
		TargetResource: "octocat",
		BaseUrl:        server.URL + "/",
		InstallationId: 12345,
	}, "oauth-code")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "octocat" {
		t.Fatalf("Unexpected verified user: %s", got)
	}
}

func TestVerifyUserToServerRejectsNonRepositoryAdmin(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/login/oauth/access_token":
			if _, err := w.Write([]byte(`{"access_token":"user-token","token_type":"bearer"}`)); err != nil {
				t.Fatalf("write token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/user":
			if _, err := w.Write([]byte(`{"login":"octocat"}`)); err != nil {
				t.Fatalf("write user response: %v", err)
			}
		case r.Method == http.MethodPost && r.URL.Path == "/app/installations/12345/access_tokens":
			if _, err := w.Write([]byte(`{"token":"installation-token"}`)); err != nil {
				t.Fatalf("write installation token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/installation/repositories":
			if _, err := w.Write([]byte(`{"total_count":1,"repositories":[{"full_name":"owner/repo","owner":{"login":"owner"}}]}`)); err != nil {
				t.Fatalf("write repositories response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/repos/owner/repo/collaborators/octocat/permission":
			if _, err := w.Write([]byte(`{"permission":"write"}`)); err != nil {
				t.Fatalf("write permission response: %v", err)
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
	useTestGitHubAppTransport(t, server)

	client, err := NewGithubClientWithGitHubAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	}, &OAuthConfig{
		ClientID:                 "client-id",
		ClientSecret:             "client-secret",
		OAuthBaseURL:             server.URL,
		APIBaseURL:               server.URL,
		AllowedOAuthBaseURLHosts: []string{serverURL.Hostname()},
		AllowedAPIBaseURLHosts:   []string{serverURL.Hostname()},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())
	_, err = client.VerifyUserToServer(ctx, &code.GitHubSetting{
		AuthMode:       code.GitHubAuthModeGitHubApp,
		Type:           code.Type_ORGANIZATION,
		TargetResource: "owner",
		BaseUrl:        server.URL + "/",
		InstallationId: 12345,
	}, "oauth-code")
	if err == nil {
		t.Fatal("Expected error but got none")
	}
	if err.Error() != "authenticated github user is not repository admin: repository_full_name=owner/repo" {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestVerifyUserToServerRejectsMissingRepository(t *testing.T) {
	_, privateKeyPEM := generateRSAPrivateKeyPEM(t)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/login/oauth/access_token":
			if _, err := w.Write([]byte(`{"access_token":"user-token","token_type":"bearer"}`)); err != nil {
				t.Fatalf("write token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/user":
			if _, err := w.Write([]byte(`{"login":"octocat"}`)); err != nil {
				t.Fatalf("write user response: %v", err)
			}
		case r.Method == http.MethodPost && r.URL.Path == "/app/installations/12345/access_tokens":
			if _, err := w.Write([]byte(`{"token":"installation-token"}`)); err != nil {
				t.Fatalf("write installation token response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/installation/repositories":
			if _, err := w.Write([]byte(`{"total_count":0,"repositories":[]}`)); err != nil {
				t.Fatalf("write repositories response: %v", err)
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
	useTestGitHubAppTransport(t, server)

	client, err := NewGithubClientWithGitHubAppAuth("default-token", &AppAuthConfig{
		AppID:               "12345",
		PrivateKey:          privateKeyPEM,
		AllowedBaseURLHosts: []string{serverURL.Hostname()},
	}, &OAuthConfig{
		ClientID:                 "client-id",
		ClientSecret:             "client-secret",
		OAuthBaseURL:             server.URL,
		APIBaseURL:               server.URL,
		AllowedOAuthBaseURLHosts: []string{serverURL.Hostname()},
		AllowedAPIBaseURLHosts:   []string{serverURL.Hostname()},
	}, logging.NewLogger())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, server.Client())
	_, err = client.VerifyUserToServer(ctx, &code.GitHubSetting{
		AuthMode:       code.GitHubAuthModeGitHubApp,
		Type:           code.Type_ORGANIZATION,
		TargetResource: "owner",
		BaseUrl:        server.URL + "/",
		InstallationId: 12345,
	}, "oauth-code")
	if err == nil {
		t.Fatal("Expected error but got none")
	}
	if err.Error() != "github app repository is required" {
		t.Fatalf("Unexpected error: %v", err)
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
	if _, err := client.appAuth.ValidateBaseURL("https://ghe.example.com/api/v3/"); err != nil {
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
	got, err := client.appAuth.ValidateBaseURL("https://ghe.example.com/api/v3")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got.String() != "https://ghe.example.com/api/v3/" {
		t.Fatalf("Unexpected base URL: %s", got.String())
	}
}
