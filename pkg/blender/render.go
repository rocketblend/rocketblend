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
			Format:  RenderFormat(opts.Format),
			Threads: opts.Threads,
			Engine:  convertRenderEngine(opts.Engine),
		},
	}

	b.logger.Info("rendering", map[string]interface{}{
		"blendFile": opts.BlendFile.Path,
		"output":    opts.Output,
		"start":     opts.Start,
		"end":       opts.End,
		"step":      opts.Step,
		"format":    opts.Format,
		"threads":   opts.Threads,
		"engine":    opts.Engine,
	})

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

func createRenderEvent(b *Blender, info *renderInfo) *types.RenderEvent {
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

	return &types.RenderEvent{
		Frame:      info.FrameNumber,
		Memory:     info.Memory,
		PeakMemory: info.PeakMemory,
		Time:       info.Time,
		Operation:  info.Operation,
		Data:       info.Data,
	}
}
