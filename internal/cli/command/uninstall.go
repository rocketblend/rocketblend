package command

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
			return uninstallWithUI(cmd.Context(), uninstallPackageOpts{
				commandOpts: opts,
				Reference:   args[0],
			})
		},
	}

	return cc
}

// uninstallWithUI runs uninstallPackage asynchronously and shows a Bubble Tea UI.
func uninstallWithUI(ctx context.Context, opts uninstallPackageOpts) error {
	if opts.Global.Verbose {
		return uninstallPackage(ctx, opts)
	}

	eventChan := make(chan ui.ProgressEvent, 10)
	opts.ProgressChan = eventChan
	uninstallCtx, cancelUninstall := context.WithCancel(ctx)
	defer cancelUninstall()

	go func() {
		defer close(eventChan)
		if err := uninstallPackage(uninstallCtx, opts); err != nil {
			eventChan <- ui.ErrorEvent{Message: err.Error()}
		}
	}()

	m := ui.NewProgressModel(eventChan, cancelUninstall)
	program := tea.NewProgram(&m, tea.WithContext(ctx))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}

	return nil
}

// uninstallPackage performs the uninstallation steps and sends events after each step.
func uninstallPackage(ctx context.Context, opts uninstallPackageOpts) error {
	emit := func(ev ui.ProgressEvent) {
		if opts.ProgressChan != nil {
			opts.ProgressChan <- ev
		}
	}

	emit(ui.StepEvent{Message: "Initializing..."})
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
