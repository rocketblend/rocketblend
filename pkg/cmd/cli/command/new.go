package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewCreateCommand(srv *client.Client) *cobra.Command {
	var name string
	var path string
	var build string

	c := &cobra.Command{
		Use:   "new",
		Short: "Creates a new project",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := srv.CreateProject(name, path, build); err != nil {
				fmt.Printf("Error creating project: %v\n", err)
				return
			}

			fmt.Printf("Project %s created\n", name)
		},
	}

	c.Flags().StringVarP(&name, "name", "n", "", "Name of the project")
	if err := c.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}

	c.Flags().StringVarP(&path, "path", "p", "", "Path to create the project at")
	c.Flags().StringVarP(&build, "build", "b", "", "Build reference to use")

	return c
}
