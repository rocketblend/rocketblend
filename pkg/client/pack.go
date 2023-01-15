package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

func (c *Client) FetchPackByReference(ref reference.Reference) error {
	err := c.pack.FetchByReference(ref)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) PullPackByReference(ref reference.Reference) error {
	err := c.pack.PullByReference(ref)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) FindPackByReference(ref reference.Reference) (*rocketpack.RocketPack, error) {
	pack, err := c.pack.FindByReference(ref)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (c *Client) GetPackByReference(ref reference.Reference) error {
	err := c.pack.FetchByReference(ref)
	if err != nil {
		return err
	}

	err = c.pack.PullByReference(ref)
	if err != nil {
		return err
	}

	return nil
}
