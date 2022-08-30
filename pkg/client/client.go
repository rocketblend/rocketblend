package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/library"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/scribble"
)

type (
	InstallService interface {
		FindAll() ([]*install.Install, error)
		FindByID(id string) (*install.Install, error)
		Create(i *install.Install) error
		Remove(id string) error
	}

	LibraryService interface {
		FindBuildByPath(path string) (*library.Build, error)
		FindPackageByPath(path string) (*library.Build, error)
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
		Platform        runtime.Platform
	}

	Client struct {
		install    InstallService
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

	platform := runtime.DetectPlatform()
	if platform == runtime.Undefined {
		return nil, fmt.Errorf("cannot detect platform")
	}

	appDir := filepath.Join(home, fmt.Sprintf(".%s", "rocketblend"))
	conf := Config{
		InstallationDir: filepath.Join(appDir, "installations"),
		DBDir:           filepath.Join(appDir, "data"),
		Platform:        platform,
	}

	return &conf, nil
}

func (c *Client) FindInstall(repo string) (*install.Install, error) {
	id := c.encoder.Hash(repo)
	install, err := c.install.FindByID(id)
	if err != nil {
		return nil, err
	}

	return install, nil
}

func (c *Client) AddInstall(install *install.Install) error {
	err := c.install.Create(install)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) InstallBuild(repo string) error {
	// Check if install already exists
	inst, err := c.FindInstall(repo)
	if err != nil {
		return err
	}

	if inst != nil {
		return fmt.Errorf("already installed")
	}

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
	source := build.GetSourceForPlatform(c.conf.Platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", c.conf.Platform)
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
		Build: repo,
		Path:  outPath,
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

// func (c *Client) FindOrFetchBuild(build string) (*library.Build, error) {
// 	var b *library.Build
// 	var err error

// 	install, _ := c.FindInstall(build)
// 	if install != nil {
// 		b, err = c.library.FindBuildByPath(install.Build)
// 	} else {
// 		b, err = c.library.FetchBuild(build)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return b, nil
// }

func (c *Client) FindBuildByPath(build string) (*library.Build, error) {
	return c.library.FindBuildByPath(build)
}

func (c *Client) FetchPackage(source string) (*library.Package, error) {
	pack, err := c.library.FetchPackage(source)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (c *Client) Platform() runtime.Platform {
	return c.conf.Platform
}
