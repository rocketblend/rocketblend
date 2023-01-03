package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/rocketblend/pkg/core/build"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/jot"
)

func NewAddonService(driver *jot.Driver) *addon.Service {
	return addon.NewService(driver)
}

func NewBuildService(driver *jot.Driver, platform runtime.Platform, addonService build.AddonService) *build.Service {
	return build.NewService(driver, platform, addonService)
}

func NewResourceService(dir string) *resource.Service {
	srv := resource.NewService(dir)
	return srv
}
