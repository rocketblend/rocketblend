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
		Short: "Opens a rocketfile",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("open called")
		},
	}

	c.Flags().StringVarP(&filePath, "file", "f", ".rocket", "The path to the rocketfile to open")
	c.Flags().StringVarP(&blenderArgs, "args", "a", "", "Arguments to pass to blender")

	return c
}
