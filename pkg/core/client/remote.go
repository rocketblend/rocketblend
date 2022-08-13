package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/remote"
	"github.com/rocketblend/rocketblend/pkg/scribble"
)

func NewRemoteService(driver *scribble.Driver) *remote.Service {
	repo := remote.NewScribbleRepository(driver)
	conf := remote.LoadConfig()
	srv := remote.NewService(conf, repo)

	return srv
}
