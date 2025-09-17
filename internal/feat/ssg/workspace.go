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
	var dirs []string
	env := w.Cfg().StrValOrDef(am.Key.AppEnv, "prod")
	w.Log().Info("Read environment mode", "key", am.Key.AppEnv, "value", env)

	if env == "dev" {
		w.Log().Info("Running in DEV mode, using local paths.")
		wd, err := os.Getwd()
		if err != nil {
			w.Log().Error("Cannot get working directory", "error", err)
			return err
		}
		base := filepath.Join(wd, "_workspace")
		dbDir := filepath.Join(base, "db")
		dirs = []string{
			filepath.Join(base, "config"),
			dbDir,
			filepath.Join(base, "documents", "markdown"),
			filepath.Join(base, "documents", "html"),
			filepath.Join(base, "documents", "assets", "images"),
		}

		// Override config values for dev mode
		devDSN := "file:" + filepath.Join(dbDir, "clio.db") + "?cache=shared&mode=rwc"
		w.Cfg().Set(am.Key.DBSQLiteDSN, devDSN)
		w.Log().Info("Overriding config for DEV mode", "key", am.Key.DBSQLiteDSN, "value", devDSN)

	} else {
		w.Log().Info("Running in PROD mode, using system paths.")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			w.Log().Error("Cannot get user home directory", "error", err)
			return err
		}
		dirs = []string{
			filepath.Join(homeDir, ".config", "clio"),
			filepath.Join(homeDir, ".clio"),
			filepath.Join(homeDir, "Documents", "Clio", "markdown"),
			filepath.Join(homeDir, "Documents", "Clio", "html"),
			filepath.Join(homeDir, "Documents", "Clio", "assets", "images"),
		}
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