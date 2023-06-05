package helpers

import (
	"fmt"
	"path/filepath"
	"sort"
)

func RemoveDuplicateStr(strs []string) []string {
	sort.Strings(strs)
	for i := len(strs) - 1; i > 0; i-- {
		if strs[i] == strs[i-1] {
			strs = append(strs[:i], strs[i+1:]...)
		}
	}

	return strs
}

func FindFilePathForExt(dir string, ext string) (string, error) {
	// Get a list of all files in the current directory.

	files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		return "", fmt.Errorf("failed to list files in current directory: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in directory")
	}

	return files[0], nil
}
