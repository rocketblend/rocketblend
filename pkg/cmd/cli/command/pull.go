package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func NewPullCommand(client *client.Client) *cobra.Command {
	var ref string

	c := &cobra.Command{
		Use:   "pull",
		Short: "pulls the actual data defined by a pack",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := client.PullPackByReference(reference.Reference(ref))
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	c.Flags().StringVarP(&ref, "ref", "r", "", "reference of the addon to get (required)")
	if err := c.MarkFlagRequired("ref"); err != nil {
		fmt.Println(err)
	}

	return c
}
