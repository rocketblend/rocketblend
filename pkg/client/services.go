package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/rocketblend/pkg/core/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/jot"
)

func NewResourceService(dir string) *resource.Service {
	srv := resource.NewService(dir)
	return srv
}

func NewPackService(driver *jot.Driver, platform runtime.Platform) *rocketpack.Service {
	return rocketpack.NewService(driver, platform)
}
