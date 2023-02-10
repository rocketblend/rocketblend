package cli

import (
	cc "github.com/ivanpirog/coloredcobra"
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
