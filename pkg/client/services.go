package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/rocketblend/pkg/core/build"
	"github.com/rocketblend/rocketblend/pkg/core/preference"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/scribble"
)

func NewPreferenceService(driver *scribble.Driver) *preference.Service {
	repo := preference.NewScribbleRepository(driver)
	srv := preference.NewService(repo)

	return srv
}

func NewAddonService(driver *jot.Driver) *addon.Service {
	return addon.NewService(driver)
}

func NewBuildService(driver *jot.Driver, addonService build.AddonService) *build.Service {
	return build.NewService(driver, addonService)
}

func NewResourceService(dir string) *resource.Service {
	srv := resource.NewService(dir)
	return srv
}
