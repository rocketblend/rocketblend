package helpers

import (
	"fmt"
	"path/filepath"
)

func FindFilePathForExt(ext string) (string, error) {
	// Get a list of all files in the current directory.
	files, err := filepath.Glob("*")
	if err != nil {
		return "", fmt.Errorf("failed to list files in current directory: %w", err)
	}

	// Iterate through the list of files and check if any have a .blend extension.
	for _, file := range files {
		if filepath.Ext(file) == ext {
			// Found a .blend file. Return the full path.
			return filepath.Abs(file)
		}
	}

	// No .blend files found. Return an error.
	return "", fmt.Errorf("no files found with given extension in directory")
}
