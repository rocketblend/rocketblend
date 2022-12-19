package client

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/addon"
)

func (c *Client) FindAllAddons() ([]*addon.Addon, error) {
	addons, err := c.addon.FindAll()
	if err != nil {
		return nil, err
	}

	return addons, nil
}

func (c *Client) AddAddon(path string) error {
	pack, err := c.library.FindPackageByPath(path)
	if err != nil {
		return err
	}

	a, _ := c.findAddon(pack.Reference)
	if a != nil {
		return fmt.Errorf("addon already installed")
	}

	err = c.addon.Create(c.newAddon(pack.Reference, path))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveAddon(hash string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) findAddon(ref string) (*addon.Addon, error) {
	id := c.encoder.Hash(ref)
	addon, err := c.addon.FindByID(id)
	if err != nil {
		return nil, err
	}

	return addon, nil
}

func (c *Client) getAddonMapByReferences(ref []string) (map[string]string, error) {
	addonMap := make(map[string]string)
	for _, a := range ref {
		name, path, err := c.getAddonNamePathByReference(a)
		if err != nil {
			return nil, fmt.Errorf("failed to find package: %s", err)
		}

		addonMap[name] = path
	}

	return addonMap, nil
}

func (c *Client) getAddonNamePathByReference(ref string) (string, string, error) {
	addon, err := c.findAddon(ref)
	if err != nil {
		return "", "", fmt.Errorf("failed to find addon: %s", err)
	}

	pack, err := c.library.FindPackageByPath(addon.Path)
	if err != nil {
		return "", "", fmt.Errorf("failed to find package: %s", err)
	}

	return pack.Name, filepath.Join(addon.Path, pack.Source.File), nil
}

func (c *Client) newAddon(ref string, path string) *addon.Addon {
	return &addon.Addon{
		Id:      c.encoder.Hash(ref),
		Package: ref,
		Path:    path,
	}
}
