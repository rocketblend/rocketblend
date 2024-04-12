package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/types"
	"sigs.k8s.io/yaml"
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
	if err := yaml.Unmarshal(f, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %s", err)
	}

	if err := validator.Validate(result); err != nil {
		return nil, err
	}

	return &result, nil
}

func Save[T any](validator types.Validator, filePath string, object *T) error {
	if validator == nil {
		return errors.New("validator is required")
	}

	if err := validator.Validate(object); err != nil {
		return fmt.Errorf("failed to validate object: %w", err)
	}

	bytes, err := yaml.Marshal(&object)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %s", err)
	}

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}
