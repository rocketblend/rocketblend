package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli/client"
	"github.com/spf13/cobra"
)

func NewCreateCommand(srv *client.Client) *cobra.Command {
	var projectName string
	var path string

	c := &cobra.Command{
		Use:   "create",
		Short: "Creates a new rocketfile",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create called")
		},
	}

	c.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	if err := c.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}

	c.Flags().StringVarP(&path, "path", "p", "", "Path to create the project at")

	return c
}
