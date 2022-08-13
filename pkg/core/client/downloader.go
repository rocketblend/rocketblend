package client

import "github.com/rocketblend/rocketblend/pkg/core/downloader"

func NewDownloaderService(dir string) *downloader.Service {
	conf := downloader.NewConfig(dir, true)
	srv := downloader.NewService(conf)

	return srv
}
