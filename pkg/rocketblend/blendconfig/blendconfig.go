package blendconfig

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/rocketblend/helpers"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketfile"
)

const BlenderFileExtension = ".blend"

type (
	BlendConfig struct {
		ProjectPath   string                 `json:"projectPath"`
		BlendFileName string                 `json:"blendFileName"`
		RocketFile    *rocketfile.RocketFile `json:"rocketFile"`
	}
)

func (blendFile *BlendConfig) BlendFilePath() string {
	return filepath.Join(blendFile.ProjectPath, blendFile.BlendFileName)
}

func New(ProjectPath string, blendFileName string, RocketFile *rocketfile.RocketFile) (*BlendConfig, error) {
	blendFile := &BlendConfig{
		ProjectPath:   ProjectPath,
		BlendFileName: blendFileName,
		RocketFile:    RocketFile,
	}

	err := Validate(blendFile)
	if err != nil {
		return nil, fmt.Errorf("failed to validate blend file: %s", err)
	}

	return blendFile, nil
}

func Load(blendFilePath string, rocketFilePath string) (*BlendConfig, error) {
	err := validateBlendFilePath(blendFilePath)
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
		ProjectPath:   filepath.Dir(blendFilePath),
		BlendFileName: filepath.Base(blendFilePath),
		RocketFile:    rocketfile,
	}, nil
}

func Validate(blendFile *BlendConfig) error {
	if err := validateProjectPath(blendFile.ProjectPath); err != nil {
		return fmt.Errorf("failed to validate project path: %w", err)
	}

	if err := validateBlendFilePath(blendFile.BlendFileName); err != nil {
		return fmt.Errorf("failed to validate blend file path: %w", err)
	}

	if err := rocketfile.Validate(blendFile.RocketFile); err != nil {
		return fmt.Errorf("failed to validate rocket file: %w", err)
	}

	return nil
}

func validateProjectPath(projectPath string) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	return nil
}

func validateBlendFilePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("blend file path cannot be empty")
	}

	ext := filepath.Ext(filePath)
	if ext != BlenderFileExtension {
		return fmt.Errorf("invalid file extension (must be '%s'): %s", BlenderFileExtension, ext)
	}

	return nil
}
