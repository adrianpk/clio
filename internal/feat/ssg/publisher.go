package ssg

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianpk/clio/internal/am"
)

// PublisherConfig holds all configuration needed for a publishing operation.
type PublisherConfig struct {
	RepoURL      string // Full URL to the GitHub repository
	Branch       string // Target branch for publishing (e.g., "gh-pages")
	PagesSubdir  string // Subdirectory within the repo (e.g., "" for root, "docs")
	Auth         am.GitAuth
	CommitAuthor am.GitCommit
}

// Publisher defines the interface for orchestrating the publishing process.
type Publisher interface {
	// Validate checks if the provided configuration is valid for publishing.
	Validate(cfg PublisherConfig) error

	// Publish takes the source directory (containing generated HTML) and publishes it
	// to the configured GitHub repo.
	Publish(ctx context.Context, cfg PublisherConfig, sourceDir string) (commitURL string, err error)

	// Plan performs a dry-run, showing what changes would be made without
	// actually pushing to the remote.
	Plan(ctx context.Context, cfg PublisherConfig, sourceDir string) (PlanReport, error)
}

// PlanReport summarizes the changes that would be made during a dry-run.
type PlanReport struct {
	Added    []string
	Modified []string
	Removed  []string
	Summary  string
}

// publisher implements the Publisher interface.
type publisher struct {
	am.Core
	gitClient am.GitClient // Injected GitHub client
}

// NewPublisher creates a new Publisher instance.
func NewPublisher(gitClient am.GitClient, opts ...am.Option) *publisher {
	return &publisher{
		Core:      am.NewCore("ssg-pub", opts...),
		gitClient: gitClient,
	}
}

// Validate implementation
func (p *publisher) Validate(cfg PublisherConfig) error {
	// NOTE: See if we can use am.Validator approach here
	if cfg.RepoURL == "" {
		return fmt.Errorf("repo URL cannot be empty")
	}

	if cfg.Branch == "" {
		return fmt.Errorf("publish branch cannot be empty")
	}

	return nil
}

// Publish implementation
func (p *publisher) Publish(ctx context.Context, cfg PublisherConfig, sourceDir string) (commitURL string, err error) {
	p.Log().Info("Starting publish process")

	// Tenp dir
	tempDir, err := os.MkdirTemp("", "clio-publish-*")
	if err != nil {
		return "", fmt.Errorf("cannot create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)
	p.Log().Info("Temp dir created", "path", tempDir)

	// Clone
	if err := p.gitClient.Clone(ctx, cfg.RepoURL, tempDir, cfg.Auth); err != nil {
		return "", fmt.Errorf("cannot clone repo: %w", err)
	}
	p.Log().Info("Repo cloned")

	// Checkout target branch
	if err := p.gitClient.Checkout(ctx, tempDir, cfg.Branch, true); err != nil {
		return "", fmt.Errorf("cannot checkout branch: %w", err)
	}
	p.Log().Info("Checked out branch", "branch", cfg.Branch)

	// Clean and copy source dir content
	targetDir := filepath.Join(tempDir, cfg.PagesSubdir)
	p.Log().Info("Cleaning target directory", "path", targetDir)
	if err := os.RemoveAll(targetDir); err != nil {
		return "", fmt.Errorf("cannot clean target dir: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create target dir: %w", err)
	}

	p.Log().Info("Copying generated site to target directory")
	if err := copyDir(sourceDir, targetDir); err != nil {
		return "", fmt.Errorf("cannot copy site content: %w", err)
	}

	// Stage
	p.Log().Info("Staging changes")
	if err := p.gitClient.Add(ctx, tempDir, "."); err != nil {
		return "", fmt.Errorf("cannot stage changes: %w", err)
	}

	// GitCommit
	p.Log().Info("Committing changes")
	commitHash, err := p.gitClient.Commit(ctx, tempDir, cfg.CommitAuthor)
	if err != nil {
		return "", fmt.Errorf("cannot commit changes: %w", err)
	}
	p.Log().Info("Changes committed", "hash", commitHash)

	// Push
	p.Log().Info("Pushing changes to remote")
	if err := p.gitClient.Push(ctx, tempDir, cfg.Auth); err != nil {
		return "", fmt.Errorf("cannot push changes: %w", err)
	}

	// // NOTE: We need to find a neater way to do this
	commitURL = fmt.Sprintf("%s/commit/%s", cfg.RepoURL, commitHash)
	p.Log().Info("Publish process completed successfully", "commit_url", commitURL)

	return commitURL, nil
}

// Plan implementation
func (p *publisher) Plan(ctx context.Context, cfg PublisherConfig, sourceDir string) (PlanReport, error) {
	p.Log().Info("Starting plan dry-run process")

	var report PlanReport

	// Tenp dir
	tempDir, err := os.MkdirTemp("", "clio-plan-*")
	if err != nil {
		return PlanReport{}, fmt.Errorf("cannot create temp dir for plan: %w", err)
	}
	defer os.RemoveAll(tempDir)
	p.Log().Info("Temp dir created for plan", "path", tempDir)

	// Clone
	if err := p.gitClient.Clone(ctx, cfg.RepoURL, tempDir, cfg.Auth); err != nil {
		return PlanReport{}, fmt.Errorf("cannot clone repo for plan: %w", err)
	}
	p.Log().Info("Repo cloned for plan")

	// Checkout target branch
	if err := p.gitClient.Checkout(ctx, tempDir, cfg.Branch, true); err != nil {
		return PlanReport{}, fmt.Errorf("cannot checkout branch for plan: %w", err)
	}
	p.Log().Info("Checked out branch for plan", "branch", cfg.Branch)

	// Clean and copy source dir content
	targetDir := filepath.Join(tempDir, cfg.PagesSubdir)
	p.Log().Info("Cleaning target directory for plan", "path", targetDir)
	if err := os.RemoveAll(targetDir); err != nil {
		return PlanReport{}, fmt.Errorf("cannot clean target dir for plan: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return PlanReport{}, fmt.Errorf("cannot create target dir for plan: %w", err)
	}

	p.Log().Info("Copying generated site to target directory for plan")
	if err := copyDir(sourceDir, targetDir); err != nil {
		return PlanReport{}, fmt.Errorf("cannot copy site content for plan: %w", err)
	}

	// Stage
	p.Log().Info("Staging changes for plan")
	if err := p.gitClient.Add(ctx, tempDir, "."); err != nil {
		return PlanReport{}, fmt.Errorf("cannot stage changes for plan: %w", err)
	}

	// Get status
	p.Log().Info("Getting git status for plan")
	statusOutput, err := p.gitClient.Status(ctx, tempDir)
	if err != nil {
		return PlanReport{}, fmt.Errorf("cannot get git status for plan: %w", err)
	}

	// Parse status output
	lines := strings.Split(statusOutput, "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		status := line[0:2]
		filename := strings.TrimSpace(line[3:])

		switch status {
		case "A ":
			report.Added = append(report.Added, filename)
		case "M ":
			report.Modified = append(report.Modified, filename)
		case "D ":
			report.Removed = append(report.Removed, filename)
		case "??": // Untracked files, which should be added by `git add .`
			report.Added = append(report.Added, filename)
		}
	}

	report.Summary = fmt.Sprintf("Added: %d, Modified: %d, Removed: %d", len(report.Added), len(report.Modified), len(report.Removed))
	p.Log().Info("Plan dry-run process completed successfully", "summary", report.Summary)

	return report, nil
}

// copyDir copies the contents of src to dst. It is not recursive!
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}
