package client

import "github.com/rocketblend/rocketblend/pkg/jot/reference"

func (c *Client) InstallBuild(ref reference.Reference) error {
	err := c.build.FetchByReference(ref)
	if err != nil {
		return err
	}

	err = c.build.PullByReference(ref, c.conf.Platform)
	if err != nil {
		return err
	}

	return nil
}
