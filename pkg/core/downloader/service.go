package downloader

import (
	"github.com/rocketblend/rocketblend/pkg/downloader"
)

type (
	Config struct {
		DownloadDir string
	}

	Service struct {
		conf Config
	}
)

func NewService(conf Config) *Service {
	return &Service{
		conf: conf,
	}
}

func NewConfig(dir string) Config {
	return Config{
		DownloadDir: dir,
	}
}

func (s *Service) Download(url string) error {
	err := downloader.DownloadFile(s.conf.DownloadDir, url)
	if err != nil {
		return err
	}

	return nil
}
