package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newFetchCommand() *cobra.Command {
	var referenceStr string

	c := &cobra.Command{
		Use:   "fetch [flags]",
		Short: "Fetch a package defination via the given reference",
		Long:  `fetch a package defination via the given reference and stores it in the global cache`,
		Run: func(cmd *cobra.Command, args []string) {
			reference, err := reference.Parse(referenceStr)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = srv.driver.FetchPackByReference(reference)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	c.Flags().StringVarP(&referenceStr, "reference", "r", "", "reference of the addon to get (required)")
	c.MarkFlagRequired("reference")

	return c
}
