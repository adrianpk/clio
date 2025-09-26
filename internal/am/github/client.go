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

// Client implements the am.Client interface using command-line git.
type Client struct {
	am.Core
}

// NewClient creates a new git Client.
func NewClient(core am.Core) *Client {
	return &Client{
		Core: core,
	}
}

func (c *Client) Clone(ctx context.Context, repoURL, localPath string, auth am.GitAuth) error {
	if auth.Method == am.AuthToken {
		u, err := url.Parse(repoURL)
		if err != nil {
			return fmt.Errorf("cannot parse repo URL: %w", err)
		}
		u.User = url.UserPassword("oauth2", auth.Token)
		repoURL = u.String()
	}

	cmd := exec.CommandContext(ctx, "git", "clone", repoURL, localPath)

	return c.runCommand(cmd)
}

func (c *Client) Checkout(ctx context.Context, localRepoPath, branch string, create bool) error {
	args := []string{"checkout"}
	if create {
		args = append(args, "-b")
	}
	args = append(args, branch)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = localRepoPath

	return c.runCommand(cmd)
}

func (c *Client) Add(ctx context.Context, localRepoPath, pathspec string) error {
	cmd := exec.CommandContext(ctx, "git", "add", pathspec)
	cmd.Dir = localRepoPath
	return c.runCommand(cmd)
}

func (c *Client) Commit(ctx context.Context, localRepoPath string, commit am.GitCommit) (string, error) {
	configUserCmd := exec.CommandContext(ctx, "git", "config", "user.name", commit.UserName)
	configUserCmd.Dir = localRepoPath
	if err := c.runCommand(configUserCmd); err != nil {
		return "", fmt.Errorf("cannot set git user name: %w", err)
	}

	configEmailCmd := exec.CommandContext(ctx, "git", "config", "user.email", commit.UserEmail)
	configEmailCmd.Dir = localRepoPath
	if err := c.runCommand(configEmailCmd); err != nil {
		return "", fmt.Errorf("cannot set git user email: %w", err)
	}

	commitCmd := exec.CommandContext(ctx, "git", "commit", "-m", commit.Message)
	commitCmd.Dir = localRepoPath
	if err := c.runCommand(commitCmd); err != nil {
		return "", fmt.Errorf("cannot commit changes: %w", err)
	}

	hashCmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	hashCmd.Dir = localRepoPath
	var out bytes.Buffer
	hashCmd.Stdout = &out
	if err := c.runCommand(hashCmd); err != nil {
		return "", fmt.Errorf("cannot get commit hash: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}

func (c *Client) Push(ctx context.Context, localRepoPath string, auth am.GitAuth) error {
	cmd := exec.CommandContext(ctx, "git", "push")
	cmd.Dir = localRepoPath
	return c.runCommand(cmd)
}

func (c *Client) Status(ctx context.Context, localRepoPath string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmd.Dir = localRepoPath
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := c.runCommand(cmd); err != nil {
		return "", fmt.Errorf("cannot get git status: %w", err)
	}
	return stdout.String(), nil
}

// runCommand is a helper to execute commands and return a detailed error.
func (c *Client) runCommand(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cannot execute command %q: %w: %s", cmd.String(), err, stderr.String())
	}

	return nil
}
