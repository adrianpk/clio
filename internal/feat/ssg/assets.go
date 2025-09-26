package ssg

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyStaticAssets(assetsFS embed.FS, targetDir string) error {

	return fs.WalkDir(assetsFS, "assets/ssg/static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		destPath := filepath.Join(targetDir, path[len("assets/ssg"):])

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("cannot create directory: %w", err)
			}
			return nil
		}

		return copyFile(assetsFS, path, destPath)
	})
}

func copyFile(assetsFS embed.FS, srcPath, dstPath string) error {
	srcFile, err := assetsFS.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("cannot copy file: %w", err)
	}

	return nil
}
