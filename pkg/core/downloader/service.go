package downloader

import (
	"fmt"
	"path/filepath"

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

func (s *Service) Download(url string) (string, error) {
	name := filepath.Base(url)
	path := filepath.Join(s.conf.DownloadDir, name)

	err := downloader.DownloadFile(path, url)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %s", url, err)
	}

	return path, nil
}
