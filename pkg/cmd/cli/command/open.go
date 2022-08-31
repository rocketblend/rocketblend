package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewOpenCommand(srv *client.Client) *cobra.Command {
	var project string
	var build string
	var arguments string

	c := &cobra.Command{
		Use:   "open",
		Short: "Opens blender with the specified version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := srv.OpenProject(project, build, arguments); err != nil {
				fmt.Printf("error trying to open blender: %v\n", err)
				return
			}
		},
	}

	c.Flags().StringVarP(&build, "build", "b", "", "Build reference override")
	c.Flags().StringVarP(&project, "project", "p", "", "The path to a .blendfile")
	c.Flags().StringVarP(&arguments, "args", "a", "", "Pass through arguments for blender")

	return c
}
