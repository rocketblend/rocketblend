package command

import (
	"context"

	"github.com/rocketblend/rocketblend/internal/cli/ui"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type runProjectOpts struct {
	commandOpts
	ProgressChan chan<- ui.ProgressEvent
}

// newRunCommand creates a new cobra command for running the project.
func newRunCommand(opts commandOpts) *cobra.Command {
	cc := &cobra.Command{
		Use:   "run",
		Short: "Runs the project",
		Long:  "Launches the project in the current working directory.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithProgressUI(
				cmd.Context(),
				opts.Global.Verbose,
				func(ctx context.Context, eventChan chan<- ui.ProgressEvent) error {
					return runProject(ctx, runProjectOpts{
						commandOpts:  opts,
						ProgressChan: eventChan,
					})
				})
		},
	}

	return cc
}

// runProject performs the steps needed to run the project and emits progress events.
func runProject(ctx context.Context, opts runProjectOpts) error {
	emit := func(ev ui.ProgressEvent) {
		if opts.ProgressChan != nil {
			opts.ProgressChan <- ev
		}
	}

	emit(ui.StepEvent{Message: "Initialising..."})
	blendFilePath, err := findFilePathForExt(opts.Global.WorkingDirectory, types.BlendFileExtension)
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

	emit(ui.StepEvent{Message: "Loading profiles..."})
	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{opts.Global.WorkingDirectory},
	})
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Resolving dependencies..."})
	resolve, err := driver.ResolveProfiles(ctx, &types.ResolveProfilesOpts{
		Profiles: profiles.Profiles,
	})
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Step: 7, Message: "Running project..."})
	blender, err := container.GetBlender()
	if err != nil {
		return err
	}

	if err := blender.Run(ctx, &types.RunOpts{
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Path:         blendFilePath,
				Dependencies: resolve.Installations[0],
				Strict:       profiles.Profiles[0].Strict,
			},
		},
	}); err != nil {
		return err
	}

	emit(ui.CompletionEvent{Message: "Exited Blender!"})
	return nil
}
