package blender

import (
	"context"
	"errors"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

func (b *Blender) Run(ctx context.Context, opts *types.RunOpts) error {
	if err := b.validator.Validate(opts); err != nil {
		return err
	}

	build := opts.BlendFile.Build()
	if build == nil {
		return errors.New("missing build")
	}

	arguments := arguments{
		Background:    opts.Background,
		BlendFilePath: opts.BlendFile.Path,
	}

	if opts.BlendFile.Addons() != nil {
		arguments.Script = startupScript()
		arguments.Rockeblend = &rocketblendArguments{
			Addons: opts.BlendFile.Addons(),
			Strict: opts.BlendFile.InjectionMode == types.StrictInjectionMode,
		}
	}

	outputChannel := make(chan string, 100)
	defer close(outputChannel)

	go ProcessChannel(outputChannel, b.processOuput)

	if err := b.execute(ctx, build.Path, &arguments, outputChannel); err != nil {
		return err
	}

	return nil
}

// processOuput parses the output from the blender process and logs the relevant information.
func (b *Blender) processOuput(output string) {
	if output == "" {
		return
	}

	// TODO: Just parse the output here and put on a channel to be consumed by the caller.
	info, err := parseRenderOutput(output)
	if err != nil {
		b.logger.Debug("blender", map[string]interface{}{
			"output": strings.ToLower(strings.TrimSpace(output)),
		})

		return
	}

	data := map[string]interface{}{
		"frame":      info.FrameNumber,
		"memory":     info.Memory,
		"peakMemory": info.PeakMemory,
		"time":       info.Time,
	}

	for key, value := range info.Data {
		data[key] = value
	}

	b.logger.Info(info.Operation, data)
}
