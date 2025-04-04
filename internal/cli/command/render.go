package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rocketblend/rocketblend/internal/cli/ui"
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
		Engine        string

		Output string
		Format string

		EventChan chan types.BlenderEvent
		commandOpts
	}

	displayRenderProjectOpts struct {
		Verbose bool
		renderProjectOpts
	}
)

// newRenderCommand creates a new cobra command for rendering the project.
func newRenderCommand(opts commandOpts) *cobra.Command {
	var frameStart int
	var frameEnd int
	var frameStep int

	var engine string

	var revision int
	var continueRendering bool

	var output string
	var format string

	var autoConfirm bool

	cc := &cobra.Command{
		Use:   "render",
		Short: "Renders the project",
		Long:  `Renders the project using the specified frame range and options.`,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if frameStart < 1 {
				return fmt.Errorf("frame start should be greater than 0")
			}

			if frameEnd == 0 {
				frameEnd = frameStart
			}

			if frameEnd < 0 {
				return fmt.Errorf("frame end should be greater than or equal to 0")
			}

			if frameStep < 1 {
				return fmt.Errorf("frame step should be greater than 0")
			}

			if frameEnd < frameStart {
				return fmt.Errorf("frame end should be greater than or equal to frame start")
			}

			if revision < 0 {
				return fmt.Errorf("revision should be greater than or equal to 0")
			}

			if continueRendering && frameStart == frameEnd {
				return fmt.Errorf("frame start and end should be different when continuing a render")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			blendFilePath, err := findFilePathForExt(opts.Global.WorkingDirectory, types.BlendFileExtension)
			if err != nil {
				return fmt.Errorf("failed to find blend file: %w", err)
			}

			// TODO: Switch to standard relative path formatting and just convert to // for Blender.
			templatePath := strings.Replace(output, "//", fmt.Sprintf("%s/", opts.Global.WorkingDirectory), 1)
			revision := calculateRevision(revision, templatePath, continueRendering)

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

			if frameStart <= existingFrame {
				promptMessage := fmt.Sprintf("The output directory already contains existing frames within the specified range (%d-%d). Are you sure you want to overwrite them?", frameStart, frameEnd)
				if !askForConfirmation(
					cmd.Context(),
					promptMessage,
					autoConfirm,
				) {
					return nil
				}
			}

			return displayRenderProject(cmd.Context(), displayRenderProjectOpts{
				Verbose: opts.Global.Verbose,
				renderProjectOpts: renderProjectOpts{
					BlendFilePath: blendFilePath,
					FrameStart:    frameStart,
					FrameEnd:      frameEnd,
					FrameStep:     frameStep,
					Engine:        engine,
					Output:        outputPath,
					Format:        format,
					commandOpts:   opts,
				},
			})
		},
	}

	cc.Flags().IntVarP(&frameStart, "start", "s", 1, "frame to start rendering from")
	cc.Flags().IntVarP(&frameEnd, "end", "e", 0, "frame to end rendering at, 0 for single frame")
	cc.Flags().IntVarP(&frameStep, "jump", "j", 1, "number of frames to step forward after each rendered frame")

	cc.Flags().IntVarP(&revision, "revision", "r", 0, "revision number for the output directory, 0 for auto-increment")
	cc.Flags().BoolVarP(&continueRendering, "continue", "c", false, "continue rendering from the last rendered frame in the output directory")

	cc.Flags().StringVarP(&engine, "engine", "g", "", "override render engine (cycles, eevee, workbench)")

	cc.Flags().StringVarP(&output, "output", "o", DefaultOutputTemplate, "output path for the rendered frames")
	cc.Flags().StringVarP(&format, "format", "f", "PNG", "output format for the rendered frames")

	cc.Flags().BoolVarP(&autoConfirm, "auto-confirm", "y", false, "overwrite any existing files without requiring confirmation")

	return cc
}

func displayRenderProject(ctx context.Context, opts displayRenderProjectOpts) error {
	if opts.Verbose {
		return renderInVerboseMode(ctx, opts)
	}

	return renderWithUI(ctx, opts)
}

func renderInVerboseMode(ctx context.Context, opts displayRenderProjectOpts) error {
	return renderProject(ctx, renderProjectOpts{
		commandOpts:   opts.commandOpts,
		BlendFilePath: opts.renderProjectOpts.BlendFilePath,
		FrameStart:    opts.renderProjectOpts.FrameStart,
		FrameEnd:      opts.renderProjectOpts.FrameEnd,
		FrameStep:     opts.renderProjectOpts.FrameStep,
		Engine:        opts.renderProjectOpts.Engine,
		Output:        opts.renderProjectOpts.Output,
		Format:        opts.renderProjectOpts.Format,
		EventChan:     nil,
	})
}

func renderWithUI(ctx context.Context, opts displayRenderProjectOpts) error {
	eventChan := make(chan types.BlenderEvent, 100)

	ctxRender, cancelRender := context.WithCancel(ctx)
	defer cancelRender()

	go func() {
		defer close(eventChan)

		if err := renderProject(ctxRender, renderProjectOpts{
			commandOpts:   opts.commandOpts,
			BlendFilePath: opts.renderProjectOpts.BlendFilePath,
			FrameStart:    opts.renderProjectOpts.FrameStart,
			FrameEnd:      opts.renderProjectOpts.FrameEnd,
			FrameStep:     opts.renderProjectOpts.FrameStep,
			Engine:        opts.renderProjectOpts.Engine,
			Output:        opts.renderProjectOpts.Output,
			Format:        opts.renderProjectOpts.Format,
			EventChan:     eventChan,
		}); err != nil {
			if ctxRender.Err() == context.Canceled {
				return
			}

			// Send error to UI
			eventChan <- &types.ErrorEvent{Message: err.Error()}
		}
	}()

	totalFrames := calculateTotalFrames(opts.renderProjectOpts.FrameStart, opts.renderProjectOpts.FrameEnd, opts.renderProjectOpts.FrameStep)

	m := ui.NewRenderProgressModel(totalFrames, eventChan, cancelRender)
	program := tea.NewProgram(&m, tea.WithContext(ctx))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}

	return nil
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
		Format: opts.Format,
		Engine: types.RenderEngine(opts.Engine),
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Path:         opts.BlendFilePath,
				Dependencies: resolve.Installations[0],
				Strict:       profiles.Profiles[0].Strict,
			},
			Background: true,
			EventChan:  opts.EventChan,
		},
	}); err != nil {
		return err
	}

	return nil
}

func calculateTotalFrames(frameStart, frameEnd, frameStep int) int {
	if frameStep <= 0 {
		frameStep = 1
	}

	return (frameEnd-frameStart)/frameStep + 1
}

func calculateRevision(revision int, templatePath string, continueRendering bool) int {
	if revision >= 1 {
		return revision
	}

	current := currentRevision(templatePath)
	if current <= 0 {
		return 1
	}

	if continueRendering {
		return current
	}

	return current + 1
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
