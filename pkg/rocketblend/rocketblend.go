package rocketblend

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/resource"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
)

const (
	Name                 = "rocketblend"
	BlenderFileExtension = ".blend"
)

type (
	Driver struct {
		logger                logger.Logger
		resource              resource.Service
		pack                  rocketpack.Service
		debug                 bool
		platform              runtime.Platform
		defaultBuild          string
		InstallationDirectory string
		addonsEnabled         bool
	}

	Options struct {
		Logger                logger.Logger
		Debug                 bool
		Platform              runtime.Platform
		InstallationDirectory string
		AddonsEnabled         bool
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithPlatform(platform runtime.Platform) Option {
	return func(o *Options) {
		o.Platform = platform
	}
}

func WithInstallationDirectory(dir string) Option {
	return func(o *Options) {
		o.InstallationDirectory = dir
	}
}

func WithDebug() Option {
	return func(o *Options) {
		o.Debug = true
	}
}

func WithAddonsEnabled() Option {
	return func(o *Options) {
		o.AddonsEnabled = true
	}
}

func New(opts ...Option) (*Driver, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	// if not installation directory is provided, use the default
	if options.InstallationDirectory == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("cannot find config directory: %v", err)
		}

		dir := filepath.Join(configDir, Name, "packages")

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create main directory: %w", err)
		}

		options.InstallationDirectory = dir
	}

	// if no platform is provided, detect it
	if options.Platform == runtime.Undefined {
		platform := runtime.DetectPlatform()
		if platform == runtime.Undefined {
			return nil, fmt.Errorf("cannot detect platform")
		}

		options.Platform = platform
	}

	jot, err := jot.New(jot.WithLogger(options.Logger), jot.WithStorageDir(options.InstallationDirectory))
	if err != nil {
		return nil, err
	}

	pack, err := rocketpack.NewService(rocketpack.WithLogger(options.Logger), rocketpack.WithStorage(jot))
	if err != nil {
		return nil, err
	}

	options.Logger.Debug("RocketBlend initialized", map[string]interface{}{
		"platform":              options.Platform.String(),
		"installationDirectory": options.InstallationDirectory,
		"addonsEnabled":         options.AddonsEnabled,
	})

	// create driver
	return &Driver{
		logger:                options.Logger,
		pack:                  pack,
		resource:              resource.NewService(),
		debug:                 options.Debug,
		InstallationDirectory: options.InstallationDirectory,
		platform:              options.Platform,
		addonsEnabled:         options.AddonsEnabled,
	}, nil
}
