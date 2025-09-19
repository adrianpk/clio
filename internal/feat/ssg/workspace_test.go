package ssg_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/adrianpk/clio/internal/am"
	"github.com/adrianpk/clio/internal/feat/ssg"
)

func TestWorkspaceSetup(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get user home directory: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "clio-test-*")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("could not change to temp dir: %v", err)
	}
	defer os.Chdir(originalWd)

	testCases := []struct {
		name          string
		env           string
		expectedPaths map[string]string
	}{
		{
			name: "dev mode",
			env:  "dev",
			expectedPaths: map[string]string{
				am.Key.DBSQLiteDSN:      "file:" + filepath.Join(tempDir, "_workspace", "db", "clio.db") + "?cache=shared&mode=rwc",
				am.Key.SSGWorkspacePath: filepath.Join(tempDir, "_workspace"),
				am.Key.SSGDocsPath:      filepath.Join(tempDir, "_workspace", "documents"),
				am.Key.SSGMarkdownPath:  filepath.Join(tempDir, "_workspace", "documents", "markdown"),
				am.Key.SSGHTMLPath:      filepath.Join(tempDir, "_workspace", "documents", "html"),
				am.Key.SSGAssetsPath:    filepath.Join(tempDir, "_workspace", "documents", "assets"),
				am.Key.SSGImagesPath:    filepath.Join(tempDir, "_workspace", "documents", "assets", "images"),
			},
		},
		{
			name: "prod mode",
			env:  "prod",
			expectedPaths: map[string]string{
				am.Key.SSGWorkspacePath: filepath.Join(homeDir, ".clio"),
				am.Key.SSGDocsPath:      filepath.Join(homeDir, "Documents", "Clio"),
				am.Key.SSGMarkdownPath:  filepath.Join(homeDir, "Documents", "Clio", "markdown"),
				am.Key.SSGHTMLPath:      filepath.Join(homeDir, "Documents", "Clio", "html"),
				am.Key.SSGAssetsPath:    filepath.Join(homeDir, "Documents", "Clio", "assets"),
				am.Key.SSGImagesPath:    filepath.Join(homeDir, "Documents", "Clio", "assets", "images"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := am.NewConfig()
			cfg.Set(am.Key.AppEnv, tc.env)

			logger := am.NewLogger("")
			ws := ssg.NewWorkspace(am.WithCfg(cfg), am.WithLog(logger))

			if err := ws.Setup(context.Background()); err != nil {
				t.Fatalf("ws.Setup() failed: %v", err)
			}

			for key, expectedPath := range tc.expectedPaths {
				actualPath := cfg.StrValOrDef(key, "")
				if actualPath != expectedPath {
					t.Errorf("config value for key %q: got %q, want %q", key, actualPath, expectedPath)
				}
			}

			if tc.name == "dev mode" {
				for key, path := range tc.expectedPaths {
					if key == am.Key.DBSQLiteDSN {
						continue
					}
					if _, err := os.Stat(path); os.IsNotExist(err) {
						t.Errorf("directory %q should have been created in dev mode", path)
					}
				}
			}
		})
	}
}