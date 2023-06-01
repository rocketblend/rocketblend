package rocketpack

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	AddonSource struct {
		File string `json:"file" validate:"required"`
		URL  string `json:"url,omitempty" validate:"url"`
	}

	Addon struct {
		Name    string          `json:"name" validate:"required"`
		Version *semver.Version `json:"version,omitempty"`
		Source  *AddonSource    `json:"source,omitempty"`
	}
)

func (a *Addon) GetDownloadUrl() (string, error) {
	if a.Source == nil {
		return "", fmt.Errorf("failed to find source for addon: %s", a.Name)
	}

	return a.Source.URL, nil
}

func (a *Addon) GetExecutableName() (string, error) {
	if a.Source == nil {
		return "", fmt.Errorf("failed to find source for addon: %s", a.Name)
	}

	return a.Source.File, nil
}
