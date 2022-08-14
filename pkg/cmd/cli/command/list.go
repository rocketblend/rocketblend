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
		Short: "List all the versions of blender that are available",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			build, err := srv.GetAvilableBuilds("windows", filterTag)
			if err != nil {
				fmt.Printf("Error getting available builds: %v\n", err)
			}

			for _, build := range build {
				fmt.Printf("(%s)\t%s\t\t%s\n", build.Hash, build.Name, build.Uri)
			}
		},
	}

	c.Flags().StringVarP(&filterTag, "tag", "t", "", "Filter by tag")

	return c
}
