package cli

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/pkg/cli/command"
	"github.com/rocketblend/rocketblend/pkg/cli/config"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"
	"github.com/spf13/cobra"
)

func New() (*cobra.Command, error) {
	cs, err := config.New()
	if err != nil {
		return nil, err
	}

	config, err := cs.Get()
	if err != nil {
		return nil, err
	}

	rocketblendOptions := rocketblend.Options{
		Debug:                 config.Debug,
		Platform:              config.Platform,
		InstallationDirectory: config.InstallDir,
		AddonsEnabled:         config.Features.Addons,
	}

	driver, err := rocketblend.New(&rocketblendOptions)
	if err != nil {
		return nil, err
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
		Aliases:         cc.Magenta,
		Example:         cc.Green + cc.Italic,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	return rootCMD, nil
}
