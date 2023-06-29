package rocketpack

import (
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Build struct {
		Args    string                       `json:"args,omitempty"`
		Version *semver.Version              `json:"version,omitempty"`
		Sources map[runtime.Platform]*Source `json:"sources,omitempty"`
		Addons  []reference.Reference        `json:"addons,omitempty"`
	}
)
