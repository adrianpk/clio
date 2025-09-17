package ssg

import (
	"context"
	"os"
	"path/filepath"

	"github.com/adrianpk/clio/internal/am"
)

type Workspace struct {
	am.Core
}

func NewWorkspace(opts ...am.Option) *Workspace {
	core := am.NewCore("ssg-workspace", opts...)
	w := &Workspace{
		Core: core,
	}
	return w
}

func (w *Workspace) Setup(ctx context.Context) error {
	return w.setupDirs()
}

func (w *Workspace) setupDirs() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		w.Log().Error("Cannot get user home directory", "error", err)
		return err
	}

	// TODO: These paths should eventually come from a configuration file.
	dirs := []string{
		filepath.Join(homeDir, ".config", "clio"),
		filepath.Join(homeDir, ".clio"),
		filepath.Join(homeDir, "Documents", "Clio", "markdown"),
		filepath.Join(homeDir, "Documents", "Clio", "html"),
		filepath.Join(homeDir, "Documents", "Clio", "assets", "images"),
	}

	w.Log().Info("Ensuring base directory structure exists...")
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			w.Log().Error("Error creating directory", "path", dir, "error", err)
			return err
		}
	}
	w.Log().Info("Base directory structure verified.")

	return nil
}
