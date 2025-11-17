package github

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v44/github"
)

func TestFilterByNamePattern(t *testing.T) {
	type args struct {
		repos   []*github.Repository
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []*github.Repository
	}{
		{
			name: "Return repositories contained risken in repository name",
			args: args{
				repos: []*github.Repository{
					{
						Name: github.String("risken-core"),
					},
					{
						Name: github.String("core"),
					},
				},
				pattern: "risken",
			},
			want: []*github.Repository{
				{
					Name: github.String("risken-core"),
				},
			},
		},
		{
			name: "Return all repositories when pattern is empty",
			args: args{
				repos: []*github.Repository{
					{
						Name: github.String("risken-core"),
					},
					{
						Name: github.String("core"),
					},
				},
				pattern: "",
			},
			want: []*github.Repository{
				{
					Name: github.String("risken-core"),
				},
				{
					Name: github.String("core"),
				},
			},
		},
		{
			name: "Return empty when no match",
			args: args{
				repos: []*github.Repository{
					{
						Name: github.String("risken-core"),
					},
					{
						Name: github.String("core"),
					},
				},
				pattern: "nonexistent",
			},
			want: []*github.Repository{},
		},
		{
			name: "Handle nil repository name",
			args: args{
				repos: []*github.Repository{
					{
						Name: nil,
					},
					{
						Name: github.String("core"),
					},
				},
				pattern: "core",
			},
			want: []*github.Repository{
				{
					Name: github.String("core"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByNamePattern(tt.args.repos, tt.args.pattern)
			if len(got) != len(tt.want) {
				t.Errorf("FilterByNamePattern() length = %v, want %v", len(got), len(tt.want))
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterByNamePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterByVisibility(t *testing.T) {
	visibilityPublic := "public"
	visibilityInternal := "internal"
	visibilityPrivate := "private"
	type args struct {
		repos        []*github.Repository
		scanPublic   bool
		scanInternal bool
		scanPrivate  bool
	}
	tests := []struct {
		name string
		args args
		want []*github.Repository
	}{
		{
			name: "Return public repositories",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("public-repo"),
						Visibility: &visibilityPublic,
					},
					{
						Name:       github.String("internal-repo"),
						Visibility: &visibilityInternal,
					},
					{
						Name:       github.String("private-repo"),
						Visibility: &visibilityPrivate,
					},
				},
				scanPublic: true,
			},
			want: []*github.Repository{
				{
					Name:       github.String("public-repo"),
					Visibility: &visibilityPublic,
				},
			},
		},
		{
			name: "Return internal repositories",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("public-repo"),
						Visibility: &visibilityPublic,
					},
					{
						Name:       github.String("internal-repo"),
						Visibility: &visibilityInternal,
					},
					{
						Name:       github.String("private-repo"),
						Visibility: &visibilityPrivate,
					},
				},
				scanInternal: true,
			},
			want: []*github.Repository{
				{
					Name:       github.String("internal-repo"),
					Visibility: &visibilityInternal,
				},
			},
		},
		{
			name: "Return private repositories",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("public-repo"),
						Visibility: &visibilityPublic,
					},
					{
						Name:       github.String("internal-repo"),
						Visibility: &visibilityInternal,
					},
					{
						Name:       github.String("private-repo"),
						Visibility: &visibilityPrivate,
					},
				},
				scanPrivate: true,
			},
			want: []*github.Repository{
				{
					Name:       github.String("private-repo"),
					Visibility: &visibilityPrivate,
				},
			},
		},
		{
			name: "Return all repositories",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("public-repo"),
						Visibility: &visibilityPublic,
					},
					{
						Name:       github.String("internal-repo"),
						Visibility: &visibilityInternal,
					},
					{
						Name:       github.String("private-repo"),
						Visibility: &visibilityPrivate,
					},
				},
				scanPublic:   true,
				scanInternal: true,
				scanPrivate:  true,
			},
			want: []*github.Repository{
				{
					Name:       github.String("public-repo"),
					Visibility: &visibilityPublic,
				},
				{
					Name:       github.String("internal-repo"),
					Visibility: &visibilityInternal,
				},
				{
					Name:       github.String("private-repo"),
					Visibility: &visibilityPrivate,
				},
			},
		},
		{
			name: "Return empty when no visibility matches",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("public-repo"),
						Visibility: &visibilityPublic,
					},
				},
				scanPublic: false,
			},
			want: []*github.Repository{},
		},
		{
			name: "Handle nil visibility",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("repo-without-visibility"),
						Visibility: nil,
					},
					{
						Name:       github.String("public-repo"),
						Visibility: &visibilityPublic,
					},
				},
				scanPublic: true,
			},
			want: []*github.Repository{
				{
					Name:       github.String("public-repo"),
					Visibility: &visibilityPublic,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByVisibility(tt.args.repos, tt.args.scanPublic, tt.args.scanInternal, tt.args.scanPrivate)
			if len(got) != len(tt.want) {
				t.Errorf("FilterByVisibility() length = %v, want %v", len(got), len(tt.want))
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterByVisibility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyFilters(t *testing.T) {
	visibilityPublic := "public"
	visibilityPrivate := "private"
	const testLimitRepositorySizeKb = 100000
	type args struct {
		repos []*github.Repository
		opts  *FilterOptions
	}
	tests := []struct {
		name string
		args args
		want []*github.Repository
	}{
		{
			name: "Apply both visibility and name pattern filters",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("risken-public"),
						Visibility: &visibilityPublic,
					},
					{
						Name:       github.String("risken-private"),
						Visibility: &visibilityPrivate,
					},
					{
						Name:       github.String("other-public"),
						Visibility: &visibilityPublic,
					},
				},
				opts: &FilterOptions{
					RepositoryPattern: "risken",
					ScanPublic:        true,
					ScanPrivate:       false,
				},
			},
			want: []*github.Repository{
				{
					Name:       github.String("risken-public"),
					Visibility: &visibilityPublic,
				},
			},
		},
		{
			name: "Return all when opts is nil",
			args: args{
				repos: []*github.Repository{
					{
						Name: github.String("repo1"),
					},
					{
						Name: github.String("repo2"),
					},
				},
				opts: nil,
			},
			want: []*github.Repository{
				{
					Name: github.String("repo1"),
				},
				{
					Name: github.String("repo2"),
				},
			},
		},
		{
			name: "Apply only visibility filter",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("public-repo"),
						Visibility: &visibilityPublic,
					},
					{
						Name:       github.String("private-repo"),
						Visibility: &visibilityPrivate,
					},
				},
				opts: &FilterOptions{
					ScanPublic:  true,
					ScanPrivate: false,
				},
			},
			want: []*github.Repository{
				{
					Name:       github.String("public-repo"),
					Visibility: &visibilityPublic,
				},
			},
		},
		{
			name: "Apply only name pattern filter",
			args: args{
				repos: []*github.Repository{
					{
						Name:       github.String("risken-core"),
						Visibility: &visibilityPublic,
					},
					{
						Name:       github.String("other-repo"),
						Visibility: &visibilityPublic,
					},
				},
				opts: &FilterOptions{
					RepositoryPattern: "risken",
					ScanPublic:        true,
					ScanPrivate:       true,
				},
			},
			want: []*github.Repository{
				{
					Name:       github.String("risken-core"),
					Visibility: &visibilityPublic,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplyFilters(tt.args.repos, tt.args.opts, testLimitRepositorySizeKb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
