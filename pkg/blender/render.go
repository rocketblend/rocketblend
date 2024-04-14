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
		Background:    opts.Background,
		BlendFilePath: opts.BlendFile.Path,
		Render: &renderArguments{
			Start:   opts.Start,
			End:     opts.End,
			Step:    opts.Step,
			Output:  outputPath,
			Format:  opts.Format,
			Devices: opts.CyclesDevices,
			Threads: opts.Threads,
		},
	}

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
