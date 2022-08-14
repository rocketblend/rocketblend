package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewInstallCommand(client *client.Client) *cobra.Command {
	var build string

	c := &cobra.Command{
		Use:   "install",
		Short: "Installs a new verison of blender into the local repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.InstallBuild(build); err != nil {
				fmt.Printf("Error installing build: %v\n", err)
			}
		},
	}

	c.Flags().StringVarP(&build, "build", "b", "", "Build hash of the version to install")
	if err := c.MarkFlagRequired("build"); err != nil {
		fmt.Println(err)
	}

	return c
}
