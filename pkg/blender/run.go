package blender

import (
	"context"
	"errors"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

func (b *blender) Run(ctx context.Context, opts *types.RunOpts) error {
	if err := b.validator.Validate(opts); err != nil {
		return err
	}

	builds := opts.BlendFile.FindAll(types.PackageBuild)
	if len(builds) == 0 {
		return errors.New("no builds found")
	}

	outputChannel := make(chan string, 100)
	defer close(outputChannel)

	go ProcessChannel(outputChannel, b.processOuput)

	if err := b.execute(ctx, builds[0].Path, &arguments{}, outputChannel); err != nil {
		return err
	}

	return nil
}

// processOuput parses the output from the blender process and logs the relevant information.
func (b *blender) processOuput(output string) {
	if output == "" {
		return
	}

	// TODO: Just parse the output here and put on a channel to be consumed by the caller.
	info, err := parseRenderOutput(output)
	if err != nil {
		b.logger.Debug("blender", map[string]interface{}{
			"message": strings.TrimSpace(output),
		})

		return
	}

	b.logger.Info("rendering", map[string]interface{}{
		"frame":      info.FrameNumber,
		"memory":     info.Memory,
		"peakMemory": info.PeakMemory,
		"time":       info.Time,
		"operation":  info.Operation,
	})
}
