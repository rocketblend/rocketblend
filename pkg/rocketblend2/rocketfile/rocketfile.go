package rocketfile

import (
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/helpers"
	"sigs.k8s.io/yaml"
)

const FileName = "rocketfile.yaml"

type (
	RocketFile struct {
		Build   reference.Reference   `json:"build"`
		ARGS    string                `json:"args,omitempty"`
		Version string                `json:"version,omitempty"`
		Addons  []reference.Reference `json:"addons,omitempty"`
	}
)

func Load(filePath string) (*RocketFile, error) {
	if err := validateFilePath(filePath); err != nil {
		return nil, fmt.Errorf("failed to validate file path: %s", err)
	}

	if err := helpers.FileExists(filePath); err != nil {
		return nil, fmt.Errorf("failed to find blend file: %s", err)
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	var rocketFile RocketFile
	if err := yaml.Unmarshal(f, &rocketFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	if err := Validate(&rocketFile); err != nil {
		return nil, fmt.Errorf("failed to validate rocketfile: %w", err)
	}

	return &rocketFile, nil
}

func Save(filePath string, rocketfile *RocketFile) error {
	err := Validate(rocketfile)
	if err != nil {
		return fmt.Errorf("failed to validate rocketfile: %s", err)
	}

	err = validateFilePath(filePath)
	if err != nil {
		return fmt.Errorf("failed to validate file path: %s", err)
	}

	f, err := yaml.Marshal(rocketfile)
	if err != nil {
		return fmt.Errorf("failed to marshal rocketfile: %s", err)
	}

	if err := os.WriteFile(filePath, f, 0644); err != nil {
		return fmt.Errorf("failed to write rocketfile: %s", err)
	}

	return nil
}

func Validate(r *RocketFile) error {
	return nil
}

func validateFilePath(filePath string) error {
	return helpers.ValidateFilePath(filePath, FileName)
}
