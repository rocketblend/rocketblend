package factory

import (
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/cli/config"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"
)

type (
	Factory interface {
		CreateConfigService() (*config.Service, error)
		CreateRocketBlendService() (*rocketblend.Driver, error)
	}

	factory struct {
	}
)

func New() Factory {
	return &factory{}
}

func (f *factory) CreateConfigService() (*config.Service, error) {
	return config.New()
}

func (f *factory) CreateRocketBlendService() (*rocketblend.Driver, error) {
	cs, err := config.New()
	if err != nil {
		return nil, err
	}

	config, err := cs.Get()
	if err != nil {
		return nil, err
	}

	opts := []rocketblend.Option{
		rocketblend.WithInstallationDirectory(config.InstallDir),
		rocketblend.WithPlatform(config.Platform),
	}

	if config.Debug {
		opts = append(opts, rocketblend.WithLogger(logger.New(logger.WithPretty(), logger.WithLogLevel("debug"))))
	}

	if config.Features.Addons {
		opts = append(opts, rocketblend.WithAddonsEnabled())
	}

	return rocketblend.New(opts...)
}
