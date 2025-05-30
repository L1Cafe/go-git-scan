package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/L1Cafe/go-git-scan/lib"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {

	appName := filepath.Base(os.Args[0])

	if len(os.Args) < 3 {
		log.Fatal("Usage: " + appName + " {github,gitlab,codeberg,sourcehut} username_or_userid\n" +
			"Alternatively: " + appName + " url \"https://example.com/user/username\"\n")
	}
	platform := strings.ToLower(os.Args[1])
	identifier := os.Args[2]

	ctx := context.Background()
	var reposToProcess []lib.RepoInfo

	switch platform {
	case "github":
		ghToken := ""
		if os.Getenv("GITHUB_TOKEN") != "" {
			ghToken = os.Getenv("GITHUB_TOKEN")
		}
		var err error
		reposToProcess, err = lib.FetchGitHubRepos(ctx, ghToken, identifier)
		if err != nil {
			log.Fatalf("Error fetching GitHub repositories: %v", err)
		}
	case "gitlab":
		log.Fatalf("GitLab support is not yet implemented in this version.")
	case "codeberg":
		log.Fatalf("Codeberg support is not yet implemented in this version.")
	case "sourcehut":
		log.Fatalf("SourceHut support is not yet implemented in this version.")
	case "url":
		log.Fatalf("URL support is not yet implemented.")
	default:
		log.Fatalf("Unsupported platform: %s.", platform)
	}

	if len(reposToProcess) == 0 {
		fmt.Printf("No repositories/projects found for identifier %s on %s.\n", identifier, platform)
		return
	}
	fmt.Printf("Found %d repositories/projects to process.\n", len(reposToProcess))

	tempDir, err := os.MkdirTemp("", "git-repos-*")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		fmt.Printf("Cleaning up temporary directory: %s\n", tempDir)
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Printf("Warning: Failed to remove temporary directory %s: %v", tempDir, err)
		}
	}()
	fmt.Printf("Cloning repositories into: %s\n", tempDir)

	commitAuthors := make(map[string]int)

	for i, repo := range reposToProcess {
		// Common filters (can be adjusted based on platform nuances)
		if repo.IsFork && platform == "github" {
			fmt.Printf("\nSkipping fork (%d/%d): %s\n", i+1, len(reposToProcess), repo.Name)
			continue
		}
		if repo.IsDisabled && platform == "github" { // GitHub specific concept
			fmt.Printf("\nSkipping disabled repository (%d/%d): %s\n", i+1, len(reposToProcess), repo.Name)
			continue
		}
		if repo.CloneURL == "" {
			fmt.Printf("\nSkipping repository/project with no clone URL (%d/%d): %s\n", i+1, len(reposToProcess), repo.Name)
			continue
		}

		fmt.Printf("\nProcessing (%d/%d): %s\n", i+1, len(reposToProcess), repo.Name)
		repoPath := filepath.Join(tempDir, repo.Name) // TODO: Consider sanitizing repo.Name for path safety

		fmt.Printf("  Cloning %s...\n", repo.CloneURL)
		gitRepo, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
			URL:      repo.CloneURL,
			Progress: nil,
			Depth:    0,
		})
		if err != nil {
			log.Printf("  Failed to clone %s: %v. Skipping.", repo.Name, err)
			// Attempt to clean up partially cloned directory to prevent issues on retry or with same name
			_ = os.RemoveAll(repoPath)
			continue
		}

		commitIter, err := gitRepo.Log(&git.LogOptions{All: true})
		if err != nil {
			log.Printf("  Failed to get commit log for %s: %v. Skipping.", repo.Name, err)
			continue
		}

		commitCountInRepo := 0
		err = commitIter.ForEach(func(c *object.Commit) error {
			authorKey := fmt.Sprintf("%s <%s>", strings.TrimSpace(c.Author.Name), strings.TrimSpace(c.Author.Email))
			commitAuthors[authorKey]++
			commitCountInRepo++
			return nil
		})
		if err != nil {
			log.Printf("  Error iterating commits for %s: %v.", repo.Name, err)
		}
		fmt.Printf("  Processed %d commits in %s.\n", commitCountInRepo, repo.Name)
	}

	fmt.Printf("\n--- Commit Author Summary ---\n")
	if len(commitAuthors) == 0 {
		fmt.Println("No commits found or processed.")
	} else {
		for author, count := range commitAuthors {
			fmt.Printf("%s: %d commits\n", author, count)
		}
	}
}
