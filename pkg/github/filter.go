package github

import (
	"strings"

	"github.com/google/go-github/v44/github"
)

const (
	githubVisibilityPublic   string = "public"
	githubVisibilityInternal string = "internal"
	githubVisibilityPrivate  string = "private"
)

// FilterOptions represents options for filtering repositories
type FilterOptions struct {
	// RepositoryPattern filters repositories by name pattern (substring match)
	// If empty, no filtering is applied
	RepositoryPattern string
	// ScanPublic indicates whether to include public repositories
	ScanPublic bool
	// ScanInternal indicates whether to include internal repositories
	ScanInternal bool
	// ScanPrivate indicates whether to include private repositories
	ScanPrivate bool
}

// FilterByNamePattern filters repositories by name pattern
// If pattern is empty, returns all repositories
func FilterByNamePattern(repos []*github.Repository, pattern string) []*github.Repository {
	if pattern == "" {
		return repos
	}

	filteredRepos := make([]*github.Repository, 0)
	for _, repo := range repos {
		if repo.Name != nil && strings.Contains(*repo.Name, pattern) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos
}

// FilterByVisibility filters repositories by visibility
func FilterByVisibility(repos []*github.Repository, scanPublic, scanInternal, scanPrivate bool) []*github.Repository {
	filteredRepos := make([]*github.Repository, 0)
	for _, repo := range repos {
		if repo.Visibility == nil {
			continue
		}
		visibility := *repo.Visibility
		if scanPublic && visibility == githubVisibilityPublic {
			filteredRepos = append(filteredRepos, repo)
		}
		if scanInternal && visibility == githubVisibilityInternal {
			filteredRepos = append(filteredRepos, repo)
		}
		if scanPrivate && visibility == githubVisibilityPrivate {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos
}

// ApplyFilters applies all filters specified in FilterOptions to the repository list
func ApplyFilters(repos []*github.Repository, opts *FilterOptions) []*github.Repository {
	if opts == nil {
		return repos
	}

	// Apply visibility filter first
	repos = FilterByVisibility(repos, opts.ScanPublic, opts.ScanInternal, opts.ScanPrivate)

	// Then apply name pattern filter
	repos = FilterByNamePattern(repos, opts.RepositoryPattern)

	return repos
}
