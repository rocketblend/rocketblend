package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/build"
)

func NewBuildService() *build.Service {
	conf := build.NewConfig()
	http := build.NewHttpClient()
	srv := build.NewService(conf, http)

	return srv
}
