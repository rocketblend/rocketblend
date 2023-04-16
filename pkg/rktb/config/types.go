package config

import (
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
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
