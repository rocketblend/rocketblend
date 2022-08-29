package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewLibraryInstallCommand(client *client.Client) *cobra.Command {
	var build string

	c := &cobra.Command{
		Use:   "library",
		Short: "Installs a new build of blender from the library",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			b, err := client.FetchBuild(build)
			if err != nil {
				fmt.Printf("Error fetching build: %v\n", err)
				return
			}

			fmt.Println(b)
			// if err := client.InstallBuild(build); err != nil {
			// 	fmt.Printf("Error installing build: %v\n", err)
			// }
		},
	}

	c.Flags().StringVarP(&build, "build", "b", "", "build to install")
	if err := c.MarkFlagRequired("build"); err != nil {
		fmt.Println(err)
	}

	return c
}
