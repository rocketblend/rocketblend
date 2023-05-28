package cli

import (
	"github.com/flowshot-io/x/pkg/logger"
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

	opts := []rocketblend.Option{
		rocketblend.WithInstallationDirectory(config.InstallDir),
		rocketblend.WithPlatform(config.Platform),
		rocketblend.WithLogger(logger.New(logger.WithPretty())),
	}

	// TODO: Remove this and just use log level
	if config.Debug {
		opts = append(opts, rocketblend.WithDebug())
	}

	if config.Features.Addons {
		opts = append(opts, rocketblend.WithAddonsEnabled())
	}

	driver, err := rocketblend.New(opts...)
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
