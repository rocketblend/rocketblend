package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/blender"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

const DefaultOutputTemplate = "//output/" + blender.RevisionTempalteVariable + "/" + blender.NameTemplateVariable + "-#####"

type (
	renderProjectOpts struct {
		BlendFilePath string
		FrameStart    int
		FrameEnd      int
		FrameStep     int

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

	var revision int
	var continueRendering bool

	var output string
	var format string

	var autoConfirm bool

	cc := &cobra.Command{
		Use:   "render",
		Short: "Renders the project",
		Long:  `Renders the project from the specified start frame to the end frame, with the given step. Outputs the render in the provided format.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			blendFilePath, err := findFilePathForExt(opts.Global.WorkingDirectory, types.BlendFileExtension)
			if err != nil {
				return fmt.Errorf("failed to find blend file: %w", err)
			}

			// TODO: Switch to standard relative path formatting and just convert to // for Blender.
			templatePath := strings.Replace(output, "//", fmt.Sprintf("%s/", opts.Global.WorkingDirectory), 1)

			if revision < 1 {
				revision = currentRevision(templatePath) + 1
			}

			outputPath, err := helpers.ParseTemplateWithData(templatePath, &blender.TemplatedOutputData{
				Name:     helpers.ExtractName(blendFilePath),
				Revision: helpers.PadWithZero(revision, 5),
			})
			if err != nil {
				return fmt.Errorf("failed to parse output template: %w", err)
			}

			existingFrame, err := existingFrameNumber(outputPath)
			if err != nil {
				return fmt.Errorf("failed to find existing frame: %w", err)
			}

			if continueRendering && existingFrame > 0 {
				if existingFrame < frameEnd {
					frameStart = existingFrame + 1
				}
			}

			if existingFrame != 0 && existingFrame >= frameStart && (frameEnd == 0 || existingFrame <= frameEnd) {
				promptMessage := fmt.Sprintf("The output directory already contains existing frames within the specified range (%d-%d). Are you sure you want to overwrite them?", frameStart, frameEnd)
				if frameEnd == 0 {
					promptMessage = fmt.Sprintf("The output directory already contains frames starting from %d. Are you sure you want to overwrite them?", frameStart)
				}

				if !askForConfirmation(
					cmd.Context(),
					fmt.Sprintf(promptMessage, frameStart, frameEnd),
					autoConfirm,
				) {
					return nil
				}
			}

			return runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				if err := renderProject(ctx, renderProjectOpts{
					commandOpts:   opts,
					BlendFilePath: blendFilePath,
					FrameStart:    frameStart,
					FrameEnd:      frameEnd,
					FrameStep:     frameStep,
					Output:        outputPath,
					Format:        format,
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

	cc.Flags().IntVarP(&revision, "revision", "r", 0, "revision subfolder for output. Defaults to auto-increment.")
	cc.Flags().BoolVarP(&continueRendering, "continue", "c", false, "continue rendering from the last rendered image within the specified frame range.")

	cc.Flags().StringVarP(&output, "output", "o", DefaultOutputTemplate, "output path for the rendered frames")
	cc.Flags().StringVarP(&format, "format", "f", "PNG", "output format for the rendered frames")

	cc.Flags().BoolVarP(&autoConfirm, "auto-confirm", "y", false, "overwrite existing files without confirmation")

	return cc
}

func renderProject(ctx context.Context, opts renderProjectOpts) error {
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

	blend, err := container.GetBlender()
	if err != nil {
		return err
	}

	if err := blend.Render(ctx, &types.RenderOpts{
		Start:  opts.FrameStart,
		End:    opts.FrameEnd,
		Step:   opts.FrameStep,
		Output: opts.Output,
		Format: types.RenderFormat(opts.Format),
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Path:         opts.BlendFilePath,
				Dependencies: resolve.Installations[0],
			},
			Background: true,
		},
	}); err != nil {
		return err
	}

	return nil
}

func currentRevision(templatedPath string) int {
	revision, err := blender.FindMaxRevision(templatedPath)
	if err != nil {
		return 0
	}

	return revision
}

func existingFrameNumber(path string) (int, error) {
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		return 0, nil
	}

	return blender.FindMaxFrameNumber(path)
}
