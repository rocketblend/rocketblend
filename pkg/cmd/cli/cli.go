package cli

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command"
	"github.com/rocketblend/rocketblend/pkg/core"
)

func Execute() error {
	driver, err := core.New(nil)
	if err != nil {
		return err
	}

	srv := command.NewService(driver)
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
