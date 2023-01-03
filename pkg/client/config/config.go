package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/spf13/viper"
)

type (
	Defaults struct {
		Build string `mapstructure:"build"`
	}

	Directories struct {
		Installations string `mapstructure:"installations"`
		Resources     string `mapstructure:"resources"`
	}

	Features struct {
		Addons bool `mapstructure:"addons"`
	}

	Config struct {
		Debug       bool             `mapstructure:"debug"`
		Platform    runtime.Platform `mapstructure:"platform"`
		Defaults    *Defaults        `mapstructure:"defaults"`
		Directories *Directories     `mapstructure:"directories"`
		Features    *Features        `mapstructure:"features"`
	}
)

func PlatformHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		// Check that the data is string
		if f.Kind() != reflect.String {
			return data, nil
		}

		// Check that the target type is our custom type
		if t != reflect.TypeOf(runtime.Platform(0)) {
			return data, nil
		}

		// Return the parsed value
		p := runtime.Platform(0)
		return p.FromString(data.(string)), nil
	}
}

func Load() (config *Config, err error) {
	v := viper.New()

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find user home directory: %v", err)
	}

	platform := runtime.DetectPlatform()
	if platform == runtime.Undefined {
		return nil, fmt.Errorf("cannot detect platform")
	}

	appDir := filepath.Join(home, "rocketblend")

	v.SetDefault("debug", false)
	v.SetDefault("platform", platform)
	v.SetDefault("defaults.build", "")
	v.SetDefault("directories.installations", filepath.Join(appDir, "installations"))
	v.SetDefault("directories.resources", filepath.Join(appDir, "resources"))
	v.SetDefault("features.addons", false)

	v.SetConfigName("settings") // Set the name of the configuration file
	v.AddConfigPath(appDir)     // Look for the configuration file at the home directory
	v.SetConfigType("yml")      // Set the config type to YAML

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = v.Unmarshal(&config, viper.DecodeHook(PlatformHookFunc()))

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, err
}
