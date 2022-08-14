package client

import (
	"fmt"
	"sort"

	"github.com/blang/semver/v4"
	"go.lsp.dev/uri"

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
		FetchAll(req build.FetchRequest) ([]*build.Build, error)
		Find(hash string) (*build.Build, error)
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

func LoadConfig() Config {
	// Viper config
	var conf Config
	return conf
}

func (c *Client) AddInstall(install *install.Install) error {
	err := c.install.Create(install)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) InstallBuild(hash string) error {
	inst, err := c.install.FindByHash(hash)
	if err != nil {
		return err
	}

	if inst != nil {
		return fmt.Errorf("already installed")
	}

	build, err := c.build.Find(hash)
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
		Path:    dir,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveInstall(hash string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) GetAvilableBuilds(platform string, tag string) ([]*Available, error) {
	re, err := c.remote.FindAll()
	if err != nil {
		return nil, err
	}

	req := build.FetchRequest{
		Remotes:  re,
		Platform: platform,
		Tag:      tag,
	}

	builds, err := c.build.FetchAll(req)
	if err != nil {
		return nil, err
	}

	insts, err := c.install.FindAll(install.FindRequest{})
	if err != nil {
		return nil, err
	}

	var available []*Available
	for _, inst := range insts {
		var t *Available
		t.Hash = inst.Hash
		t.Name = inst.Name
		t.Version, _ = semver.Parse(inst.Version)
		t.Uri = uri.File(inst.Path)
		available = append(available, t)
	}

	for _, b := range builds {
		isExisting := false
		for _, existing := range available {
			if b.Hash == existing.Hash {
				isExisting = true
				break
			}
		}

		if !isExisting {
			var t *Available
			t.Hash = b.Hash
			t.Name = b.Name
			t.Version, _ = semver.Parse(b.Version)
			t.Uri, _ = uri.Parse(b.DownloadUrl)
			available = append(available, t)
		}
	}

	sort.SliceStable(available, func(i, j int) bool {
		return available[i].Version.GT(available[j].Version)
	})

	return available, nil
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
