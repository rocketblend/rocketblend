package config

import (
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
)

var (
	DefaultBuild = "github.com/rocketblend/official-library/packages/blender/builds/stable/3.4.1"
)

type (
	Features struct {
		Addons bool `mapstructure:"addons"`
	}

	Config struct {
		Debug        bool             `mapstructure:"debug"`
		Platform     runtime.Platform `mapstructure:"platform"`
		InstallDir   string           `mapstructure:"installDir"`
		DefaultBuild string           `mapstructure:"defaultBuild"`
		Features     *Features        `mapstructure:"features"`
	}
)
