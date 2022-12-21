package command

import (
	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewOpenCommand(srv *client.Client) *cobra.Command {
	var path string

	c := &cobra.Command{
		Use:   "open",
		Short: "Opens blender with the specified version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := srv.Open(path); err != nil {
				cmd.PrintErrln(err)
			}
		},
	}

	c.Flags().StringVarP(&path, "path", "p", "", "The path to a .blendfile")

	return c
}
