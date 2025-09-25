package github

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/adrianpk/clio/internal/am"
)

// AuthMethod defines the authentication strategy.
type AuthMethod int

const (
	AuthToken AuthMethod = iota // HTTPS with Personal Access Token
	AuthSSH                     // SSH with agent
)

// Auth holds authentication details.
type Auth struct {
	Method AuthMethod
	Token  string
}

// Commit holds details for making a commit.
type Commit struct {
	UserName  string
	UserEmail string
	Message   string
}

// Client provides a low-level interface for Git operations.
type Client interface {
	Clone(ctx context.Context, repoURL, localPath string, auth Auth) error
	Checkout(ctx context.Context, localRepoPath, branch string, create bool) error
	Add(ctx context.Context, localRepoPath, pathspec string) error
	Commit(ctx context.Context, localRepoPath string, commit Commit) (string, error)
	Push(ctx context.Context, localRepoPath string, auth Auth) error
	Status(ctx context.Context, localRepoPath string) (string, error)
}

// client implements the Client interface using command-line git.
type client struct {
	am.Core
}

// NewClient creates a new git client.
func NewClient(opts ...am.Option) *client {
	return &client{
		Core: am.NewCore("git-client", opts...),
	}
}

func (c *client) Clone(ctx context.Context, repoURL, localPath string, auth Auth) error {
	if auth.Method == AuthToken {
		u, err := url.Parse(repoURL)
		if err != nil {
			return fmt.Errorf("invalid repo URL: %w", err)
		}
		u.User = url.UserPassword("oauth2", auth.Token)
		repoURL = u.String()
	}

	cmd := exec.CommandContext(ctx, "git", "clone", repoURL, localPath)

	return c.runCommand(cmd)
}

func (c *client) Checkout(ctx context.Context, localRepoPath, branch string, create bool) error {
	args := []string{"checkout"}
	if create {
		args = append(args, "-b")
	}
	args = append(args, branch)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = localRepoPath

	return c.runCommand(cmd)
}

func (c *client) Add(ctx context.Context, localRepoPath, pathspec string) error {
	cmd := exec.CommandContext(ctx, "git", "add", pathspec)
	cmd.Dir = localRepoPath
	return c.runCommand(cmd)
}

func (c *client) Commit(ctx context.Context, localRepoPath string, commit Commit) (string, error) {
	configUserCmd := exec.CommandContext(ctx, "git", "config", "user.name", commit.UserName)
	configUserCmd.Dir = localRepoPath
	if err := c.runCommand(configUserCmd); err != nil {
		return "", fmt.Errorf("failed to set git user name: %w", err)
	}

	configEmailCmd := exec.CommandContext(ctx, "git", "config", "user.email", commit.UserEmail)
	configEmailCmd.Dir = localRepoPath
	if err := c.runCommand(configEmailCmd); err != nil {
		return "", fmt.Errorf("failed to set git user email: %w", err)
	}

	commitCmd := exec.CommandContext(ctx, "git", "commit", "-m", commit.Message)
	commitCmd.Dir = localRepoPath
	if err := c.runCommand(commitCmd); err != nil {
		return "", fmt.Errorf("git commit failed: %w", err)
	}

	hashCmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	hashCmd.Dir = localRepoPath
	var out bytes.Buffer
	hashCmd.Stdout = &out
	if err := c.runCommand(hashCmd); err != nil {
		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}

func (c *client) Push(ctx context.Context, localRepoPath string, auth Auth) error {
	cmd := exec.CommandContext(ctx, "git", "push")
	cmd.Dir = localRepoPath
	return c.runCommand(cmd)
}

func (c *client) Status(ctx context.Context, localRepoPath string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmd.Dir = localRepoPath
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := c.runCommand(cmd); err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}
	return stdout.String(), nil
}

// runCommand is a helper to execute commands and return a detailed error.
func (c *client) runCommand(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(`error executing command: %s\nerror: %w\noutput: %s`, cmd.String(), err, stderr.String())
	}

	return nil
}
