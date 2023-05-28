package command

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
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
			if err := srv.validateName(args[0]); err != nil {
				return fmt.Errorf("validation failed for name '%s': %w", args[0], err)
			}

			ref, err := srv.parseReference(srv.config.GetValueByString("defaultBuild"))
			if err != nil {
				return fmt.Errorf("failed to parse default build reference: %w", err)
			}

			if err := srv.createProject(cmd.Context(), args[0], ref, skipInstall); err != nil {
				return fmt.Errorf("failed to create project '%s': %w", args[0], err)
			}

			return nil
		},
	}

	c.Flags().BoolVarP(&skipInstall, "skip-install", "s", false, "skip installing dependencies")

	return c
}

// createProject uses the driver to create a new project.
func (srv *Service) createProject(ctx context.Context, name string, buildRef *reference.Reference, skipInstall bool) error {
	if err := srv.driver.Create(ctx, name, srv.flags.workingDirectory, *buildRef, skipInstall); err != nil {
		return fmt.Errorf("driver failed to create project: %w", err)
	}

	return nil
}

// validateName checks if the provided name is valid.
func (srv *Service) validateName(name string) error {
	if filepath.IsAbs(name) || strings.Contains(name, string(filepath.Separator)) {
		return fmt.Errorf("%q is not a valid project name, it should not contain any path separators", name)
	}

	if ext := filepath.Ext(name); ext != "" {
		return fmt.Errorf("%q is not a valid project name, it should not contain any file extension", name)
	}

	return nil
}
