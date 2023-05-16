package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newInstallCommand() *cobra.Command {
	var global bool
	var force bool

	c := &cobra.Command{
		Use:   "install [reference]",
		Short: "Install project dependencies",
		Long:  `Adds dependencies to the current project and installs them.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var ref *reference.Reference
			var err error

			if len(args) > 0 {
				ref, err = srv.parseReference(args[0])
				if err != nil {
					return err
				}
			}

			if global {
				return srv.installGlobal(ref, force)
			}

			return srv.installLocal(ref, force)
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "install dependencies globally")
	c.Flags().BoolVarP(&force, "force", "f", false, "force install dependencies (even if they are already installed)")

	return c
}

// installGlobal installs a package globally by its reference.
func (srv *Service) installGlobal(ref *reference.Reference, force bool) error {
	if ref != nil {
		err := srv.driver.InstallPackByReference(*ref, force)
		if err != nil {
			return fmt.Errorf("failed to install package: %w", err)
		}
	}

	return nil
}

// installLocal installs dependencies of the current project by reference.
func (srv *Service) installLocal(ref *reference.Reference, force bool) error {
	err := srv.driver.InstallDependencies(srv.flags.workingDirectory, ref, force)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}
