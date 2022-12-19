package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/library"
	"github.com/rocketblend/rocketblend/pkg/core/preference"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
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

	AddonService interface {
		FindAll() ([]*addon.Addon, error)
		FindByID(id string) (*addon.Addon, error)
		Create(i *addon.Addon) error
		Remove(id string) error
	}

	PreferenceService interface {
		Find() (*preference.Settings, error)
		Create(i *preference.Settings) error
	}

	LibraryService interface {
		FindBuildByPath(path string) (*library.Build, error)
		FindPackageByPath(path string) (*library.Package, error)
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

	ResourceService interface {
		FindByName(name string) (*resource.Resource, error)
		SaveOut() error
	}

	Config struct {
		DBDir           string
		InstallationDir string
		ResourceDir     string
		Debug           bool
		Platform        runtime.Platform
	}

	Client struct {
		install    InstallService
		addon      AddonService
		preference PreferenceService
		library    LibraryService
		downloader DownloadService
		archiver   ArchiverService
		encoder    EncoderService
		resource   ResourceService
		conf       Config
	}
)

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
		ResourceDir:     filepath.Join(appDir, "resources"),
		Debug:           false,
		Platform:        platform,
	}

	return &conf, nil
}

func NewClient(conf Config) (*Client, error) {
	db, err := scribble.New(conf.DBDir, nil)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(conf.InstallationDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create installation directory: %w", err)
	}

	if err := os.MkdirAll(conf.ResourceDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create resource directory: %w", err)
	}

	client := &Client{
		install:    NewInstallService(db),
		addon:      NewAddonService(db),
		preference: NewPreferenceService(db),
		library:    NewLibraryService(),
		downloader: NewDownloaderService(conf.InstallationDir),
		archiver:   NewArchiverService(true),
		encoder:    NewEncoderService(),
		resource:   NewResourceService(conf.ResourceDir),
		conf:       conf,
	}

	return client, nil
}

func (c *Client) Initialize() error {
	err := c.resource.SaveOut()
	if err != nil {
		return err
	}

	return nil
}

// func (c *Client) Platform() runtime.Platform {
// 	return c.conf.Platform
// }
