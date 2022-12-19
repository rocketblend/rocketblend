package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/rocketblend/pkg/core/archiver"
	"github.com/rocketblend/rocketblend/pkg/core/downloader"
	"github.com/rocketblend/rocketblend/pkg/core/encoder"
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/library"
	"github.com/rocketblend/rocketblend/pkg/core/preference"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/scribble"
)

func NewAddonService(driver *scribble.Driver) *addon.Service {
	repo := addon.NewScribbleRepository(driver)
	srv := addon.NewService(repo)

	return srv
}

func NewInstallService(driver *scribble.Driver) *install.Service {
	repo := install.NewScribbleRepository(driver)
	conf := install.LoadConfig()
	srv := install.NewService(conf, repo)

	return srv
}

func NewPreferenceService(driver *scribble.Driver) *preference.Service {
	repo := preference.NewScribbleRepository(driver)
	srv := preference.NewService(repo)

	return srv
}

func NewLibraryService() *library.Service {
	conf := library.NewClientConfig()
	http := library.NewClient(conf)
	repo := library.NewRepository()
	srv := library.NewService(http, repo)

	return srv
}

func NewArchiverService(delete bool) *archiver.Service {
	conf := archiver.NewConfig(delete)
	srv := archiver.NewService(conf)

	return srv
}

func NewDownloaderService(dir string) *downloader.Service {
	conf := downloader.NewConfig(dir)
	srv := downloader.NewService(conf)

	return srv
}

func NewResourceService(dir string) *resource.Service {
	srv := resource.NewService(dir)
	return srv
}

func NewEncoderService() *encoder.Service {
	srv := encoder.NewService()

	return srv
}
