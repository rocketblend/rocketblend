package downloader

import (
	"github.com/rocketblend/rocketblend/pkg/archiver"
	"github.com/rocketblend/rocketblend/pkg/downloader"
)

type (
	Config struct {
		Unarchive   bool
		DownloadDir string
	}

	Service struct {
		conf Config
	}
)

func NewService(conf Config) *Service {
	srv := &Service{
		conf: conf,
	}

	return srv
}

func NewConfig(dir string, unarchive bool) Config {
	return Config{
		Unarchive:   unarchive,
		DownloadDir: dir,
	}
}

func (s *Service) Download(url string) error {
	err := downloader.DownloadFile(s.conf.DownloadDir, url)
	if err != nil {
		return err
	}

	if s.conf.Unarchive {
		err = archiver.Extract(s.conf.DownloadDir, true)
		if err != nil {
			return err
		}
	}

	return nil
}
