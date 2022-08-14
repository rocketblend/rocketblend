package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewOpenCommand(srv *client.Client) *cobra.Command {
	var filePath string
	var blenderArgs string

	c := &cobra.Command{
		Use:   "open",
		Short: "Opens a project",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("open called")
		},
	}

	c.Flags().StringVarP(&filePath, "file", "f", "", "The path to the .blendfile to open")
	c.Flags().StringVarP(&blenderArgs, "args", "a", "", "Arguments to pass to blender")

	return c
}
