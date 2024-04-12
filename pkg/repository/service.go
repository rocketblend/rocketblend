package repository

import (
	"errors"
	"os"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
)

type (
	Options struct {
		Logger    logger.Logger
		Validator types.Validator

		Platform         runtime.Platform
		StoragePath      string
		InstallationPath string

		Downloader types.Downloader
		Extractor  types.Extractor
	}

	Option func(*Options)

	repository struct {
		logger           logger.Logger
		validator        types.Validator
		downloader       types.Downloader
		extractor        types.Extractor
		platform         runtime.Platform
		storagePath      string
		installationPath string
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func WithStoragePath(storagePath string) Option {
	return func(o *Options) {
		o.StoragePath = storagePath
	}
}

func WithInstallationPath(installationPath string) Option {
	return func(o *Options) {
		o.InstallationPath = installationPath
	}
}

func NewService(opts ...Option) (*repository, error) {
	options := &Options{
		Logger:    logger.NoOp(),
		Validator: validator.New(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Downloader == nil {
		return nil, errors.New("downloader is nil")
	}

	if options.Extractor == nil {
		return nil, errors.New("extractor is nil")
	}

	if options.StoragePath == "" {
		return nil, errors.New("storage path is empty")
	}

	if options.InstallationPath == "" {
		return nil, errors.New("installation path is empty")
	}

	if err := os.MkdirAll(options.StoragePath, 0755); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(options.InstallationPath, 0755); err != nil {
		return nil, err
	}

	options.Logger.Debug("initializing rocketpack service", map[string]interface{}{
		"storagePath": options.StoragePath,
	})

	return &repository{
		logger:           options.Logger,
		validator:        options.Validator,
		downloader:       options.Downloader,
		extractor:        options.Extractor,
		platform:         options.Platform,
		storagePath:      options.StoragePath,
		installationPath: options.InstallationPath,
	}, nil
}
