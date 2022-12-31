package library

import (
	"github.com/Masterminds/semver/v3"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
)

type (
	Source struct {
		Platform   runtime.Platform `json:"platform"`
		Executable string           `json:"executable"`
		URL        string           `json:"url"`
	}

	Build struct {
		Reference      string         `json:"reference"`
		Args           string         `json:"args"`
		BlenderVersion semver.Version `json:"blenderVersion"`
		Source         []*Source      `json:"source"`
		Packages       []string       `json:"packages"`
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
