package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewListCommand(srv *client.Client) *cobra.Command {
	var filterTag string

	c := &cobra.Command{
		Use:   "list",
		Short: "List all the versions currently installed",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			installs, err := srv.FindAllInstalls()
			if err != nil {
				fmt.Printf("Error finding blender installs: %v\n", err)
			}

			for _, install := range installs {
				fmt.Printf("%s - %s\n", install.Path, install.Build)
			}
		},
	}

	c.Flags().StringVarP(&filterTag, "tag", "t", "", "Filter by tag")

	return c
}
