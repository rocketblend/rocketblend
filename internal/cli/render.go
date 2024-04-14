package cli

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/renderoptions"
	"github.com/spf13/cobra"
)

// newRenderCommand creates a new cobra command for rendering the project.
// It sets up all necessary flags and executes the rendering through the driver.
func (c *cli) newRenderCommand() *cobra.Command {
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
			rocketblend, err := c.getDriver()
			if err != nil {
				return err
			}

			return c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				return rocketblend.Render(
					ctx,
					renderoptions.WithFrameRange(frameStart, frameEnd, frameStep),
					renderoptions.WithOutput(output),
					renderoptions.WithFormat(format),
				)
			}, &spinnerOptions{Suffix: "Rendering project..."})
		},
	}

	cc.Flags().IntVarP(&frameStart, "frame-start", "s", 0, "start frame")
	cc.Flags().IntVarP(&frameEnd, "frame-end", "e", 0, "end frame")
	cc.Flags().IntVarP(&frameStep, "frame-step", "t", 1, "frame step")

	cc.Flags().StringVarP(&output, "output", "o", "//output/{{.Project}}-#####", "set the render path and file name")
	cc.Flags().StringVarP(&format, "format", "f", "PNG", "set the render format")

	return cc
}
