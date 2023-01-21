package core

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/core/config"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/rocketblend/pkg/core/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/jot"
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

	if opts.Config == nil {
		conf, err := config.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
		opts.Config = conf
	}

	// if no resource service is provided, create a default
	if opts.ResourceService == nil {
		srv := resource.NewService()
		opts.ResourceService = srv
	}

	// if no pack service is provided, create a default
	if opts.PackService == nil {
		jot, err := jot.New(opts.Config.Directories.Installations, nil)
		if err != nil {
			return nil, err
		}
		opts.PackService = rocketpack.NewService(jot, opts.Config.Platform)
	}

	// create driver
	driver := Driver{
		conf:     opts.Config,
		log:      opts.Logger,
		pack:     opts.PackService,
		resource: opts.ResourceService,
	}

	return &driver, nil
}
