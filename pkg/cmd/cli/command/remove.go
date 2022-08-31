package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewRemoveCommand(srv *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "remove",
		Short: "Removes a version of blender from the local repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := srv.RemoveInstall(args[0]); err != nil {
				fmt.Printf("Error removing build: %v\n", err)
				return
			}

			fmt.Printf("Build %s removed\n", args[0])
		},
	}

	return c
}
