package client

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/executable"
	"github.com/rocketblend/rocketblend/pkg/core/preference"
)

func (c *Client) SetDefaultExecutable(ref string) error {
	settings, err := c.preference.Find()
	if err != nil {
		return err
	}

	if settings == nil {
		settings = &preference.Settings{}
	}

	settings.DefaultBuild = ref

	err = c.preference.Create(settings)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) findDefaultExecutable() (*executable.Executable, error) {
	settings, err := c.preference.Find()
	if err != nil {
		return nil, err
	}

	if settings.DefaultBuild == "" {
		// TODO: Get latest build and set as default
		return nil, fmt.Errorf("no default executable set")
	}

	executable, err := c.findExecutableByBuildReference(settings.DefaultBuild)
	if err != nil {
		return nil, err
	}

	return executable, nil
}

func (c *Client) findExecutableByBuildReference(ref string) (*executable.Executable, error) {
	if ref == "" {
		return nil, fmt.Errorf("invalid build reference")
	}

	install, err := c.findInstall(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to find install: %s", err)
	}

	build, err := c.library.FindBuildByPath(install.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to find build: %s", err)
	}

	addonMap, err := c.getAddonMapByReferences(build.Packages)
	if err != nil {
		return nil, fmt.Errorf("failed to find all addons for build: %s", err)
	}

	return &executable.Executable{
		Path:   filepath.Join(install.Path, build.GetSourceForPlatform(c.conf.Platform).Executable),
		Addons: addonMap,
		ARGS:   build.Args,
	}, nil
}
