package command

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	installPackageOpts struct {
		commandOpts
		Reference string
		Update    bool
	}
)

// newInstallCommand creates a new cobra command for installing project dependencies.
func newInstallCommand(opts commandOpts) *cobra.Command {
	//var forceUpdate bool

	cc := &cobra.Command{
		Use:   "install [reference]",
		Short: "Installs project dependencies",
		Long:  `Adds the specified dependencies to the current project and installs them. If no reference is provided, all dependencies in the project are installed instead.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			suffix := "Installing dependencies..."
			if len(args) > 0 {
				suffix = "Installing package..."
			}

			return runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				if err := installPackage(ctx, installPackageOpts{
					commandOpts: opts,
					Reference:   args[0],
					Update:      false,
				}); err != nil {
					return fmt.Errorf("failed to install project dependencies: %w", err)
				}

				return nil
			}, &spinnerOptions{
				Suffix:  suffix,
				Verbose: opts.Global.Verbose,
			})
		},
	}

	//cc.Flags().BoolVarP(&forceUpdate, "update", "u", false, "refreshes the package definition before installing it")

	return cc
}

func installPackage(ctx context.Context, opts installPackageOpts) error {
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

	if opts.Reference != "" {
		ref, err := reference.Parse(opts.Reference)
		if err != nil {
			return err
		}

		profiles.Profiles[0].AddDependencies(&types.Dependency{
			Reference: ref,
		})
	}

	if err := driver.TidyProfiles(ctx, &types.TidyProfilesOpts{
		Profiles: profiles.Profiles,
	}); err != nil {
		return err
	}

	if err := driver.InstallProfiles(ctx, &types.InstallProfilesOpts{
		Profiles: profiles.Profiles,
	}); err != nil {
		return err
	}

	if err := driver.SaveProfiles(ctx, &types.SaveProfilesOpts{
		Profiles: map[string]*types.Profile{
			opts.Global.WorkingDirectory: profiles.Profiles[0],
		},
	}); err != nil {
		return err
	}

	return nil
}
