package command

import (
	"fmt"
	"path/filepath"
)

func findFilePathForExt(dir string, ext string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		return "", fmt.Errorf("failed to list files in current directory: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in directory")
	}

	return files[0], nil
}
