package types

import (
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/runtime"
)

var (
	DefaultBuild = "github.com/rocketblend/official-library/packages/v0/builds/blender/4.2.2"
)

type (
	Config struct {
		Platform          runtime.Platform    `mapstructure:"platform"`
		DefaultBuild      reference.Reference `mapstructure:"defaultBuild"`
		InstallationsPath string              `mapstructure:"installationsPath"`
		PackagesPath      string              `mapstructure:"packagesPath"`
	}

	Configurator interface {
		Get() (config *Config, err error)
		GetAllValues() map[string]interface{}
		GetValueByString(key string) string
		SetValueByString(key string, value string) error
		Save(config *Config) error
		Path() string
	}
)
