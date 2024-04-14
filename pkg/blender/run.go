package blender

import (
	"context"
	"errors"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	runArguments struct {
		background    bool
		blendFilePath string
		script        string
	}
)

func (b *blender) Run(ctx context.Context, opts *types.RunOpts) error {
	if err := b.validator.Validate(opts); err != nil {
		return err
	}

	build := opts.BlendFile.Build()
	if build == nil {
		return errors.New("missing build")
	}

	addons := opts.BlendFile.Addons()
	if !opts.ModifyAddons {
		addons = nil
	}

	outputChannel := make(chan string, 100)
	defer close(outputChannel)

	go ProcessChannel(outputChannel, b.processOuput)

	if err := b.execute(ctx, build.Path, &arguments{
		preArguments: &preArguments{
			background:    opts.Background,
			blendFilePath: opts.BlendFile.Path,
		},
		postArguments: &postArguments{
			addons: addons,
			script: startupScript(),
		},
	}, outputChannel); err != nil {
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
