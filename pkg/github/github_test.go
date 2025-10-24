package github

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/code"
	"github.com/google/go-github/v44/github"
)

type fakeGitHubRepoService struct {
	repos []*github.Repository
	resp  *github.Response
	err   error
}

func makeGitHubRepository(name, login string) github.Repository {
	return github.Repository{
		Name: &name,
		Owner: &github.User{
			Login: &login,
		},
	}
}

func PointerString(input string) *string {
	return &input
}

func newfakeGitHubRepoService(empty bool, name, login string, err error) *fakeGitHubRepoService {
	if empty {
		return &fakeGitHubRepoService{
			resp: &github.Response{
				NextPage: 0,
			},
		}
	}
	repo := makeGitHubRepository(name, login)
	return &fakeGitHubRepoService{
		err: err,
		repos: []*github.Repository{
			&repo,
		},
		resp: &github.Response{
			NextPage: 0,
		},
	}
}

func (f *fakeGitHubRepoService) List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return f.repos, f.resp, f.err
}
func (f *fakeGitHubRepoService) ListByOrg(ctx context.Context, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	return f.repos, f.resp, f.err
}

func (f *fakeGitHubRepoService) Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	if len(f.repos) > 0 {
		return f.repos[0], f.resp, f.err
	}
	return nil, f.resp, f.err
}

func Test_listRepositoryForUserWithOption(t *testing.T) {
	cases := []struct {
		name       string
		repository GitHubRepoService
		login      string
		isAuthUser bool
		want       []*github.Repository
		wantError  bool
	}{
		{
			name:       "OK",
			login:      "owner",
			isAuthUser: true,
			repository: newfakeGitHubRepoService(false, "repo", "owner", nil),
			want: []*github.Repository{
				{
					Name:  PointerString("repo"),
					Owner: &github.User{Login: PointerString("owner")},
				},
			},
		},
		{
			name:       "OK empty",
			login:      "owner",
			isAuthUser: true,
			repository: newfakeGitHubRepoService(true, "", "", nil),
			want:       []*github.Repository{},
		},
		{
			name:       "NG List Error",
			login:      "owner",
			isAuthUser: true,
			repository: newfakeGitHubRepoService(false, "", "", errors.New("something error")),
			want:       []*github.Repository{},
			wantError:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			githubClient := NewGithubClient("token", logging.NewLogger())

			got, err := githubClient.listRepositoryForUserWithOption(ctx, c.repository, c.login, c.isAuthUser)
			if c.wantError && err == nil {
				t.Fatal("Unexpected no error")
			}
			if !c.wantError && err != nil {
				t.Fatalf("Unexpected error occurred, err=%+v", err)
			}
			if len(got) != len(c.want) {
				t.Fatalf("Unexpected not matching: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func Test_listRepositoryForOrgWithOption(t *testing.T) {
	cases := []struct {
		name       string
		repository GitHubRepoService
		login      string
		want       []*github.Repository
		wantError  bool
	}{
		{
			name:       "OK",
			login:      "owner",
			repository: newfakeGitHubRepoService(false, "repo", "owner", nil),
			want: []*github.Repository{
				{
					Name:  PointerString("repo"),
					Owner: &github.User{Login: PointerString("owner")},
				},
			},
		},
		{
			name:       "OK empty",
			repository: newfakeGitHubRepoService(true, "", "", nil),
			want:       []*github.Repository{},
		},
		{
			name:       "NG List Error",
			login:      "owner",
			repository: newfakeGitHubRepoService(false, "", "", errors.New("something error")),
			want:       []*github.Repository{},
			wantError:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			githubClient := NewGithubClient("token", logging.NewLogger())
			got, err := githubClient.listRepositoryForOrgWithOption(ctx, c.repository, c.login)
			if c.wantError && err == nil {
				t.Fatal("Unexpected no error")
			}
			if !c.wantError && err != nil {
				t.Fatalf("Unexpected error occured, err=%+v", err)
			}
			if len(got) != len(c.want) {
				t.Fatalf("Unexpected not matching: want=%+v, got=%+v", c.want, got)
			}

		})
	}
}

func TestGetSingleRepository(t *testing.T) {
	cases := []struct {
		name          string
		repository    GitHubRepoService
		config        *code.GitHubSetting
		repoName      string
		want          *github.Repository
		wantError     bool
		expectedError string
	}{
		{
			name:     "OK - valid repository",
			repoName: "owner/repo",
			config: &code.GitHubSetting{
				TargetResource: "owner",
				Type:           code.Type_USER,
			},
			repository: newfakeGitHubRepoService(false, "repo", "owner", nil),
			want: &github.Repository{
				Name:  PointerString("repo"),
				Owner: &github.User{Login: PointerString("owner")},
			},
		},
		{
			name:     "NG - invalid repository name format",
			repoName: "invalid-format",
			config: &code.GitHubSetting{
				TargetResource: "owner",
				Type:           code.Type_USER,
			},
			repository:    newfakeGitHubRepoService(false, "repo", "owner", nil),
			wantError:     true,
			expectedError: "invalid repository name format",
		},
		{
			name:     "NG - repository does not belong to target resource",
			repoName: "other-owner/repo",
			config: &code.GitHubSetting{
				TargetResource: "owner",
				Type:           code.Type_USER,
			},
			repository:    newfakeGitHubRepoService(false, "repo", "other-owner", nil),
			wantError:     true,
			expectedError: "repository other-owner/repo does not belong to USER owner",
		},
		{
			name:     "NG - GitHub API error",
			repoName: "owner/repo",
			config: &code.GitHubSetting{
				TargetResource: "owner",
				Type:           code.Type_USER,
			},
			repository:    newfakeGitHubRepoService(false, "repo", "owner", errors.New("GitHub API error")),
			wantError:     true,
			expectedError: "failed to get repository owner/repo",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			githubClient := NewGithubClient("token", logging.NewLogger())

			client := &GitHubV3Client{
				Repositories: c.repository,
			}

			got, err := githubClient.GetSingleRepository(ctx, client, c.config, c.repoName)

			if c.wantError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				if c.expectedError != "" && !strings.Contains(err.Error(), c.expectedError) {
					t.Fatalf("Expected error to contain '%s', but got: %v", c.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if got == nil {
					t.Fatal("Expected repository but got nil")
				}
				if got.Name == nil || *got.Name != *c.want.Name {
					t.Fatalf("Expected repository name %s, got %v", *c.want.Name, got.Name)
				}
				if got.Owner == nil || got.Owner.Login == nil || *got.Owner.Login != *c.want.Owner.Login {
					t.Fatalf("Expected owner %s, got %v", *c.want.Owner.Login, got.Owner.Login)
				}
			}
		})
	}
}
