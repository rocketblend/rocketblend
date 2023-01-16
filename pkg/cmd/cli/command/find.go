package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newFindCommand() *cobra.Command {
	var ref string

	c := &cobra.Command{
		Use:   "find",
		Short: "Finds a packs details locally and validates it",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			pack, err := srv.driver.FindPackByReference(reference.Reference(ref))
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(pack.ToString())
		},
	}

	c.Flags().StringVarP(&ref, "ref", "r", "", "reference of the addon to get (required)")
	if err := c.MarkFlagRequired("ref"); err != nil {
		fmt.Println(err)
	}

	return c
}
