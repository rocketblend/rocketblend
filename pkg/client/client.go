package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/build"
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/library"
	"github.com/rocketblend/rocketblend/pkg/core/remote"
	"github.com/rocketblend/scribble"
)

type (
	InstallService interface {
		FindAll() ([]*install.Install, error)
		FindBySource(source string) (*install.Install, error)
		Create(i *install.Install) error
		Remove(source string) error
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

	LibraryService interface {
		FetchBuild(str string) (*library.Build, error)
		FetchPackage(str string) (*library.Package, error)
	}

	DownloadService interface {
		Download(url string, path string) error
	}

	ArchiverService interface {
		Extract(path string) error
	}

	EncoderService interface {
		Hash(str string) string
	}

	Config struct {
		DBDir           string
		InstallationDir string
	}

	Client struct {
		install    InstallService
		remote     RemoteService
		build      BuildService
		library    LibraryService
		downloader DownloadService
		archiver   ArchiverService
		encoder    EncoderService
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
		library:    NewLibraryService(),
		downloader: NewDownloaderService(conf.InstallationDir),
		archiver:   NewArchiverService(true),
		encoder:    NewEncoderService(),
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
		InstallationDir: filepath.Join(appDir, "installations-2"),
		DBDir:           filepath.Join(appDir, "data-2"),
	}

	return &conf, nil
}

func (c *Client) FindInstall(source string) (*install.Install, error) {
	ints, err := c.install.FindBySource(source)
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

func (c *Client) InstallBuild(repo string) error {
	// inst, err := c.install.FindByHash(hash)
	// if err != nil {
	// 	return err
	// }

	// if inst != nil {
	// 	return fmt.Errorf("already installed")
	// }

	// Fetch build from repo
	build, err := c.library.FetchBuild(repo)
	if err != nil {
		return err
	}

	if build == nil {
		return fmt.Errorf("invalid build")
	}

	// Output directory
	outPath := filepath.Join(c.conf.InstallationDir, repo)

	// Create output directories
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return err
	}

	// build info for current platform
	platform := "windows" // TODO: detect platform
	source := build.GetSourceForPlatform(platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", platform)
	}

	// Download URL
	downloadURL := source.URL

	// Download file path
	name := filepath.Base(downloadURL)
	filePath := filepath.Join(outPath, name)

	// Download file to file path
	err = c.downloader.Download(downloadURL, filePath)
	if err != nil {
		return err
	}

	// Extract the archived file
	if err := c.archiver.Extract(filePath); err != nil {
		return err
	}

	// Markshal build
	js, err := json.Marshal(build)
	if err != nil {
		return err
	}

	// Write out build config.json
	if err := os.WriteFile(filepath.Join(outPath, "config.json"), js, os.ModePerm); err != nil {
		return err
	}

	// Add install to database
	err = c.AddInstall(&install.Install{
		Id:    c.encoder.Hash(repo),
		Path:  outPath,
		Build: build,
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
	ints, err := c.install.FindAll()
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

func (c *Client) FetchBuild(source string) (*library.Build, error) {
	build, err := c.library.FetchBuild(source)
	if err != nil {
		return nil, err
	}

	return build, nil
}

func (c *Client) FetchPackage(source string) (*library.Package, error) {
	pack, err := c.library.FetchPackage(source)
	if err != nil {
		return nil, err
	}

	return pack, nil
}
