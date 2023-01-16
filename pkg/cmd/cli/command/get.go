package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newGetCommand() *cobra.Command {
	var ref string

	c := &cobra.Command{
		Use:   "get",
		Short: "gets a packge and installs it ready for use",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := srv.getPackByReference(reference.Reference(ref))
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
