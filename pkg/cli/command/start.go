package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newStartCommand creates a new cobra.Command that starts the project located in the working directory.
// It optionally allows the project to be started in the background.
func (srv *Service) newStartCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start project",
		Long:  `Starts the project located in the current working directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("unable to locate project file: %w", err)
			}

			err = srv.driver.Start(blend, []string{})
			if err != nil {
				return fmt.Errorf("failed to start project: %w", err)
			}

			return nil
		},
	}

	return c
}
