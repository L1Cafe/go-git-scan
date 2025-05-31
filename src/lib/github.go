// go-git-scan - public commit analyser
// Copyright (C) 2025 Leon Castillejos

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

	client.UserAgent = "go-git-scan/0.1"

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
