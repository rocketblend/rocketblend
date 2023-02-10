package cli

import (
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/config"
	"github.com/rocketblend/rocketblend/pkg/core"
)

func Execute() error {
	cs, err := config.New()
	if err != nil {
		return err
	}

	config, err := cs.Get()
	if err != nil {
		return err
	}

	rocketblendOptions := core.Options{
		Debug:         config.Debug,
		Platform:      config.Platform,
		DefaultBuild:  config.DefaultBuild,
		AddonsEnabled: config.Features.Addons,
	}

	driver, err := core.New(&rocketblendOptions)
	if err != nil {
		return err
	}

	srv := command.NewService(cs, driver)
	rootCMD := srv.NewCommand()

	return rootCMD.Execute()
}
