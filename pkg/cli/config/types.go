package config

import (
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
)

var (
	DefaultBuild = "github.com/rocketblend/official-library/packages/blender/builds/stable/3.4.1"
)

type (
	Config struct {
		Platform          runtime.Platform `mapstructure:"platform"`
		DefaultBuild      string           `mapstructure:"defaultBuild"`
		LogLevel          string           `mapstructure:"logLevel"`
		InstallationsPath string           `mapstructure:"installationsPath"`
		PackagesPath      string           `mapstructure:"packagesPath"`
	}
)
