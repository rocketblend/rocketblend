package blendconfig

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/rocketblend2/helpers"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketfile"
)

const BlenderFileExtension = ".blend"

type (
	BlendConfig struct {
		BlendFilePath string                 `json:"blendFilepath"`
		RocketFile    *rocketfile.RocketFile `json:"rocketFile,omitempty"`
	}
)

func New(blendFilePath string, RocketFile *rocketfile.RocketFile) (*BlendConfig, error) {
	blendFile := &BlendConfig{
		BlendFilePath: blendFilePath,
		RocketFile:    RocketFile,
	}

	err := Validate(blendFile)
	if err != nil {
		return nil, fmt.Errorf("failed to validate blend file: %s", err)
	}

	return blendFile, nil
}

func Load(blendFilePath string, rocketFilePath string) (*BlendConfig, error) {
	err := validateBlendConfigPath(blendFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to validate blend file path: %s", err)
	}

	err = helpers.FileExists(blendFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to find blend file: %s", err)
	}

	rocketfile, err := rocketfile.Load(rocketFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to validate rocket file path: %s", err)
	}

	return &BlendConfig{
		BlendFilePath: blendFilePath,
		RocketFile:    rocketfile,
	}, nil
}

func Validate(blendFile *BlendConfig) error {
	if err := validateBlendConfigPath(blendFile.BlendFilePath); err != nil {
		return fmt.Errorf("failed to validate blend file path: %w", err)
	}

	if err := rocketfile.Validate(blendFile.RocketFile); err != nil {
		return fmt.Errorf("failed to validate rocket file: %w", err)
	}

	return nil
}

func validateBlendConfigPath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("blend file path cannot be empty")
	}

	ext := filepath.Ext(filePath)
	if ext != BlenderFileExtension {
		return fmt.Errorf("invalid file extension (must be '%s'): %s", BlenderFileExtension, ext)
	}

	return nil
}
