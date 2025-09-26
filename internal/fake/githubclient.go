package fake

import (
	"context"

	"github.com/adrianpk/clio/internal/am"
)

// GithubClient is a fake implementation of am.GitClient for testing.
type GithubClient struct {
	// Expected results
	CloneFn    func(ctx context.Context, repoURL, localPath string, auth am.GitAuth) error
	CheckoutFn func(ctx context.Context, localRepoPath, branch string, create bool) error
	AddFn      func(ctx context.Context, localRepoPath, pathspec string) error
	CommitFn   func(ctx context.Context, localRepoPath string, commit am.GitCommit) (string, error)
	PushFn     func(ctx context.Context, localRepoPath string, auth am.GitAuth) error
	StatusFn   func(ctx context.Context, localRepoPath string) (string, error)

	// Captured arguments
	CloneCalls []struct {
		Ctx                context.Context
		RepoURL, LocalPath string
		Auth               am.GitAuth
	}
	CheckoutCalls []struct {
		Ctx                   context.Context
		LocalRepoPath, Branch string
		Create                bool
	}
	AddCalls []struct {
		Ctx                     context.Context
		LocalRepoPath, Pathspec string
	}
	CommitCalls []struct {
		Ctx           context.Context
		LocalRepoPath string
		Commit        am.GitCommit
	}
	PushCalls []struct {
		Ctx           context.Context
		LocalRepoPath string
		Auth          am.GitAuth
	}
	StatusCalls []struct {
		Ctx           context.Context
		LocalRepoPath string
	}
}

// NewGithubClient creates a new fake GithubClient.
func NewGithubClient() *GithubClient {
	return &GithubClient{}
}

func (f *GithubClient) Clone(ctx context.Context, repoURL, localPath string, auth am.GitAuth) error {
	f.CloneCalls = append(f.CloneCalls, struct {
		Ctx                context.Context
		RepoURL, LocalPath string
		Auth               am.GitAuth
	}{Ctx: ctx, RepoURL: repoURL, LocalPath: localPath, Auth: auth})
	if f.CloneFn != nil {
		return f.CloneFn(ctx, repoURL, localPath, auth)
	}
	return nil // Default success
}

func (f *GithubClient) Checkout(ctx context.Context, localRepoPath, branch string, create bool) error {
	f.CheckoutCalls = append(f.CheckoutCalls, struct {
		Ctx                   context.Context
		LocalRepoPath, Branch string
		Create                bool
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Branch: branch, Create: create})
	if f.CheckoutFn != nil {
		return f.CheckoutFn(ctx, localRepoPath, branch, create)
	}
	return nil
}

func (f *GithubClient) Add(ctx context.Context, localRepoPath, pathspec string) error {
	f.AddCalls = append(f.AddCalls, struct {
		Ctx                     context.Context
		LocalRepoPath, Pathspec string
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Pathspec: pathspec})
	if f.AddFn != nil {
		return f.AddFn(ctx, localRepoPath, pathspec)
	}
	return nil
}

func (f *GithubClient) Commit(ctx context.Context, localRepoPath string, commit am.GitCommit) (string, error) {
	f.CommitCalls = append(f.CommitCalls, struct {
		Ctx           context.Context
		LocalRepoPath string
		Commit        am.GitCommit
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Commit: commit})
	if f.CommitFn != nil {
		return f.CommitFn(ctx, localRepoPath, commit)
	}
	return "fake-commit-hash", nil // Default hash
}

func (f *GithubClient) Push(ctx context.Context, localRepoPath string, auth am.GitAuth) error {
	f.PushCalls = append(f.PushCalls, struct {
		Ctx           context.Context
		LocalRepoPath string
		Auth          am.GitAuth
	}{Ctx: ctx, LocalRepoPath: localRepoPath, Auth: auth})
	if f.PushFn != nil {
		return f.PushFn(ctx, localRepoPath, auth)
	}
	return nil
}

func (f *GithubClient) Status(ctx context.Context, localRepoPath string) (string, error) {
	f.StatusCalls = append(f.StatusCalls, struct {
		Ctx           context.Context
		LocalRepoPath string
	}{Ctx: ctx, LocalRepoPath: localRepoPath})
	if f.StatusFn != nil {
		return f.StatusFn(ctx, localRepoPath)
	}
	return " M somefile.txt\n?? anotherfile.txt", nil // Default status
}