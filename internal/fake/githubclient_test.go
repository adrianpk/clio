package fake_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adrianpk/clio/internal/fake"
	"github.com/adrianpk/clio/internal/git/github"
)

func TestGithubClientClone(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.GithubClient)
		ctx         context.Context
		repoURL     string
		localPath   string
		auth        github.Auth
		expectedErr error
		expectCalls int
	}{
		{
			name:        "successful clone",
			setupFake:   func(f *fake.GithubClient) {},
			ctx:         context.Background(),
			repoURL:     "https://github.com/owner/repo.git",
			localPath:   "/tmp/repo",
			auth:        github.Auth{Method: github.AuthToken, Token: "test-token"},
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "clone returns error",
			setupFake: func(f *fake.GithubClient) {
				f.CloneFn = func(ctx context.Context, repoURL, localPath string, auth github.Auth) error {
					return errors.New("clone failed")
				}
			},
			ctx:         context.Background(),
			repoURL:     "https://github.com/owner/repo.git",
			localPath:   "/tmp/repo",
			auth:        github.Auth{Method: github.AuthSSH},
			expectedErr: errors.New("clone failed"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewGithubClient()
			tt.setupFake(f)

			err := f.Clone(tt.ctx, tt.repoURL, tt.localPath, tt.auth)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.CloneCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.CloneCalls))
			}
			if tt.expectCalls > 0 {
				call := f.CloneCalls[0]
				if call.RepoURL != tt.repoURL || call.LocalPath != tt.localPath || call.Auth != tt.auth {
					t.Errorf("captured call arguments mismatch")
				}
			}
		})
	}
}
