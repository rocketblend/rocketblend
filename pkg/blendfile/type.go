package blendfile

import (
	"github.com/rocketblend/rocketblend/pkg/core/executable"
)

type (
	RocketFile struct {
		Build    string   `json:"build"`
		ARGS     string   `json:"args"`
		Version  string   `json:"version"`
		Packages []string `json:"packages"`
	}

	BlendFile struct {
		Exec   *executable.Executable `json:"exec"`
		Path   string                 `json:"path"`
		Addons []string               `json:"addons"`
		ARGS   string                 `json:"args"`
	}

	AddonDict struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
)
