package build

import (
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Source struct {
		Platform   runtime.Platform `json:"platform"`
		Executable string           `json:"executable"`
		URL        string           `json:"url"`
	}

	Build struct {
		Args           string          `json:"args"`
		BlenderVersion *semver.Version `json:"blenderVersion"`
		Source         []*Source       `json:"source"`
		Addons         []string        `json:"addons"`
	}
)

func (i *Build) GetSourceForPlatform(platform runtime.Platform) *Source {
	for _, s := range i.Source {
		if s.Platform == platform {
			return s
		}
	}

	return nil
}
