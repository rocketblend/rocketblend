package repository

import (
	"errors"
	"os"

	"github.com/rocketblend/rocketblend/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/runtime"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
)

type (
	Options struct {
		Logger    types.Logger
		Validator types.Validator

		Platform         runtime.Platform
		PackagePath      string
		InstallationPath string

		Downloader types.Downloader
		Extractor  types.Extractor
	}

	Option func(*Options)

	Repository struct {
		logger           types.Logger
		validator        types.Validator
		downloader       types.Downloader
		extractor        types.Extractor
		platform         runtime.Platform
		packagePath      string
		installationPath string
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

func WithDownloader(downloader types.Downloader) Option {
	return func(o *Options) {
		o.Downloader = downloader
	}
}

func WithExtractor(extractor types.Extractor) Option {
	return func(o *Options) {
		o.Extractor = extractor
	}
}

func WithPackagePath(packagePath string) Option {
	return func(o *Options) {
		o.PackagePath = packagePath
	}
}

func WithInstallationPath(installationPath string) Option {
	return func(o *Options) {
		o.InstallationPath = installationPath
	}
}

func WithPlatform(platform runtime.Platform) Option {
	return func(o *Options) {
		o.Platform = platform
	}
}

func New(opts ...Option) (*Repository, error) {
	options := &Options{
		Logger:    logger.NoOp(),
		Validator: validator.New(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Validator == nil {
		return nil, errors.New("validator is nil")
	}

	if options.Downloader == nil {
		return nil, errors.New("downloader is nil")
	}

	if options.Extractor == nil {
		return nil, errors.New("extractor is nil")
	}

	if options.PackagePath == "" {
		return nil, errors.New("storage path is empty")
	}

	if options.InstallationPath == "" {
		return nil, errors.New("installation path is empty")
	}

	if err := os.MkdirAll(options.PackagePath, 0755); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(options.InstallationPath, 0755); err != nil {
		return nil, err
	}

	options.Logger.Debug("initializing repository", map[string]interface{}{
		"packagePath": options.PackagePath,
	})

	return &Repository{
		logger:           options.Logger,
		validator:        options.Validator,
		downloader:       options.Downloader,
		extractor:        options.Extractor,
		platform:         options.Platform,
		packagePath:      options.PackagePath,
		installationPath: options.InstallationPath,
	}, nil
}
