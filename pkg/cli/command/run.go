package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newRunCommand creates a new 'run' command.
// This command runs the project located in the working directory.
// It optionally allows the project to be run in the background.
func (srv *Service) newRunCommand() *cobra.Command {
	var background bool

	c := &cobra.Command{
		Use:   "run",
		Short: "Run project",
		Long:  `Runs the project located in the current working directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("unable to locate project file: %w", err)
			}

			err = srv.driver.Run(blend, background, []string{})
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
