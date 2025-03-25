package command

import (
	"context"

	"github.com/rocketblend/rocketblend/internal/cli/ui"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type uninstallPackageOpts struct {
	commandOpts
	Reference    string
	ProgressChan chan<- ui.ProgressEvent
}

// newUninstallCommand creates a new cobra.Command that uninstalls dependencies from the current project.
func newUninstallCommand(opts commandOpts) *cobra.Command {
	cc := &cobra.Command{
		Use:   "uninstall [reference]",
		Short: "Remove project dependencies",
		Long:  "Removes dependencies from the current project.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithProgressUI(
				cmd.Context(),
				opts.Global.Verbose,
				func(ctx context.Context, eventChan chan<- ui.ProgressEvent) error {
					return uninstallPackage(ctx, uninstallPackageOpts{
						commandOpts:  opts,
						Reference:    args[0],
						ProgressChan: eventChan,
					})
				})
		},
	}

	return cc
}

// uninstallPackage performs the uninstallation steps and sends events after each step.
func uninstallPackage(ctx context.Context, opts uninstallPackageOpts) error {
	emit := func(ev ui.ProgressEvent) {
		if opts.ProgressChan != nil {
			opts.ProgressChan <- ev
		}
	}

	emit(ui.StepEvent{Message: "Initialising..."})
	container, err := getContainer(containerOpts{
		AppName:     opts.AppName,
		Development: opts.Development,
		Level:       opts.Global.Level,
		Verbose:     opts.Global.Verbose,
	})
	if err != nil {
		return err
	}

	configurator, err := container.GetConfigurator()
	if err != nil {
		return err
	}

	config, err := configurator.Get()
	if err != nil {
		return err
	}

	ref, err := reference.Aliased(opts.Reference, config.Aliases)
	if err != nil {
		return err
	}

	driver, err := container.GetDriver()
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Loading profiles..."})
	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{opts.Global.WorkingDirectory},
	})
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Removing dependency..."})
	profiles.Profiles[0].RemoveDependencies(&types.Dependency{
		Reference: ref,
	})

	emit(ui.StepEvent{Message: "Saving profiles..."})
	if err := driver.SaveProfiles(ctx, &types.SaveProfilesOpts{
		Profiles: map[string]*types.Profile{
			opts.Global.WorkingDirectory: profiles.Profiles[0],
		},
		Overwrite: true,
	}); err != nil {
		return err
	}

	emit(ui.CompletionEvent{Message: "Dependency removed!"})
	return nil
}
