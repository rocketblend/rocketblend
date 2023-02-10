package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/rocketblend/rocketblend/pkg/core"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/spf13/viper"
)

type Service struct {
	viper     *viper.Viper
	validator *validator.Validate
}

func New() (*Service, error) {
	v, err := load()
	if err != nil {
		return nil, err
	}

	return &Service{
		viper:     v,
		validator: validator.New(),
	}, nil
}

func (srv *Service) Get() (config *Config, err error) {
	err = srv.viper.Unmarshal(&config, viper.DecodeHook(platformHookFunc()))
	if err != nil {
		return nil, err
	}

	err = srv.validate(config)
	if err != nil {
		return nil, err
	}

	return config, err
}

func (srv *Service) GetAllValues() map[string]interface{} {
	return srv.viper.AllSettings()
}

func (srv *Service) GetValueByString(key string) string {
	return fmt.Sprint(srv.viper.Get(key))
}

func (srv *Service) SetValueByString(key string, value string) error {
	srv.viper.Set(key, value)

	_, err := srv.Get()
	if err != nil {
		return err
	}

	srv.viper.WriteConfig()

	return nil
}

func (srv *Service) Save(config *Config) error {
	err := srv.validate(config)
	if err != nil {
		return err
	}

	var m map[string]interface{}
	mapstructure.Decode(config, &m)

	srv.viper.MergeConfigMap(m)

	return srv.viper.WriteConfig()
}

func (srv *Service) validate(config *Config) error {
	if err := srv.validator.Struct(config); err != nil {
		return err
	}

	return nil
}

func platformHookFunc() mapstructure.DecodeHookFuncType {
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

func load() (*viper.Viper, error) {
	v := viper.New()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find config directory: %v", err)
	}

	platform := runtime.DetectPlatform()
	if platform == runtime.Undefined {
		return nil, fmt.Errorf("cannot detect platform")
	}

	appDir := filepath.Join(configDir, core.Name)

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create main directory: %w", err)
	}

	v.SetDefault("debug", false)
	v.SetDefault("platform", platform.String())
	v.SetDefault("defaultBuild", core.DefaultBuild)
	v.SetDefault("features.addons", false)

	v.SetConfigName("settings") // Set the name of the configuration file
	v.AddConfigPath(appDir)     // Look for the configuration file at the home directory
	v.SetConfigType("json")     // Set the config type to JSON

	v.SafeWriteConfig()

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}
