package client

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/executable"
	"github.com/rocketblend/rocketblend/pkg/core/preference"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
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
	build, err := c.build.FindByReference(reference.Reference(ref))
	if err != nil {
		return nil, fmt.Errorf("failed to find build: %s", err)
	}

	addonMap, err := c.getExecutableAddonsByReference(build.Packages)
	if err != nil {
		return nil, fmt.Errorf("failed to find all addons for build: %s", err)
	}

	return &executable.Executable{
		Path:   filepath.Join(c.conf.InstallationDir, ref, build.GetSourceForPlatform(c.conf.Platform).Executable),
		Addons: addonMap,
		ARGS:   build.Args,
	}, nil
}

func (c *Client) getExecutableAddonsByReference(ref []string) (*[]executable.Addon, error) {
	addons := []executable.Addon{}
	for _, r := range ref {
		addon, err := c.getExecutableAddonByReference(r)
		if err != nil {
			return nil, fmt.Errorf("failed to find addon: %s", err)
		}

		addons = append(addons, *addon)
	}

	return &addons, nil
}

func (c *Client) getExecutableAddonByReference(ref string) (*executable.Addon, error) {
	pack, err := c.addon.FindByReference(reference.Reference(ref))
	if err != nil {
		return nil, fmt.Errorf("failed to find addon: %s", err)
	}

	return &executable.Addon{
		Name:    pack.Name,
		Version: pack.AddonVersion,
		Path:    filepath.Join(c.conf.InstallationDir, ref, pack.Source.File),
	}, nil
}
