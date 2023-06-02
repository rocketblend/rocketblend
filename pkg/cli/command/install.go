package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newInstallCommand creates a new cobra command for installing project dependencies.
// It allows for global installation, and forceful installation of dependencies.
func (srv *Service) newInstallCommand() *cobra.Command {
	var global bool
	var force bool

	c := &cobra.Command{
		Use:   "install [reference]",
		Short: "Installs project dependencies",
		Long:  `Adds the specified dependencies to the current project and installs them. If no reference is provided, all dependencies in the project are installed instead.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}

	// Add 'global' and 'force' flags to the command
	c.Flags().BoolVarP(&global, "global", "g", false, "Installs dependencies globally, affecting all projects.")
	c.Flags().BoolVarP(&force, "force", "f", false, "Forces the installation of dependencies, even if they are already installed.")

	return c
}
