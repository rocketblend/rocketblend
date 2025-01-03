package command

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type uninstallPackageOpts struct {
	commandOpts
	Reference string
}

// newUninstallCommand creates a new cobra.Command that uninstalls dependencies from the current project.
func newUninstallCommand(opts commandOpts) *cobra.Command {
	cc := &cobra.Command{
		Use:   "uninstall [reference]",
		Short: "Remove project dependencies",
		Long:  "Removes dependencies from the current project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				if err := uninstallPackage(ctx, uninstallPackageOpts{
					commandOpts: opts,
					Reference:   args[0],
				}); err != nil {
					return fmt.Errorf("failed to uninstall package: %w", err)
				}

				return nil
			}, &spinnerOptions{Suffix: "Removing package..."})
		},
	}

	return cc
}

func uninstallPackage(ctx context.Context, opts uninstallPackageOpts) error {
	ref, err := reference.Parse(opts.Reference)
	if err != nil {
		return err
	}

	container, err := getContainer(containerOpts{
		AppName:     opts.AppName,
		Development: opts.Development,
		Level:       opts.Global.Level,
		Verbose:     opts.Global.Verbose,
	})
	if err != nil {
		return err
	}

	driver, err := container.GetDriver()
	if err != nil {
		return err
	}

	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{opts.Global.WorkingDirectory},
	})
	if err != nil {
		return err
	}

	profiles.Profiles[0].RemoveDependencies(&types.Dependency{
		Reference: ref,
	})

	if err := driver.SaveProfiles(ctx, &types.SaveProfilesOpts{
		Profiles: map[string]*types.Profile{
			opts.Global.WorkingDirectory: profiles.Profiles[0],
		},
		Overwrite: true,
	}); err != nil {
		return err
	}

	return nil
}
