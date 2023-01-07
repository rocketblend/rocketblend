package client

import "github.com/rocketblend/rocketblend/pkg/jot/reference"

func (c *Client) InstallAddon(ref reference.Reference) error {
	err := c.addon.FetchByReference(ref)
	if err != nil {
		return err
	}

	err = c.addon.PullByReference(ref)
	if err != nil {
		return err
	}

	return nil
}
