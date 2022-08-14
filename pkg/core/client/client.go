package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/build"
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/remote"
	"github.com/rocketblend/rocketblend/pkg/scribble"
)

type (
	InstallService interface {
		FindAll(req install.FindRequest) ([]*install.Install, error)
		FindByHash(hash string) (*install.Install, error)
		Create(i *install.Install) error
		Remove(hash string) error
	}

	RemoteService interface {
		FindAll() ([]*remote.Remote, error)
		Add(remote *remote.Remote) error
		Remove(name string) error
	}

	BuildService interface {
		Fetch(req build.FetchRequest) ([]*build.Build, error)
	}

	DownloadService interface {
		Download(url string) error
	}

	ArchiverService interface {
		Extract(path string) error
	}

	Config struct {
		DBDir           string
		InstallationDir string
	}

	Client struct {
		install    InstallService
		remote     RemoteService
		build      BuildService
		downloader DownloadService
		archiver   ArchiverService
	}
)

func NewClient(conf Config) (*Client, error) {
	db, err := scribble.New(conf.DBDir, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		install:    NewInstallService(db),
		remote:     NewRemoteService(db),
		build:      NewBuildService(),
		downloader: NewDownloaderService(conf.InstallationDir),
		archiver:   NewArchiverService(true),
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

func (c *Client) RemoveInstall(hash string) error {
	return nil
}

func (c *Client) GetAvilableBuilds() ([]*build.Build, error) {
	return nil, nil
}

func (c *Client) GetRemotes() ([]*remote.Remote, error) {
	remotes, err := c.remote.FindAll()
	if err != nil {
		return nil, err
	}

	return remotes, nil
}

func (c *Client) AddRemote(name string, url string) error {
	remote := &remote.Remote{
		Name: name,
		URL:  url,
	}

	if err := c.remote.Add(remote); err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveRemote(name string) error {
	if err := c.remote.Remove(name); err != nil {
		return err
	}

	return nil
}
