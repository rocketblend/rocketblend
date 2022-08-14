package client

import "github.com/rocketblend/rocketblend/pkg/core/downloader"

func NewDownloaderService(dir string) *downloader.Service {
	conf := downloader.NewConfig(dir)
	srv := downloader.NewService(conf)

	return srv
}
