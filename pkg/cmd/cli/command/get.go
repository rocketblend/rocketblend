package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newGetCommand() *cobra.Command {
	var referenceStr string

	c := &cobra.Command{
		Use:   "get [flags]",
		Short: "Add a package dependency to the current project and install it",
		Long: `resolves the given reference to a package, updates rocketblend.yaml to require the reference,
and downloads the package source to the global cache`,
		Run: func(cmd *cobra.Command, args []string) {
			reference, err := reference.Parse(referenceStr)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = srv.getPackByReference(reference)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	c.Flags().StringVarP(&referenceStr, "reference", "r", "", "reference of the package to get (required)")
	c.MarkFlagRequired("reference")

	return c
}

func (srv *Service) getPackByReference(ref reference.Reference) error {
	err := srv.driver.FetchPackByReference(ref)
	if err != nil {
		return err
	}

	err = srv.driver.PullPackByReference(ref)
	if err != nil {
		return err
	}

	return nil
}
