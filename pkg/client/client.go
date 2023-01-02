package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/rocketblend/pkg/core/build"
	"github.com/rocketblend/rocketblend/pkg/core/preference"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/scribble"
)

type (
	PreferenceService interface {
		Find() (*preference.Settings, error)
		Create(i *preference.Settings) error
	}

	ResourceService interface {
		FindByName(name string) (*resource.Resource, error)
		SaveOut() error
	}

	AddonService interface {
		FindByReference(ref reference.Reference) (*addon.Package, error)
		FetchByReference(ref reference.Reference) error
		PullByReference(ref reference.Reference) error
	}

	BuildService interface {
		FindByReference(ref reference.Reference) (*build.Build, error)
		FetchByReference(ref reference.Reference) error
		PullByReference(ref reference.Reference, platform runtime.Platform) error
	}

	Config struct {
		DBDir           string
		InstallationDir string
		ResourceDir     string
		Debug           bool
		Platform        runtime.Platform
	}

	Client struct {
		preference PreferenceService
		resource   ResourceService
		addon      AddonService
		build      BuildService
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

	appDir := filepath.Join(home, "rocketblend")
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

	jot, err := jot.New(conf.InstallationDir, nil)
	if err != nil {
		return nil, err
	}

	addonService := NewAddonService(jot)

	client := &Client{
		preference: NewPreferenceService(db),
		addon:      addonService,
		build:      NewBuildService(jot, addonService),
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
