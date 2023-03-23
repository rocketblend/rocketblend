package cli

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/config"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/generator"
	"github.com/rocketblend/rocketblend/pkg/core"
	"github.com/spf13/cobra"
)

func setup() (*cobra.Command, error) {
	cs, err := config.New()
	if err != nil {
		return nil, err
	}

	config, err := cs.Get()
	if err != nil {
		return nil, err
	}

	rocketblendOptions := core.Options{
		Debug:         config.Debug,
		Platform:      config.Platform,
		DefaultBuild:  config.DefaultBuild,
		AddonsEnabled: config.Features.Addons,
	}

	driver, err := core.New(&rocketblendOptions)
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
		Aliases:         cc.Bold,
		Example:         cc.Italic,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	return rootCMD, nil
}

func GenerateDocs(path string) error {
	cmd, err := setup()
	if err != nil {
		return err
	}

	// err = doc.GenMarkdownTree(cmd, path)
	// if err != nil {
	// 	return err
	// }

	err = generator.MarkdownTree(cmd, path)
	if err != nil {
		return err
	}

	return nil
}

func Execute() error {
	rootCMD, err := setup()
	if err != nil {
		return err
	}

	return rootCMD.Execute()
}
