package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (srv *Service) newStartCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Starts the project",
		Long:  `Starts the project located in the current working directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}

	return c
}
