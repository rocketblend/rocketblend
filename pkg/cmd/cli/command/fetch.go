package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func NewFetchCommand(client *client.Client) *cobra.Command {
	var ref string

	c := &cobra.Command{
		Use:   "fetch",
		Short: "fetchs a packs details",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := client.FetchPackByReference(reference.Reference(ref))
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
