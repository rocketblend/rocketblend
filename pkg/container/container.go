package container

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/blender"
	"github.com/rocketblend/rocketblend/pkg/configurator"
	"github.com/rocketblend/rocketblend/pkg/repository"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
)

type (
	holder[T any] struct {
		instance *T
		once     sync.Once
	}

	Options struct {
		Logger     logger.Logger
		Validator  types.Validator
		Downloader types.Downloader
		Extractor  types.Extractor

		ApplicationName    string
		ApplicationVersion string
	}

	Option func(*Options)

	Container struct {
		logger     types.Logger
		validator  types.Validator
		downloader types.Downloader
		extractor  types.Extractor

		configuratorHolder *holder[configurator.Configurator]
		repositoryHolder   *holder[repository.Repository]
		blenderHolder      *holder[blender.Blender]

		applicationDir string
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

func WithApplication(name string, version string) Option {
	return func(o *Options) {
		o.ApplicationName = name
		o.ApplicationVersion = version
	}
}

func New(opts ...Option) (*Container, error) {
	options := &Options{
		Logger:    logger.NoOp(),
		Validator: validator.New(),
	}

	for _, opt := range opts {
		opt(options)
	}

	applicationDir, err := setupApplicationDir(options.ApplicationName, options.ApplicationVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to setup application directory: %w", err)
	}

	return &Container{
		logger:             options.Logger,
		validator:          options.Validator,
		applicationDir:     applicationDir,
		configuratorHolder: &holder[configurator.Configurator]{},
		repositoryHolder:   &holder[repository.Repository]{},
		blenderHolder:      &holder[blender.Blender]{},
	}, nil
}

func (f *Container) GetLogger() (types.Logger, error) {
	return f.logger, nil
}

func (f *Container) GetValidator() (types.Validator, error) {
	return f.validator, nil
}

func (f *Container) GetDownloader() (types.Downloader, error) {
	return f.downloader, nil
}

func (f *Container) GetExtractor() (types.Extractor, error) {
	return f.extractor, nil
}

func (f *Container) GetConfigurator() (types.Configurator, error) {
	return f.getConfigurator()
}

func (f *Container) GetRepository() (types.Repository, error) {
	return f.getRepository()
}

func (f *Container) GetBlender() (types.Blender, error) {
	return f.getBlender()
}

func (f *Container) getConfigurator() (*configurator.Configurator, error) {
	var err error
	f.configuratorHolder.once.Do(func() {
		f.configuratorHolder.instance, err = configurator.New(
			configurator.WithLogger(f.logger),
			configurator.WithValidator(f.validator),
			configurator.WithLocation(f.applicationDir),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create configurator: %w", err)
	}

	return f.configuratorHolder.instance, nil
}

func (f *Container) getRepository() (*repository.Repository, error) {
	var err error
	f.repositoryHolder.once.Do(func() {
		configurator, errConfig := f.getConfigurator()
		if errConfig != nil {
			err = errConfig
			return
		}

		config, errConfig := configurator.Get()
		if errConfig != nil {
			err = errConfig
			return
		}

		f.repositoryHolder.instance, err = repository.New(
			repository.WithLogger(f.logger),
			repository.WithValidator(f.validator),
			repository.WithDownloader(f.downloader),
			repository.WithExtractor(f.extractor),
			repository.WithPackagePath(config.PackagesPath),
			repository.WithInstallationPath(config.InstallationsPath),
			repository.WithPlatform(config.Platform),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create repository: %w", err)
	}

	return f.repositoryHolder.instance, nil
}

func (f *Container) getBlender() (*blender.Blender, error) {
	var err error
	f.blenderHolder.once.Do(func() {
		f.blenderHolder.instance, err = blender.New(
			blender.WithLogger(f.logger),
			blender.WithValidator(f.validator),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create blender: %w", err)
	}

	return f.blenderHolder.instance, nil
}

func setupApplicationDir(name string, version string) (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot find config directory: %v", err)
	}

	appDir := filepath.Join(userConfigDir, name)
	if version == "dev" {
		appDir = filepath.Join(appDir, "dev")
	}

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create app directory: %w", err)
	}

	return appDir, nil
}
