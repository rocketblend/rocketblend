package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newPullCommand() *cobra.Command {
	var referenceStr string

	c := &cobra.Command{
		Use:   "pull [flags]",
		Short: "Pull the source files defined by a package",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			reference, err := reference.Parse(referenceStr)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = srv.driver.PullPackByReference(reference)
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
