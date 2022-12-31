package client

import (
	"fmt"

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

func (c *Client) newAddon(ref string, path string) *addon.Addon {
	return &addon.Addon{
		Id:      c.encoder.Hash(ref),
		Package: ref,
		Path:    path,
	}
}
