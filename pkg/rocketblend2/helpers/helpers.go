package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

func ValidateFilePath(filePath string, requiredFileName string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if filepath.Base(filePath) != requiredFileName && requiredFileName != "" {
		return fmt.Errorf("invalid file name (must be '%s'): %s", requiredFileName, filepath.Base(filePath))
	}

	return nil
}

func FileExists(filePath string) error {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}

	if info.IsDir() {
		return fmt.Errorf("file is a directory")
	}

	return nil
}
