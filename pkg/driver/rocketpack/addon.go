package rocketpack

import (
	"net/url"

	"github.com/rocketblend/rocketblend/pkg/driver/source"
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

func (a *Addon) IsLocalOnly() bool {
	return a.IsPreInstalled() || a.Source.URL == ""
}

func (a *Addon) IsPreInstalled() bool {
	return a.Source == nil
}

func (a *Addon) GetSource() (*source.Source, error) {
	var url *url.URL
	if !a.IsPreInstalled() {
		var err error
		url, err = url.Parse(a.Source.URL)
		if err != nil {
			return nil, err
		}
	}

	return &source.Source{
		FileName: a.Source.File,
		URI:      url,
	}, nil
}

func (a *Addon) GetDownloadUrl() (string, error) {
	if a.IsPreInstalled() {
		return "", nil
	}

	return a.Source.URL, nil
}

func (a *Addon) GetExecutableName() (string, error) {
	if a.IsPreInstalled() {
		return "", nil
	}

	return a.Source.File, nil
}
