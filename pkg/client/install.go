package client

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/core/install"
)

func (c *Client) FindAllInstalls() ([]*install.Install, error) {
	ints, err := c.install.FindAll()
	if err != nil {
		return nil, err
	}

	return ints, nil
}

func (c *Client) AddInstall(path string) error {
	build, err := c.library.FindBuildByPath(path)
	if err != nil {
		return err
	}

	i, _ := c.findInstall(build.Reference)
	if i != nil {
		return fmt.Errorf("build already installed")
	}

	err = c.install.Create(c.newInstall(build.Reference, path, build.Packages))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveInstall(hash string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) findInstall(ref string) (*install.Install, error) {
	id := c.encoder.Hash(ref)
	install, err := c.install.FindByID(id)
	if err != nil {
		return nil, err
	}

	return install, nil
}

func (c *Client) newInstall(ref string, path string, packs []string) *install.Install {
	return &install.Install{
		Id:       c.encoder.Hash(ref),
		Build:    ref,
		Path:     path,
		Packages: packs,
	}
}
