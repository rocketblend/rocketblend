package command

import (
	"github.com/rocketblend/rocketblend/pkg/rocketblend/reference"
	"github.com/spf13/cobra"
)

// newInstallCommand creates a new cobra command for installing project dependencies.
func (srv *Service) newInstallCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "install [reference]",
		Short: "Installs project dependencies",
		Long:  `Adds the specified dependencies to the current project and installs them. If no reference is provided, all dependencies in the project are installed instead.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rocketblend, err := srv.getDriver()
			if err != nil {
				return err
			}

			if len(args) > 0 {
				ref, err := reference.Parse(args[0])
				if err != nil {
					return err
				}

				err = rocketblend.AddDependencies(cmd.Context(), ref)
				if err != nil {
					return err
				}
			}

			return rocketblend.InstallDependencies(cmd.Context())
		},
	}

	return c
}
