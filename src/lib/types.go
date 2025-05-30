package lib

import "time"

type RepoInfo struct {
	Name         string
	CloneURL     string
	IsFork       bool
	IsDisabled   bool // Not supported in GitLab, Codeberg, or SourceHut
	LastActivity time.Time
}
