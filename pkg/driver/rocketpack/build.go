package rocketpack

import (
	"fmt"
	"net/url"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/driver/source"
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

func (i *Build) IsLocalOnly(platform runtime.Platform) (bool, error) {
	source := i.GetSourceForPlatform(platform)
	if source == nil {
		return false, fmt.Errorf("failed to find source for platform: %s", platform)
	}

	return source.URL == "", nil
}

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

func (i *Build) GetSource(platform runtime.Platform) (*source.Source, error) {
	buildSource := i.GetSourceForPlatform(platform)
	if buildSource == nil {
		return nil, fmt.Errorf("failed to find source for platform: %s", platform)
	}

	url, err := url.Parse(buildSource.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s", err)
	}

	return &source.Source{
		FileName: buildSource.Executable,
		URI:      url,
	}, nil
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