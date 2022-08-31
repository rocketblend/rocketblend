package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewGetCommand(client *client.Client) *cobra.Command {
	var pack string

	c := &cobra.Command{
		Use:   "get",
		Short: "Gets a new packge/addon for blender",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.InstallPackage(pack); err != nil {
				fmt.Printf("Error installing package: %v\n", err)
				return
			}

			fmt.Printf("Package %s added to collection\n", pack)
		},
	}

	c.Flags().StringVarP(&pack, "package", "p", "", "package of the addon get")
	if err := c.MarkFlagRequired("package"); err != nil {
		fmt.Println(err)
	}

	return c
}
