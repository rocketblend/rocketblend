package command

import (
	"github.com/spf13/cobra"
)

// newRunCommand creates a new cobra command for running the project.
func (srv *Service) newRunCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "run",
		Short: "Runs the project",
		Long:  `Launches the project in the current working directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rocketblend, err := srv.getDriver()
			if err != nil {
				return err
			}

			return rocketblend.Run(cmd.Context())
		},
	}

	return c
}
