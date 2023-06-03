package factory

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/cli/build"
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
		CreateDriver(blendConfig *blendconfig.BlendConfig) (rocketblend.Driver, error)
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
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find config directory: %v", err)
	}

	appDir := filepath.Join(userConfigDir, build.AppName)
	if build.Version == "dev" {
		appDir = filepath.Join(appDir, "dev")
	}

	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create app directory: %w", err)
	}

	configService, err := config.New(appDir)
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
		installation.WithPlatform(config.Platform),
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
	// TODO: make this either return a new instance or a cached instance
	return f.configService, nil
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

func (f *factory) CreateDriver(blendConfig *blendconfig.BlendConfig) (rocketblend.Driver, error) {
	return rocketblend.New(
		rocketblend.WithLogger(f.logger),
		rocketblend.WithRocketPackService(f.rocketPackService),
		rocketblend.WithInstallationService(f.installationService),
		rocketblend.WithBlendFileService(f.blendFileService),
		rocketblend.WithBlendConfig(blendConfig),
	)
}
