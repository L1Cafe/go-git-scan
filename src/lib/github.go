package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v72/github"
	"golang.org/x/oauth2"
)

func FetchGitHubRepos(ctx context.Context, token string, username string) ([]RepoInfo, error) {
	var client *github.Client
	if token == "" {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	opt := &github.RepositoryListByUserOptions{
		Type:        "owner",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allGHRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByUser(ctx, username, opt)
		if err != nil {
			return nil, fmt.Errorf("listing GitHub repositories: %w", err)
		}
		allGHRepos = append(allGHRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var results []RepoInfo
	for _, r := range allGHRepos {
		var lastActivity time.Time
		if r.PushedAt != nil {
			lastActivity = r.PushedAt.Time
		} else if r.UpdatedAt != nil {
			lastActivity = r.UpdatedAt.Time
		} else if r.CreatedAt != nil {
			lastActivity = r.CreatedAt.Time
		}
		results = append(results, RepoInfo{
			Name:         r.GetName(),
			CloneURL:     r.GetCloneURL(),
			IsFork:       r.GetFork(),
			IsDisabled:   r.GetDisabled(),
			LastActivity: lastActivity,
		})
	}
	return results, nil
}
