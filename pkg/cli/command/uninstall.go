package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newUninstallCommand creates a new cobra.Command that uninstalls dependencies from the current project.
func (srv *Service) newUninstallCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "uninstall [reference]",
		Short: "Remove project dependencies",
		Long:  "Removes dependencies from the current project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}

	return c
}
