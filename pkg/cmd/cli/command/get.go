package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func NewGetCommand(client *client.Client) *cobra.Command {
	var ref string

	c := &cobra.Command{
		Use:   "get",
		Short: "Gets a addon for blender",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.InstallAddon(reference.Reference(ref)); err != nil {
				fmt.Printf("Error installing addon: %v\n", err)
				return
			}

			fmt.Printf("addon %s added to collection\n", ref)
		},
	}

	c.Flags().StringVarP(&ref, "ref", "r", "", "reference of the addon to get (required)")
	if err := c.MarkFlagRequired("ref"); err != nil {
		fmt.Println(err)
	}

	return c
}
