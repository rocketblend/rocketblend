package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newFindCommand() *cobra.Command {
	var referenceStr string

	c := &cobra.Command{
		Use:   "find [flags]",
		Short: "Find if a package definition already exists",
		Long:  `find if a package definition already exists in the global cache`,
		Run: func(cmd *cobra.Command, args []string) {
			reference, err := reference.Parse(referenceStr)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			pack, err := srv.driver.FindPackByReference(reference)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(pack.ToString())
		},
	}

	c.Flags().StringVarP(&referenceStr, "reference", "r", "", "reference of the addon to get (required)")
	c.MarkFlagRequired("reference")

	return c
}
