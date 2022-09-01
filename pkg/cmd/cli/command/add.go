package command

import (
	"fmt"
	"path"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/core/library"
	"github.com/spf13/cobra"
)

func NewAddCommand(client *client.Client) *cobra.Command {
	var dir string

	c := &cobra.Command{
		Use:   "add",
		Short: "Adds a local build/package to the database",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			name := path.Base(dir)

			switch name {
			case library.BuildFile:
				err := addInstall(client, dir)
				if err != nil {
					fmt.Printf("Error adding build: %v\n", err)
				}
			case library.PackgeFile:
				err := addAddon(client, dir)
				if err != nil {
					fmt.Printf("Error adding build: %v\n", err)
				}
			default:
				fmt.Printf("Unknown file type: %s\n", name)
			}
		},
	}

	c.Flags().StringVarP(&dir, "path", "p", "", "path to the build/package to add")
	if err := c.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	return c
}

func addInstall(client *client.Client, path string) error {
	if err := client.AddInstall(path); err != nil {
		return fmt.Errorf("error adding build: %v", err)
	}

	fmt.Printf("Build %s added\n", path)

	return nil
}

func addAddon(client *client.Client, path string) error {
	if err := client.AddAddon(path); err != nil {
		return fmt.Errorf("error adding package: %v", err)
	}

	fmt.Printf("Package %s added\n", path)

	return nil
}
