package blender

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	outputPathData struct {
		Name string `json:"name"`
	}

	renderArguments struct {
		start        int
		end          int
		step         int
		output       string
		format       types.RenderFormat
		cyclesDevice []types.CyclesDevice
		Threads      int
	}
)

func (a *renderArguments) ARGS() []string {
	if a.start == 0 && a.end == 0 {
		return nil
	}

	args := []string{}
	if a.start != 0 {
		args = append(args, "--frame-start", fmt.Sprint(a.start))
	}

	if a.end != 0 {
		args = append(args, "--frame-end", fmt.Sprint(a.end))
	}

	if a.step != 0 {
		args = append(args, "--frame-jump", fmt.Sprint(a.step))
	}

	if a.output != "" {
		args = append(args, "--render-output", a.output)
	}

	if a.format != "" {
		args = append(args, "--render-format", string(a.format), "-x", "1")
	}

	if len(a.cyclesDevice) > 0 {
		devices := []string{"--cycles-device"}
		for _, device := range a.cyclesDevice {
			devices = append(devices, string(device))
		}

		args = append(args, strings.Join(devices, "+"))
	}

	if a.Threads > 0 {
		args = append(args, "-t", strconv.Itoa(a.Threads))
	}

	return append(args, "-a")
}

func (b *blender) Render(ctx context.Context, opts *types.RenderOpts) error {
	if err := b.validator.Validate(opts); err != nil {
		return err
	}

	outputPath, err := helpers.ParseTemplateWithData(opts.Output, &outputPathData{
		Name: opts.BlendFile.Name,
	})
	if err != nil {
		return err
	}
}

func outputPath(template string, data outputTemplateData) (string, error) {
	result, err := helpers.ParseTemplateWithData(template, data)
	if err != nil {
		return "", err
	}

	return result, nil
}
