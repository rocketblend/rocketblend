package client

import (
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/client/config"
	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/rocketblend/pkg/core/build"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/rocketblend/pkg/core/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

type (
	ResourceService interface {
		FindByName(name string) (*resource.Resource, error)
		SaveOut() error
	}

	AddonService interface {
		FindByReference(ref reference.Reference) (*addon.Addon, error)
		FetchByReference(ref reference.Reference) error
		PullByReference(ref reference.Reference) error
	}

	BuildService interface {
		FindByReference(ref reference.Reference) (*build.Build, error)
		FetchByReference(ref reference.Reference) error
		PullByReference(ref reference.Reference) error
	}

	PackService interface {
		FindByReference(ref reference.Reference) (*rocketpack.RocketPack, error)
		FetchByReference(ref reference.Reference) error
		PullByReference(ref reference.Reference) error
	}

	Client struct {
		resource ResourceService
		addon    AddonService
		build    BuildService
		pack     PackService
		conf     *config.Config
	}
)

func New() (*Client, error) {
	config, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := os.MkdirAll(config.Directories.Resources, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create resource directory: %w", err)
	}

	jot, err := jot.New(config.Directories.Installations, nil)
	if err != nil {
		return nil, err
	}

	addonService := NewAddonService(jot)
	buildService := NewBuildService(jot, config.Platform, addonService)
	packService := NewPackService(jot, config.Platform)
	resourceService := NewResourceService(config.Directories.Resources)

	client := &Client{
		addon:    addonService,
		build:    buildService,
		pack:     packService,
		resource: resourceService,
		conf:     config,
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
