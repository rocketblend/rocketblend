package factory

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile"
	"github.com/rocketblend/rocketblend/pkg/driver/installation"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/build"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/config"
)

type (
	Factory interface {
		GetLogger() (logger.Logger, error)
		SetLogger(logger.Logger) error
		GetConfigService() (config.Service, error)
		GetRocketPackService() (rocketpack.Service, error)
		GetInstallationService() (installation.Service, error)
		GetBlendFileService() (blendfile.Service, error)
	}

	factory struct {
		appDir string

		loggerMutex sync.Mutex
		logger      logger.Logger

		configMutex   sync.Mutex
		configService config.Service

		rocketPackMutex   sync.Mutex
		rocketPackService rocketpack.Service

		installationMutex   sync.Mutex
		installationService installation.Service

		blendFileMutex   sync.Mutex
		blendFileService blendfile.Service
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

	return &factory{
		appDir: appDir,
	}, nil
}

func (f *factory) GetConfigService() (config.Service, error) {
	f.configMutex.Lock()
	defer f.configMutex.Unlock()

	if f.configService == nil {
		service, err := config.New(f.appDir)
		if err != nil {
			return nil, err
		}

		f.configService = service
	}

	return f.configService, nil
}

func (f *factory) GetLogger() (logger.Logger, error) {
	f.loggerMutex.Lock()
	defer f.loggerMutex.Unlock()
	if f.logger == nil {
		configService, err := f.GetConfigService()
		if err != nil {
			return nil, err
		}

		config, err := configService.Get()
		if err != nil {
			return nil, err
		}

		f.logger = f.getLogger(config.LogLevel)
	}

	return f.logger, nil
}

func (f *factory) SetLogger(logger logger.Logger) error {
	f.loggerMutex.Lock()
	defer f.loggerMutex.Unlock()

	f.logger = logger
	return nil
}

func (f *factory) GetRocketPackService() (rocketpack.Service, error) {
	f.rocketPackMutex.Lock()
	defer f.rocketPackMutex.Unlock()
	if f.rocketPackService == nil {
		logger, err := f.GetLogger()
		if err != nil {
			return nil, err
		}

		configService, err := f.GetConfigService()
		if err != nil {
			return nil, err
		}

		config, err := configService.Get()
		if err != nil {
			return nil, err
		}

		service, err := rocketpack.NewService(
			rocketpack.WithLogger(logger),
			rocketpack.WithStoragePath(config.PackagesPath),
		)
		if err != nil {
			return nil, err
		}

		f.rocketPackService = service
	}
	return f.rocketPackService, nil
}

func (f *factory) GetInstallationService() (installation.Service, error) {
	f.installationMutex.Lock()
	defer f.installationMutex.Unlock()
	if f.installationService == nil {
		logger, err := f.GetLogger()
		if err != nil {
			return nil, err
		}

		configService, err := f.GetConfigService()
		if err != nil {
			return nil, err
		}

		config, err := configService.Get()
		if err != nil {
			return nil, err
		}

		service, err := installation.NewService(
			installation.WithLogger(logger),
			installation.WithStoragePath(config.InstallationsPath),
			installation.WithPlatform(config.Platform),
		)
		if err != nil {
			return nil, err
		}

		f.installationService = service
	}

	return f.installationService, nil
}

func (f *factory) GetBlendFileService() (blendfile.Service, error) {
	f.blendFileMutex.Lock()
	defer f.blendFileMutex.Unlock()
	if f.blendFileService == nil {
		logger, err := f.GetLogger()
		if err != nil {
			return nil, err
		}

		configService, err := f.GetConfigService()
		if err != nil {
			return nil, err
		}

		config, err := configService.Get()
		if err != nil {
			return nil, err
		}

		service, err := blendfile.NewService(
			blendfile.WithLogger(logger),
			blendfile.WithAddonsEnabled(config.Features.Addons),
		)
		if err != nil {
			return nil, err
		}

		f.blendFileService = service
	}

	return f.blendFileService, nil
}

func (f *factory) getLogger(logLevel string) logger.Logger {
	return logger.New(
		logger.WithLogLevel(logLevel),
		logger.WithPretty(),
	)
}
