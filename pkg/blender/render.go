package blender

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (b *Blender) Render(ctx context.Context, opts *types.RenderOpts) error {
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
		Render: &renderArguments{
			Start:   opts.Start,
			End:     opts.End,
			Step:    opts.Step,
			Output:  opts.Output,
			Format:  opts.Format,
			Devices: opts.CyclesDevices,
			Threads: opts.Threads,
		},
	}

	b.logger.Info("rendering", map[string]interface{}{
		"blendFile": opts.BlendFile.Path,
		"output":    opts.Output,
		"start":     opts.Start,
		"end":       opts.End,
		"step":      opts.Step,
		"format":    opts.Format,
		"devices":   opts.CyclesDevices,
		"threads":   opts.Threads,
	})

	if opts.BlendFile.Addons() != nil {
		arguments.Script = startupScript()
		arguments.Rockeblend = &rocketblendArguments{
			Addons: opts.BlendFile.Addons(),
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
