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
		Addons map[string]string      `json:"addons"`
		ARGS   string                 `json:"args"`
	}
)
