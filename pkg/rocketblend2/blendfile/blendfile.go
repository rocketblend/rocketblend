package blendfile

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/rocketblend2/helpers"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Addon struct {
		FilePath string         `json:"filePath"`
		Name     string         `json:"name"`
		Version  semver.Version `json:"version"`
	}

	Build struct {
		FilePath string   `json:"filePath"`
		Addons   []*Addon `json:"addons"`
		ARGS     string   `json:"args"`
	}

	BlendFile struct {
		FilePath string   `json:"filePath"`
		Build    *Build   `json:"build"`
		Addons   []*Addon `json:"addons"`
		ARGS     string   `json:"args"`
	}
)

func Validate(rocketfile *BlendFile) error {
	if rocketfile == nil {
		return fmt.Errorf("rocketfile cannot be nil")
	}

	if err := helpers.FileExists(rocketfile.FilePath); err != nil {
		return fmt.Errorf("failed to find blend file: %s", err)
	}

	if err := validateBuild(rocketfile.Build); err != nil {
		return err
	}

	for _, addon := range rocketfile.Addons {
		if err := validateAddon(addon); err != nil {
			return err
		}
	}

	return nil
}

func validateBuild(build *Build) error {
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

func validateAddon(addon *Addon) error {
	if addon == nil {
		return fmt.Errorf("addon cannot be nil")
	}

	if addon.FilePath == "" {
		return fmt.Errorf("addon file path cannot be empty")
	}

	if addon.Name == "" {
		return fmt.Errorf("addon name cannot be empty")
	}

	if err := helpers.FileExists(addon.FilePath); err != nil {
		return fmt.Errorf("failed to find addon file: %s", err)
	}

	return nil
}
