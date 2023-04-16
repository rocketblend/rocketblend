package rocketblend

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/resource"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
	"github.com/sirupsen/logrus"
)

func New(options *Options) (*Driver, error) {
	// create default options
	opts := Options{}

	// if options are passed in, use those
	if options != nil {
		opts = *options
	}

	// if no logger is provided, create a default
	if opts.Logger == nil {
		l := logrus.New()
		l.Level = logrus.InfoLevel
		opts.Logger = l
	}

	// if not installation directory is provided, use the default
	if opts.InstallationDirectory == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("cannot find config directory: %v", err)
		}

		dir := filepath.Join(configDir, Name, "packages")

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create main directory: %w", err)
		}

		opts.InstallationDirectory = dir
	}

	// if no platform is provided, detect it
	if opts.Platform == runtime.Undefined {
		platform := runtime.DetectPlatform()
		if platform == runtime.Undefined {
			return nil, fmt.Errorf("cannot detect platform")
		}

		opts.Platform = platform
	}

	jot, err := jot.New(opts.InstallationDirectory, nil)
	if err != nil {
		return nil, err
	}

	// create driver
	driver := Driver{
		log:                   opts.Logger,
		pack:                  rocketpack.NewService(jot, opts.Platform),
		resource:              resource.NewService(),
		debug:                 opts.Debug,
		InstallationDirectory: opts.InstallationDirectory,
		platform:              opts.Platform,
		addonsEnabled:         opts.AddonsEnabled,
	}

	return &driver, nil
}
