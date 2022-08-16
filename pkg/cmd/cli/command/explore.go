package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewExploreCommand(srv *client.Client) *cobra.Command {
	var filterTag string

	c := &cobra.Command{
		Use:   "explore",
		Short: "Explore all the versions of blender that are available from remotes.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			build, err := srv.FetchRemoteBuilds("windows")
			if err != nil {
				fmt.Printf("Error fetching remote builds: %v\n", err)
			}

			for _, build := range build {
				fmt.Printf("(%s)\t%s\t\t%s\n", build.Hash, build.Name, build.DownloadUrl)
			}
		},
	}

	c.Flags().StringVarP(&filterTag, "tag", "t", "", "Filter by tag")

	return c
}
