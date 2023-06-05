package blendfile

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/driver/helpers"
	"github.com/rocketblend/rocketblend/pkg/driver/installation"
)

type (
	BlendFile struct {
		ProjectName string                `json:"projectName"`
		FilePath    string                `json:"filePath"`
		Build       *installation.Build   `json:"build"`
		Addons      []*installation.Addon `json:"addons"`
		ARGS        string                `json:"args"`
	}
)

func Validate(blendFile *BlendFile) error {
	if blendFile == nil {
		return fmt.Errorf("rocketfile cannot be nil")
	}

	if blendFile.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if err := helpers.FileExists(blendFile.FilePath); err != nil {
		return fmt.Errorf("failed to find blend file: %s", err)
	}

	if err := validateBuild(blendFile.Build); err != nil {
		return err
	}

	for _, addon := range blendFile.Addons {
		if err := validateAddon(addon); err != nil {
			return err
		}
	}

	return nil
}

func validateBuild(build *installation.Build) error {
	if build == nil {
		return fmt.Errorf("build cannot be nil")
	}

	if build.FilePath == "" {
		return fmt.Errorf("build file path cannot be empty")
	}

	if err := helpers.FileExists(build.FilePath); err != nil {
		return fmt.Errorf("failed to blender executable: %s", err)
	}

	return nil
}

func validateAddon(addon *installation.Addon) error {
	if addon == nil {
		return fmt.Errorf("addon cannot be nil")
	}

	if addon.Name == "" {
		return fmt.Errorf("addon name cannot be empty")
	}

	if addon.FilePath != "" {
		// If the addon file path is not empty, then we need to make sure it exists.
		// If it is empty, then we assume the addon is pre-installed.
		if err := helpers.FileExists(addon.FilePath); err != nil {
			return fmt.Errorf("failed to find addon file (%s@%s): %s", addon.Name, addon.FilePath, err)
		}
	}

	return nil
}
