package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newRunCommand creates a new cobra command for running the project.
// It sets up all necessary flags and executes the project through the driver.
func (srv *Service) newRunCommand() *cobra.Command {
	var background bool

	c := &cobra.Command{
		Use:   "run",
		Short: "Runs the project",
		Long:  `Launches the project in the current working directory. Can optionally run in the background.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("unable to locate project file: %w", err)
			}

			err = srv.run(blend, background, []string{})
			if err != nil {
				return fmt.Errorf("failed to run project: %w", err)
			}

			return nil
		},
	}

	// background flag allows the project to be run in the background.
	c.Flags().BoolVarP(&background, "background", "b", false, "run the project in the background")

	return c
}
