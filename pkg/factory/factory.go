package factory

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
	Options struct {
		Logger     logger.Logger
		Validator  types.Validator
		Downloader types.Downloader
		Extractor  types.Extractor

		ApplicationName    string
		ApplicationVersion string
	}

	Option func(*Options)

	factory struct {
		logger     types.Logger
		validator  types.Validator
		downloader types.Downloader
		extractor  types.Extractor

		configurator *configurator.Configurator
		repository   *repository.Repository
		blender      *blender.Blender

		rwMutex sync.RWMutex

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

func New(opts ...Option) (*factory, error) {
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

	return &factory{
		logger:         options.Logger,
		validator:      options.Validator,
		applicationDir: applicationDir,
	}, nil
}

func (f *factory) GetLogger() (types.Logger, error) {
	return f.logger, nil
}

func (f *factory) GetValidator() (types.Validator, error) {
	return f.validator, nil
}

func (f *factory) GetDownloader() (types.Downloader, error) {
	return f.downloader, nil
}

func (f *factory) GetExtractor() (types.Extractor, error) {
	return f.extractor, nil
}

func (f *factory) GetConfigurator() (types.Configurator, error) {
	return initService(&f.rwMutex, &f.configurator, f.getConfigurator)
}

func (f *factory) GetRepository() (types.Repository, error) {
	return initService(&f.rwMutex, &f.repository, f.getRepository)
}

func (f *factory) GetBlender() (types.Blender, error) {
	return initService(&f.rwMutex, &f.blender, f.getBlender)
}

func (f *factory) getRepository() (*repository.Repository, error) {
	configurator, err := f.getConfigurator()
	if err != nil {
		return nil, fmt.Errorf("failed to get configurator: %w", err)
	}

	config, err := configurator.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	repository, err := repository.New(
		repository.WithLogger(f.logger),
		repository.WithValidator(f.validator),
		repository.WithDownloader(f.downloader),
		repository.WithExtractor(f.extractor),
		repository.WithPackagePath(config.PackagesPath),
		repository.WithInstallationPath(config.InstallationsPath),
		repository.WithPlatform(config.Platform),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return repository, nil
}

func (f *factory) getBlender() (*blender.Blender, error) {
	blender, err := blender.New(
		blender.WithLogger(f.logger),
		blender.WithValidator(f.validator),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create blender: %w", err)
	}

	return blender, nil
}

func (f *factory) getConfigurator() (*configurator.Configurator, error) {
	configurator, err := configurator.New(
		configurator.WithLogger(f.logger),
		configurator.WithValidator(f.validator),
		configurator.WithLocation(f.applicationDir),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create configurator: %w", err)
	}

	return configurator, nil
}

func initService[T any](rwMutex *sync.RWMutex, service **T, initFunc func() (*T, error)) (*T, error) {
	rwMutex.RLock()
	if *service != nil {
		defer rwMutex.RUnlock()
		return *service, nil
	}
	rwMutex.RUnlock()

	rwMutex.Lock()
	defer rwMutex.Unlock()
	if *service == nil {
		var err error
		*service, err = initFunc()
		if err != nil {
			return nil, err
		}
	}
	return *service, nil
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
