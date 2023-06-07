package command

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/rocketblend/helpers"
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

			return srv.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				return rocketblend.Run(ctx)
			}, &helpers.SpinnerOptions{Suffix: "Running project..."})
		},
	}

	return c
}
