package command

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	runProjectOpts struct {
		commandOpts
	}
)

// newRunCommand creates a new cobra command for running the project.
func newRunCommand(opts commandOpts) *cobra.Command {
	cc := &cobra.Command{
		Use:   "run",
		Short: "Runs the project",
		Long:  `Launches the project in the current working directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				if err := runProject(ctx, runProjectOpts{
					commandOpts: opts,
				}); err != nil {
					return fmt.Errorf("failed to run project: %w", err)
				}

				return nil
			}, &spinnerOptions{Suffix: "Running project..."})
		},
	}

	return cc
}

func runProject(ctx context.Context, opts runProjectOpts) error {
	blendFilePath, err := findFilePathForExt(opts.Global.WorkingDirectory, types.BlendFileExtension)
	if err != nil {
		return err
	}

	container, err := getContainer(opts.AppName, opts.Development, opts.Global.Verbose)
	if err != nil {
		return err
	}

	driver, err := container.GetDriver()
	if err != nil {
		return err
	}

	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{filepath.Dir(opts.Global.WorkingDirectory)},
	})
	if err != nil {
		return err
	}

	resolve, err := driver.ResolveProfiles(ctx, &types.ResolveProfilesOpts{
		Profiles: profiles.Profiles,
	})
	if err != nil {
		return err
	}

	blender, err := container.GetBlender()
	if err != nil {
		return err
	}

	if err := blender.Run(ctx, &types.RunOpts{
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Path:         blendFilePath,
				Dependencies: resolve.Installations[0],
			},
		},
	}); err != nil {
		return err
	}

	return nil
}
