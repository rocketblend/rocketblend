package library

import "github.com/rocketblend/rocketblend/pkg/core/runtime"

type (
	Source struct {
		Platform   runtime.Platform `json:"platform"`
		Executable string           `json:"executable"`
		URL        string           `json:"url"`
	}

	Build struct {
		Args     string    `json:"args"`
		Source   []*Source `json:"source"`
		Packages []string  `json:"packages"`
	}
)

type (
	Package struct {
		Source string `json:"source"`
	}
)

type (
	Install struct {
		Path string `json:"path"`
	}

	Pack struct {
		Path string `json:"path"`
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
