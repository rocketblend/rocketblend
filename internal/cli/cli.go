package cli

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/internal/cli/command"
	"github.com/spf13/cobra"
)

const Name = "rocketblend"

var Version = "dev"

func New() *cobra.Command {
	rootCMD := command.NewRootCommand(&command.RootCommandOpts{
		Name:    Name,
		Version: Version,
	})

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

	return rootCMD
}
