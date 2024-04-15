package container

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/blender"
	"github.com/rocketblend/rocketblend/pkg/configurator"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/extractor"
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
		Logger    logger.Logger
		Validator types.Validator

		ApplicationName string

		Development bool
	}

	Option func(*Options)

	Container struct {
		logger    types.Logger
		validator types.Validator

		configuratorHolder *holder[configurator.Configurator]
		downloaderHolder   *holder[downloader.Downloader]
		extractorHolder    *holder[extractor.Extractor]
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

func WithApplicationName(name string) Option {
	return func(o *Options) {
		o.ApplicationName = name
	}
}

func WithDevelopmentMode(development bool) Option {
	return func(o *Options) {
		o.Development = development
	}
}

func New(opts ...Option) (*Container, error) {
	options := &Options{
		Logger:          logger.NoOp(),
		Validator:       validator.New(),
		ApplicationName: "rocketblend",
	}

	for _, opt := range opts {
		opt(options)
	}

	applicationDir, err := setupApplicationDir(options.ApplicationName, options.Development)
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
	return f.getDownloader()
}

func (f *Container) GetExtractor() (types.Extractor, error) {
	return f.getExtractor()
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

func (f *Container) getDownloader() (*downloader.Downloader, error) {
	var err error
	f.downloaderHolder.once.Do(func() {
		f.downloaderHolder.instance = downloader.New(
			downloader.WithLogger(f.logger),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create downloader: %w", err)
	}

	return f.downloaderHolder.instance, nil
}

func (f *Container) getExtractor() (*extractor.Extractor, error) {
	var err error
	f.extractorHolder.once.Do(func() {
		f.extractorHolder.instance = extractor.New(
			extractor.WithLogger(f.logger),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create extractor: %w", err)
	}

	return f.extractorHolder.instance, nil
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

		downloader, errDownloader := f.getDownloader()
		if err != nil {
			err = errDownloader
			return
		}

		extractor, errExtractor := f.getExtractor()
		if err != nil {
			err = errExtractor
			return
		}

		f.repositoryHolder.instance, err = repository.New(
			repository.WithLogger(f.logger),
			repository.WithValidator(f.validator),
			repository.WithDownloader(downloader),
			repository.WithExtractor(extractor),
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

func setupApplicationDir(name string, development bool) (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot find config directory: %v", err)
	}

	appDir := filepath.Join(userConfigDir, name)
	if development {
		appDir = filepath.Join(appDir, "dev")
	}

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create app directory: %w", err)
	}

	return appDir, nil
}