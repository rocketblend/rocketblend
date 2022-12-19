package client

import "github.com/rocketblend/rocketblend/pkg/core/resource"

func (c *Client) FindResource(key string) (*resource.Resource, error) {
	// TODO: save out if not found.
	return c.resource.FindByName(key)
}
