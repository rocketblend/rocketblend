package installation

import (
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Addon struct {
		FilePath string          `json:"filePath"`
		Name     string          `json:"name"`
		Version  *semver.Version `json:"version"`
	}

	Build struct {
		FilePath string `json:"filePath"`
		ARGS     string `json:"args"`
	}

	Installation struct {
		Addon *Addon `json:"addon"`
		Build *Build `json:"build"`
	}
)

func (i *Installation) IsAddon() bool {
	return i.Addon != nil
}

func (i *Installation) IsBuild() bool {
	return i.Build != nil
}
