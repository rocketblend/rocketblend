package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli/client"
	"github.com/spf13/cobra"
)

func NewInstallCommand(client *client.Client) *cobra.Command {
	var buildHash string

	c := &cobra.Command{
		Use:   "install",
		Short: "Installs a new verison of blender into the local repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.InstallBuild(buildHash); err != nil {
				fmt.Printf("Error installing build: %v\n", err)
			}
		},
	}

	c.Flags().StringVarP(&buildHash, "build", "b", "", "Build hash of the version to install")
	c.MarkFlagRequired("build")

	return c
}
