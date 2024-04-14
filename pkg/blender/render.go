package blender

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	outputPathData struct {
		Name string `json:"name"`
	}
)

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

	build := opts.BlendFile.Build()
	if build == nil {
		return errors.New("missing build")
	}

	arguments := arguments{
		background:    opts.Background,
		blendFilePath: opts.BlendFile.Path,
		render: &renderArguments{
			start:   opts.Start,
			end:     opts.End,
			step:    opts.Step,
			output:  outputPath,
			format:  opts.Format,
			devices: opts.CyclesDevices,
			threads: opts.Threads,
		},
	}

	if opts.BlendFile.Addons() != nil {
		arguments.script = startupScript()
		arguments.rockeblend = &rocketblendArguments{
			addons: opts.BlendFile.Addons(),
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
