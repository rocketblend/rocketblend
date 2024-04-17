package command

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	renderProjectOpts struct {
		FrameStart int
		FrameEnd   int
		FrameStep  int

		Output string
		Format string
		commandOpts
	}
)

// newRenderCommand creates a new cobra command for rendering the project.
func newRenderCommand(opts commandOpts) *cobra.Command {
	var frameStart int
	var frameEnd int
	var frameStep int

	var output string
	var format string

	cc := &cobra.Command{
		Use:   "render",
		Short: "Renders the project",
		Long:  `Renders the project from the specified start frame to the end frame, with the given step. Outputs the render in the provided format.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				if err := renderProject(ctx, renderProjectOpts{
					commandOpts: opts,
					FrameStart:  frameStart,
					FrameEnd:    frameEnd,
					FrameStep:   frameStep,
					Output:      output,
					Format:      format,
				}); err != nil {
					return fmt.Errorf("failed to render project: %w", err)
				}

				return nil
			}, &spinnerOptions{
				Suffix:  "Rendering project...",
				Verbose: opts.Global.Verbose,
			})
		},
	}

	cc.Flags().IntVarP(&frameStart, "frame-start", "s", 0, "start frame")
	cc.Flags().IntVarP(&frameEnd, "frame-end", "e", 1, "end frame")
	cc.Flags().IntVarP(&frameStep, "frame-step", "t", 1, "frame step")

	cc.Flags().StringVarP(&output, "output", "o", "//output/{{.Name}}-#####", "set the render path and file name")
	cc.Flags().StringVarP(&format, "format", "f", "PNG", "set the render format")

	return cc
}

func renderProject(ctx context.Context, opts renderProjectOpts) error {
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

	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{opts.Global.WorkingDirectory},
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

	if err := blender.Render(ctx, &types.RenderOpts{
		Start:  opts.FrameStart,
		End:    opts.FrameEnd,
		Step:   opts.FrameStep,
		Output: opts.Output,
		Format: types.RenderFormat(opts.Format),
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Name:         helpers.ExtractName(blendFilePath),
				Path:         blendFilePath,
				Dependencies: resolve.Installations[0],
			},
			Background: true,
		},
	}); err != nil {
		return err
	}

	return nil
}
