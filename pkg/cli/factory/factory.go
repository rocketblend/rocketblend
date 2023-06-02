package factory

import (
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/cli/config"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/blendfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/installation"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketpack"
)

type (
	Factory interface {
		GetConfigService() (*config.Service, error)
		GetRocketPackService() (rocketpack.Service, error)
		GetInstallationService() (installation.Service, error)
		GetBlendFileService() (blendfile.Service, error)
		CreateRocketBlendService(blendConfig *blendconfig.BlendConfig) (rocketblend.Driver, error)
	}

	factory struct {
		logger        logger.Logger
		configService *config.Service

		rocketPackService   rocketpack.Service
		installationService installation.Service
		blendFileService    blendfile.Service
	}
)

func New() (Factory, error) {
	configService, err := config.New()
	if err != nil {
		return nil, err
	}

	config, err := configService.Get()
	if err != nil {
		return nil, err
	}

	logger := logger.New(
		logger.WithLogLevel(config.LogLevel),
		logger.WithPretty(),
	)

	rocketPackService, err := rocketpack.NewService(
		rocketpack.WithLogger(logger),
		rocketpack.WithStoragePath(config.PackagesPath),
	)
	if err != nil {
		return nil, err
	}

	installationService, err := installation.NewService(
		installation.WithLogger(logger),
		installation.WithStoragePath(config.InstallationsPath),
	)
	if err != nil {
		return nil, err
	}

	blendFileService, err := blendfile.NewService(
		blendfile.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	return &factory{
		logger:              logger,
		configService:       configService,
		rocketPackService:   rocketPackService,
		installationService: installationService,
		blendFileService:    blendFileService,
	}, nil
}

func (f *factory) GetConfigService() (*config.Service, error) {
	return config.New()
}

func (f *factory) GetRocketPackService() (rocketpack.Service, error) {
	return f.rocketPackService, nil
}

func (f *factory) GetInstallationService() (installation.Service, error) {
	return f.installationService, nil
}

func (f *factory) GetBlendFileService() (blendfile.Service, error) {
	return f.blendFileService, nil
}

func (f *factory) CreateRocketBlendService(blendConfig *blendconfig.BlendConfig) (rocketblend.Driver, error) {
	return rocketblend.New(
		rocketblend.WithLogger(f.logger),
		rocketblend.WithRocketPackService(f.rocketPackService),
		rocketblend.WithInstallationService(f.installationService),
		rocketblend.WithBlendFileService(f.blendFileService),
		rocketblend.WithBlendConfig(blendConfig),
	)
}
