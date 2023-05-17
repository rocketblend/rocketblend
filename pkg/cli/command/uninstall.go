package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

// newUninstallCommand creates a new cobra.Command that uninstalls dependencies from the current project.
// It optionally allows the dependencies to be uninstalled globally.
func (srv *Service) newUninstallCommand() *cobra.Command {
	var global bool

	c := &cobra.Command{
		Use:   "uninstall [reference]",
		Short: "Remove project dependencies",
		Long:  "Removes dependencies from the current project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := srv.parseReference(args[0])
			if err != nil {
				return fmt.Errorf("invalid reference: %w", err)
			}

			if global {
				return srv.uninstallGlobalPackage(ref)
			}

			return srv.uninstallProjectDependencies(ref)
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "uninstall dependencies globally")

	return c
}

// uninstallGlobalPackage uninstalls a package globally by its reference.
func (srv *Service) uninstallGlobalPackage(ref *reference.Reference) error {
	err := srv.driver.UninstallPackByReference(*ref)
	if err != nil {
		return fmt.Errorf("failed to uninstall global package: %w", err)
	}

	return nil
}

// uninstallProjectDependencies uninstalls dependencies of the current project by reference.
func (srv *Service) uninstallProjectDependencies(ref *reference.Reference) error {
	err := srv.driver.UninstallDependencies(srv.flags.workingDirectory, *ref)
	if err != nil {
		return fmt.Errorf("failed to uninstall project dependencies: %w", err)
	}

	return nil
}
