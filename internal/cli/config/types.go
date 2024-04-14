package config

import (
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
)

var (
	DefaultBuild = "github.com/rocketblend/official-library/packages/blender/builds/stable/3.3.8"
)

type (
	Features struct {
		Addons bool `mapstructure:"addons"`
	}

	Config struct {
		Platform          runtime.Platform    `mapstructure:"platform"`
		DefaultBuild      reference.Reference `mapstructure:"defaultBuild"`
		LogLevel          string              `mapstructure:"logLevel"`
		InstallationsPath string              `mapstructure:"installationsPath"`
		PackagesPath      string              `mapstructure:"packagesPath"`
		Features          Features            `mapstructure:"features"`
	}
)
