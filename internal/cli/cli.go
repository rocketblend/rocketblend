package cli

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/internal/cli/command"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

var Version = "dev"

func New() *cobra.Command {
	rootCMD := command.NewRootCommand(&command.RootCommandOpts{
		Name:    types.ApplicationName,
		Version: Version,
	})

	// Configure help template colours
	cc.Init(&cc.Config{
		RootCmd:         rootCMD,
		Headings:        cc.Underline,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	return rootCMD
}
