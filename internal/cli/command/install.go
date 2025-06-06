package command

import (
	"context"

	"github.com/rocketblend/rocketblend/internal/cli/ui"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type installPackageOpts struct {
	commandOpts
	Reference    string
	Pull         bool
	ProgressChan chan<- ui.ProgressEvent
}

// newInstallCommand creates a new cobra command for installing project dependencies.
func newInstallCommand(opts commandOpts) *cobra.Command {
	var update bool

	cc := &cobra.Command{
		Use:   "install [reference]",
		Short: "Installs project dependencies",
		Long:  `Adds the specified dependencies to the current project and installs them. If no reference is provided, all dependencies in the project are installed instead.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := ""
			if len(args) > 0 {
				ref = args[0]
			}

			return runWithProgressUI(
				cmd.Context(),
				opts.Global.Verbose,
				func(ctx context.Context, eventChan chan<- ui.ProgressEvent) error {
					return installPackage(ctx, installPackageOpts{
						commandOpts:  opts,
						Reference:    ref,
						Pull:         update,
						ProgressChan: eventChan,
					})
				})
		},
	}

	cc.Flags().BoolVarP(&update, "update", "u", false, "updates to the latest package definitions before installing")
	return cc
}

// installPackage performs the installation steps and sends events after each step.
func installPackage(ctx context.Context, opts installPackageOpts) error {
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

	if opts.Reference != "" {
		emit(ui.StepEvent{Message: "Updating dependencies..."})
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

		profiles.Profiles[0].AddDependencies(&types.Dependency{
			Reference: ref,
		})
	}

	emit(ui.StepEvent{Message: "Tidying profiles..."})
	if err := driver.TidyProfiles(ctx, &types.TidyProfilesOpts{
		Profiles: profiles.Profiles,
		Fetch:    opts.Pull,
	}); err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Installing dependencies..."})
	if err := driver.InstallProfiles(ctx, &types.InstallProfilesOpts{
		Profiles: profiles.Profiles,
	}); err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Saving profiles..."})
	if err := driver.SaveProfiles(ctx, &types.SaveProfilesOpts{
		Profiles: map[string]*types.Profile{
			opts.Global.WorkingDirectory: profiles.Profiles[0],
		},
		Overwrite: true,
	}); err != nil {
		return err
	}

	emit(ui.CompletionEvent{Message: "Dependencies installed!"})
	return nil
}
