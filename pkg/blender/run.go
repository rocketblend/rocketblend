package blender

import (
	"context"
	"errors"

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

	if opts.BlendFile.Addons() != nil || opts.BlendFile.Strict {
		arguments.Script = startupScript()
		arguments.Rockeblend = &rocketblendArguments{
			Addons: opts.BlendFile.Addons(),
			Strict: opts.BlendFile.Strict,
		}
	}

	outputChan := make(chan string, 100)
	defer close(outputChan)

	go processChannel(outputChan, opts.EventChan, b.processOutput)

	if err := b.execute(ctx, build.Path, &arguments, outputChan); err != nil {
		return err
	}

	return nil
}
