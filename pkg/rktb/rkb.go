package rktb

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/pkg/rktb/command"
	"github.com/rocketblend/rocketblend/pkg/rktb/config"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"
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

	rocketblendOptions := rocketblend.Options{
		Debug:                  config.Debug,
		Platform:               config.Platform,
		InstallationsDirectory: config.InstallDir,
		AddonsEnabled:          config.Features.Addons,
	}

	driver, err := rocketblend.New(&rocketblendOptions)
	if err != nil {
		return err
	}

	srv := command.NewService(cs, driver)
	rootCMD := srv.NewCommand()

	// Configure help template colours
	cc.Init(&cc.Config{
		RootCmd:         rootCMD,
		Headings:        cc.Cyan + cc.Bold + cc.Underline,
		Commands:        cc.Bold,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		Aliases:         cc.Bold,
		Example:         cc.Italic,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	return rootCMD.Execute()
}
