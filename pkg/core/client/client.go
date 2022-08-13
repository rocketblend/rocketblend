package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/scribble"
)

type (
	InstallService interface {
		FindAll(req install.FindRequest) ([]*install.Install, error)
		FindByHash(hash string) (*install.Install, error)
		Create(i *install.Install) error
		Remove(hash string) error
	}

	DownloadService interface {
		Download(url string) error
	}

	Config struct {
		DBDir           string
		InstallationDir string
	}

	Client struct {
		install    InstallService
		downloader DownloadService
	}
)

func NewClient(conf Config) (*Client, error) {
	db, err := scribble.New(conf.DBDir, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		install:    NewInstallService(db),
		downloader: NewDownloaderService(conf.InstallationDir),
	}

	return client, nil
}

func LoadConfig() Config {
	// Viper config
	var conf Config
	return conf
}

func (c *Client) InstallBuild(hash string) error {
	return nil
}
