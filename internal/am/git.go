package am

import (
	"context"
)

// GitClient defines the interface for Git operations.
type GitClient interface {
	Clone(ctx context.Context, repoURL, localPath string, auth GitAuth, env []string) error
	Checkout(ctx context.Context, localRepoPath, branch string, create bool, env []string) error
	Add(ctx context.Context, localRepoPath, pathspec string, env []string) error
	Commit(ctx context.Context, localRepoPath string, commit GitCommit, env []string) (string, error)
	Push(ctx context.Context, localRepoPath string, auth GitAuth, remote, branch string, env []string) error
	Status(ctx context.Context, localRepoPath string, env []string) (string, error)
	GitLog(ctx context.Context, localRepoPath string, args []string, env []string) (string, error) // args is now a slice
}

// GitAuth holds authentication details for Git operations.
type GitAuth struct {
	Method AuthMethod
	Token  string
}

// AuthMethod represents the authentication method.
type AuthMethod string

const (
	AuthToken AuthMethod = "token"
	AuthSSH   AuthMethod = "ssh"
)

// GitCommit holds commit details.
type GitCommit struct {
	UserName  string
	UserEmail string
	Message   string
}
