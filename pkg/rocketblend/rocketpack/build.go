package rocketpack

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/rocketblend/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	BuildSource struct {
		Platform   runtime.Platform `json:"platform" validate:"required"`
		Executable string           `json:"executable" validate:"required"`
		URL        string           `json:"url" validate:"required,url"`
	}

	Build struct {
		Args    string                `json:"args,omitempty"`
		Version *semver.Version       `json:"version,omitempty"`
		Sources []*BuildSource        `json:"sources" validate:"required"`
		Addons  []reference.Reference `json:"addons,omitempty"`
	}
)

func (i *Build) GetSourceForPlatform(platform runtime.Platform) *BuildSource {
	if i.Sources == nil {
		return nil
	}

	for _, s := range i.Sources {
		if s.Platform == platform {
			return s
		}
	}

	return nil
}

func (i *Build) GetDownloadUrl(platform runtime.Platform) (string, error) {
	source := i.GetSourceForPlatform(platform)
	if source == nil {
		return "", fmt.Errorf("failed to find source for platform: %s", platform)
	}

	return source.URL, nil
}

func (i *Build) GetExecutableName(platform runtime.Platform) (string, error) {
	source := i.GetSourceForPlatform(platform)
	if source == nil {
		return "", fmt.Errorf("failed to find source for platform: %s", platform)
	}

	return source.Executable, nil
}
