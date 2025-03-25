package configurator

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rocketblend/rocketblend/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/runtime"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
	"github.com/spf13/viper"
)

type (
	Options struct {
		Logger    types.Logger
		Validator types.Validator
		Path      string
		Name      string
		Extension string
	}

	Option func(*Options)

	Configurator struct {
		logger    types.Logger
		validator types.Validator
		viper     *viper.Viper
		path      string
		extension string
		name      string
	}
)

func WithLogger(logger types.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func WithLocation(path string) Option {
	return func(o *Options) {
		o.Path = path
	}
}

func WithApplication(name string, extenstion string) Option {
	return func(o *Options) {
		o.Name = name
		o.Extension = extenstion
	}
}

func New(opts ...Option) (*Configurator, error) {
	options := &Options{
		Logger:    logger.NoOp(),
		Validator: validator.New(),
		Extension: "json",
		Name:      "config",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Path == "" {
		return nil, errors.New("path is required")
	}

	viper, err := load(options.Path, options.Name, options.Extension)
	if err != nil {
		return nil, err
	}

	return &Configurator{
		logger:    options.Logger,
		validator: options.Validator,
		path:      options.Path,
		extension: options.Extension,
		name:      options.Name,
		viper:     viper,
	}, nil
}

func (c *Configurator) Get() (*types.Config, error) {
	var config types.Config
	err := c.viper.Unmarshal(&config, viper.DecodeHook(platformHookFunc()))
	if err != nil {
		return nil, err
	}

	if err := c.validator.Validate(config); err != nil {
		return nil, err
	}

	return &config, err
}

func (c *Configurator) GetAllValues() map[string]interface{} {
	return c.viper.AllSettings()
}

func (c *Configurator) GetValueByString(key string) string {
	return fmt.Sprint(c.viper.Get(key))
}

func (c *Configurator) SetValueByString(key string, value string) error {
	if key == "" {
		return errors.New("key is required")
	}

	c.viper.Set(key, value)

	_, err := c.Get()
	if err != nil {
		return err
	}

	c.viper.WriteConfig()

	return nil
}

func (c *Configurator) Save(config *types.Config) error {
	if err := c.validator.Validate(config); err != nil {
		return err
	}

	var m map[string]interface{}
	mapstructure.Decode(config, &m)

	c.viper.MergeConfigMap(m)
	return c.viper.WriteConfig()
}

func (c *Configurator) Path() string {
	return fmt.Sprintf("%s.%s", filepath.Join(c.path, c.name), c.extension)
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

func load(path string, name string, extension string) (*viper.Viper, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	platform := runtime.DetectPlatform()
	if platform == runtime.Undefined {
		return nil, fmt.Errorf("cannot detect platform")
	}

	v.SetDefault("platform", platform.String())
	v.SetDefault("defaultBuild", types.DefaultBuild)
	v.SetDefault("logLevel", "info")
	v.SetDefault("installationsPath", filepath.Join(path, "installations"))
	v.SetDefault("packagesPath", filepath.Join(path, "packages"))
	v.SetDefault("aliases", types.DefaultAliases)

	v.SetConfigName(name)      // Set the name of the configuration file
	v.AddConfigPath(path)      // Look for the configuration file at the home directory
	v.SetConfigType(extension) // Set the config type to JSON

	v.SafeWriteConfig()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}
