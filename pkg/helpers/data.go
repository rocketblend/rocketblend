package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/types"
)

func Load[T any](validator types.Validator, filePath string) (*T, error) {
	if validator == nil {
		return nil, errors.New("validator is required")
	}

	if err := FileExists(filePath); err != nil {
		return nil, fmt.Errorf("failed to find file: %s", err)
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	var result T
	if err := json.Unmarshal(f, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %s", err)
	}

	if err := validator.Validate(result); err != nil {
		return nil, err
	}

	return &result, nil
}

func Save[T any](validator types.Validator, filePath string, ensurePath bool, object *T) error {
	if validator == nil {
		return errors.New("validator is required")
	}

	if err := validator.Validate(object); err != nil {
		return fmt.Errorf("failed to validate object: %w", err)
	}

	bytes, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal object: %s", err)
	}

	if ensurePath {
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
	}

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}
