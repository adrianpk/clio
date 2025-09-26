package am

import (
"context"
)

// GitAuthMethod defines the authentication strategy.
type GitAuthMethod int

const (
AuthToken GitAuthMethod = iota // HTTPS with Personal Access Token
AuthSSH                        // SSH with agent
)

// GitAuth holds authentication details.
type GitAuth struct {
Method GitAuthMethod
Token  string
}

// GitCommit holds details for making a commit.
type GitCommit struct {
UserName  string
UserEmail string
Message   string
}

// GitClient provides a low-level interface for Git operations.
type GitClient interface {
Clone(ctx context.Context, repoURL, localPath string, auth GitAuth) error
Checkout(ctx context.Context, localRepoPath, branch string, create bool) error
Add(ctx context.Context, localRepoPath, pathspec string) error
Commit(ctx context.Context, localRepoPath string, commit GitCommit) (string, error)
Push(ctx context.Context, localRepoPath string, auth GitAuth) error
Status(ctx context.Context, localRepoPath string) (string, error)
}