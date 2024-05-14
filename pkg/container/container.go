package container

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/blender"
	"github.com/rocketblend/rocketblend/pkg/configurator"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/driver"
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
		Logger           logger.Logger
		Validator        types.Validator
		ProgressInterval time.Duration
		DownloadBuffer   int
		ApplicationName  string
		Development      bool
	}

	Option func(*Options)

	Container struct {
		logger         types.Logger
		validator      types.Validator
		applicationDir string

		downloadBuffer   int
		progressInterval time.Duration

		configuratorHolder *holder[configurator.Configurator]
		downloaderHolder   *holder[downloader.Downloader]
		extractorHolder    *holder[extractor.Extractor]
		repositoryHolder   *holder[repository.Repository]
		driverHolder       *holder[driver.Driver]
		blenderHolder      *holder[blender.Blender]
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

func WithProgressInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.ProgressInterval = interval
	}
}

func WithDownloadBuffer(buffer int) Option {
	return func(o *Options) {
		o.DownloadBuffer = buffer
	}
}

func New(opts ...Option) (*Container, error) {
	options := &Options{
		Logger:           logger.NoOp(),
		Validator:        validator.New(),
		ApplicationName:  types.ApplicationName,
		DownloadBuffer:   1 << 20,         // Default buffer size is 1MB
		ProgressInterval: 5 * time.Second, // Default progress interval is 5 seconds
	}

	for _, opt := range opts {
		opt(options)
	}

	applicationDir, err := setupApplicationDir(options.ApplicationName, options.Development)
	if err != nil {
		return nil, fmt.Errorf("failed to setup application directory: %w", err)
	}

	options.Logger.Debug("initializing container", map[string]interface{}{
		"path":        applicationDir,
		"development": options.Development,
	})

	return &Container{
		logger:             options.Logger,
		validator:          options.Validator,
		applicationDir:     applicationDir,
		downloadBuffer:     options.DownloadBuffer,
		progressInterval:   options.ProgressInterval,
		configuratorHolder: &holder[configurator.Configurator]{},
		downloaderHolder:   &holder[downloader.Downloader]{},
		extractorHolder:    &holder[extractor.Extractor]{},
		repositoryHolder:   &holder[repository.Repository]{},
		driverHolder:       &holder[driver.Driver]{},
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

func (f *Container) GetDriver() (types.Driver, error) {
	return f.getDriver()
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
		f.downloaderHolder.instance, err = downloader.New(
			downloader.WithLogger(f.logger),
			downloader.WithBufferSize(f.downloadBuffer),
			downloader.WithUpdateInterval(f.progressInterval),
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
		f.extractorHolder.instance, err = extractor.New(
			extractor.WithLogger(f.logger),
			extractor.WithCleanup(),
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

func (f *Container) getDriver() (*driver.Driver, error) {
	var err error
	f.driverHolder.once.Do(func() {
		repository, errRepository := f.getRepository()
		if errRepository != nil {
			err = errRepository
			return
		}

		f.driverHolder.instance, err = driver.New(
			driver.WithLogger(f.logger),
			driver.WithValidator(f.validator),
			driver.WithRepository(repository),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get/create driver: %w", err)
	}

	return f.driverHolder.instance, nil
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
