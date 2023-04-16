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
		InstallPath  string           `mapstructure:"installPath"`
		DefaultBuild string           `mapstructure:"defaultBuild"`
		Features     *Features        `mapstructure:"features"`
	}
)
