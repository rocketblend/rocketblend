package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// newRunCommand creates a new cobra command for running the project.
func (c *cli) newRunCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "run",
		Short: "Runs the project",
		Long:  `Launches the project in the current working directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rocketblend, err := c.getDriver()
			if err != nil {
				return err
			}

			return c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				return rocketblend.Run(ctx)
			}, &spinnerOptions{Suffix: "Running project..."})
		},
	}

	return cc
}
