package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (srv *Service) newInitCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initialize rocketblend",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := srv.Initialize(); err != nil {
				fmt.Printf("failed to initialize: %s", err)
				return
			}

			fmt.Println("Initialized rocketblend")
		},
	}

	return c
}
