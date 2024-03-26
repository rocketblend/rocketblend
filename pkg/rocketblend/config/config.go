package config

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/spf13/viper"
)

const (
	FileName      = "settings"
	FileExtension = "json"
)

type (
	Service interface {
		Get() (config *Config, err error)
		GetAllValues() map[string]interface{}
		GetValueByString(key string) string
		SetValueByString(key string, value string) error
		Save(config *Config) error
		Path() string
	}

	service struct {
		viper     *viper.Viper
		validator *validator.Validate
		rootPath  string
	}
)

func New(rootPath string) (Service, error) {
	v, err := load(rootPath)
	if err != nil {
		return nil, err
	}

	return &service{
		viper:     v,
		validator: validator.New(),
		rootPath:  rootPath,
	}, nil
}

func (srv *service) Get() (config *Config, err error) {
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

func (srv *service) GetAllValues() map[string]interface{} {
	return srv.viper.AllSettings()
}

func (srv *service) GetValueByString(key string) string {
	return fmt.Sprint(srv.viper.Get(key))
}

func (srv *service) SetValueByString(key string, value string) error {
	srv.viper.Set(key, value)

	_, err := srv.Get()
	if err != nil {
		return err
	}

	srv.viper.WriteConfig()

	return nil
}

func (srv *service) Save(config *Config) error {
	err := srv.validate(config)
	if err != nil {
		return err
	}

	var m map[string]interface{}
	mapstructure.Decode(config, &m)

	srv.viper.MergeConfigMap(m)

	return srv.viper.WriteConfig()
}

func (srv *service) Path() string {
	return fmt.Sprintf("%s.%s", filepath.Join(srv.rootPath, FileName), FileExtension)
}

func (srv *service) validate(config *Config) error {
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
		return runtime.PlatformFromString(data.(string)), nil
	}
}

func load(rootPath string) (*viper.Viper, error) {
	v := viper.New()

	platform := runtime.DetectPlatform()
	if platform == runtime.Undefined {
		return nil, fmt.Errorf("cannot detect platform")
	}

	v.SetDefault("platform", platform.String())
	v.SetDefault("defaultBuild", DefaultBuild)
	v.SetDefault("logLevel", "info")
	v.SetDefault("installationsPath", filepath.Join(rootPath, "installations"))
	v.SetDefault("packagesPath", filepath.Join(rootPath, "packages"))
	v.SetDefault("features.addons", false)

	v.SetConfigName(FileName)      // Set the name of the configuration file
	v.AddConfigPath(rootPath)      // Look for the configuration file at the home directory
	v.SetConfigType(FileExtension) // Set the config type to JSON

	v.SafeWriteConfig()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}
