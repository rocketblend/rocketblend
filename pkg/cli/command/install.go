package command

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
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
			var ref *reference.Reference
			var err error

			// Parse the reference if provided
			if len(args) > 0 {
				ref, err = srv.parseReference(args[0])
				if err != nil {
					return err
				}
			}

			// Depending on the 'global' flag, perform a global or local installation
			if global {
				return srv.installGlobal(cmd.Context(), ref, force)
			}

			return srv.installLocal(cmd.Context(), ref, force)
		},
	}

	// Add 'global' and 'force' flags to the command
	c.Flags().BoolVarP(&global, "global", "g", false, "Installs dependencies globally, affecting all projects.")
	c.Flags().BoolVarP(&force, "force", "f", false, "Forces the installation of dependencies, even if they are already installed.")

	return c
}

// installGlobal installs a package globally by its reference.
func (srv *Service) installGlobal(ctx context.Context, ref *reference.Reference, force bool) error {
	if ref != nil {
		err := srv.driver.InstallPackByReference(ctx, *ref, force)
		if err != nil {
			return fmt.Errorf("failed to install package: %w", err)
		}
	}

	return nil
}

// installLocal installs dependencies of the current project by reference.
func (srv *Service) installLocal(ctx context.Context, ref *reference.Reference, force bool) error {
	err := srv.driver.InstallDependencies(ctx, srv.flags.workingDirectory, ref, force)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}
