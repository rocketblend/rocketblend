package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return srv.validateProjectName(args[0])
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := srv.getConfig()
			if err != nil {
				return err
			}

			blendConfig, err := blendconfig.New(
				srv.flags.workingDirectory,
				srv.ensureBlendExtension(args[0]),
				rocketfile.New(config.DefaultBuild),
			)
			if err != nil {
				return err
			}

			driver, err := srv.createDriver(blendConfig)
			if err != nil {
				return err
			}

			return driver.Create(cmd.Context())
		},
	}

	c.Flags().BoolVarP(&skipInstall, "skip-install", "s", false, "skip installing dependencies")

	return c
}

// validateProjectName checks if the project name is valid.
func (srv *Service) validateProjectName(projectName string) error {
	if filepath.IsAbs(projectName) || strings.Contains(projectName, string(filepath.Separator)) {
		return fmt.Errorf("%q is not a valid project name, it should not contain any path separators", projectName)
	}

	if ext := filepath.Ext(projectName); ext != "" {
		return fmt.Errorf("%q is not a valid project name, it should not contain any file extension", projectName)
	}

	return nil
}

// ensureBlendExtension adds ".blend" extension to filename if it does not already have it.
func (srv *Service) ensureBlendExtension(filename string) string {
	if !strings.HasSuffix(filename, ".blend") {
		filename += ".blend"
	}

	return filename
}
