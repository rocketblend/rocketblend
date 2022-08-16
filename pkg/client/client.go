package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/build"
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/remote"
	"github.com/rocketblend/scribble"
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
		FetchAll(req build.FetchRequest) ([]*build.Build, error)
		Find(remotes []*remote.Remote, hash string) (*build.Build, error)
	}

	DownloadService interface {
		Download(url string) (string, error)
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
		conf       Config
	}
)

func NewClient(conf Config) (*Client, error) {
	db, err := scribble.New(conf.DBDir, nil)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(conf.InstallationDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create installation directory: %w", err)
	}

	client := &Client{
		install:    NewInstallService(db),
		remote:     NewRemoteService(db),
		build:      NewBuildService(),
		downloader: NewDownloaderService(conf.InstallationDir),
		archiver:   NewArchiverService(true),
		conf:       conf,
	}

	return client, nil
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find user home directory: %v", err)
	}

	appDir := filepath.Join(home, fmt.Sprintf(".%s", "rocketblend"))
	conf := Config{
		InstallationDir: filepath.Join(appDir, "installations"),
		DBDir:           filepath.Join(appDir, "data"),
	}

	return &conf, nil
}

func (c *Client) FindInstall(hash string) (*install.Install, error) {
	ints, err := c.install.FindByHash(hash)
	if err != nil {
		return nil, err
	}

	return ints, nil
}

func (c *Client) AddInstall(install *install.Install) error {
	err := c.install.Create(install)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) InstallBuild(hash string) error {
	// inst, err := c.install.FindByHash(hash)
	// if err != nil {
	// 	return err
	// }

	// if inst != nil {
	// 	return fmt.Errorf("already installed")
	// }

	remotes, err := c.remote.FindAll()
	if err != nil {
		return err
	}

	build, err := c.build.Find(remotes, hash)
	if err != nil {
		return err
	}

	if build == nil {
		return fmt.Errorf("invalid build")
	}

	dir, err := c.downloader.Download(build.DownloadUrl)
	if err != nil {
		return err
	}

	if err := c.archiver.Extract(dir); err != nil {
		return err
	}

	err = c.AddInstall(&install.Install{
		Hash:    hash,
		Name:    build.Name,
		Version: build.Version,
		Path:    filepath.Join(c.conf.InstallationDir, build.Name),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveInstall(hash string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) FindAllInstalls() ([]*install.Install, error) {
	ints, err := c.install.FindAll(install.FindRequest{})
	if err != nil {
		return nil, err
	}

	return ints, nil
}

func (c *Client) FetchRemoteBuilds(platform string) ([]*build.Build, error) {
	re, err := c.remote.FindAll()
	if err != nil {
		return nil, err
	}

	req := build.FetchRequest{
		Remotes:  re,
		Platform: platform,
	}

	builds, err := c.build.FetchAll(req)
	if err != nil {
		return nil, err
	}

	return builds, nil
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
