package cli

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/command"
	"github.com/spf13/cobra"
)

func New() (*cobra.Command, error) {
	cmd, err := command.NewService()
	if err != nil {
		return nil, err
	}

	rootCMD := cmd.NewRootCommand()

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
