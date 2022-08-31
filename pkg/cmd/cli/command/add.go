package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewAddCommand(client *client.Client) *cobra.Command {
	var path string

	c := &cobra.Command{
		Use:   "add",
		Short: "Adds a local build of blender to the database",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.AddInstall(path); err != nil {
				fmt.Printf("Error adding build: %v\n", err)
				return
			}

			fmt.Printf("Build %s added\n", path)
		},
	}

	c.Flags().StringVarP(&path, "path", "p", "", "path to the version to add")
	if err := c.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	return c
}
