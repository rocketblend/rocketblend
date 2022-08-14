package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/scribble"
)

func NewInstallService(driver *scribble.Driver) *install.Service {
	repo := install.NewScribbleRepository(driver)
	conf := install.LoadConfig()
	srv := install.NewService(conf, repo)

	return srv
}
