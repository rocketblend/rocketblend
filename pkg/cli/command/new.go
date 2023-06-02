package command

import (
	"github.com/rocketblend/rocketblend/pkg/rocketblend/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketfile"
	"github.com/spf13/cobra"
)

// newNewCommand creates a new cobra.Command object initialized for creating a new project.
// It expects a single argument which is the name of the project.
// It uses the 'skip-install' flag to decide whether or not to install dependencies.
func (srv *Service) newNewCommand() *cobra.Command {
	var skipInstall bool

	c := &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new project",
		Long:  `Creates a new project with a specified name. Use the 'skip-install' flag to skip installing dependencies.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := srv.getConfig()
			if err != nil {
				return err
			}

			blendConfig, err := blendconfig.New(srv.flags.workingDirectory, args[0], rocketfile.New(config.DefaultBuild))
			if err != nil {
				return err
			}

			driver, err := srv.factory.CreateDriver(blendConfig)
			if err != nil {
				return err
			}

			return driver.Create(cmd.Context())
		},
	}

	c.Flags().BoolVarP(&skipInstall, "skip-install", "s", false, "skip installing dependencies")

	return c
}
