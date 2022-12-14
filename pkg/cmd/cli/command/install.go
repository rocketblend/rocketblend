package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func NewInstallCommand(client *client.Client) *cobra.Command {
	var build string

	c := &cobra.Command{
		Use:   "install",
		Short: "Installs a new version of blender from build",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.InstallBuild(reference.Reference(build)); err != nil {
				fmt.Printf("Error installing build: %v\n", err)
				return
			}
		},
	}

	c.Flags().StringVarP(&build, "build", "b", "", "build of the version to install")
	if err := c.MarkFlagRequired("build"); err != nil {
		fmt.Println(err)
	}

	return c
}
