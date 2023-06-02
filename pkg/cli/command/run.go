package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newRunCommand creates a new cobra command for running the project.
func (srv *Service) newRunCommand() *cobra.Command {
	var background bool

	c := &cobra.Command{
		Use:   "run",
		Short: "Runs the project",
		Long:  `Launches the project in the current working directory. Can optionally run in the background.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}

	// background flag allows the project to be run in the background.
	c.Flags().BoolVarP(&background, "background", "b", false, "run the project in the background")

	return c
}
