package command

import (
	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewOpenCommand(srv *client.Client) *cobra.Command {
	var path string
	// var output string = "json"

	c := &cobra.Command{
		Use:   "open",
		Short: "Opens blender with the specified version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := srv.Open(path, "cmd"); err != nil {
				cmd.PrintErrln(err)
			}
		},
	}

	c.Flags().StringVarP(&path, "path", "p", "", "The path to a .blendfile")
	// c.Flags().StringVarP(&output, "output", "o", "", "Output type of command")

	return c
}
