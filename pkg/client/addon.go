package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/scribble"
)

func NewAddonService(driver *scribble.Driver) *addon.Service {
	repo := addon.NewScribbleRepository(driver)
	srv := addon.NewService(repo)

	return srv
}
