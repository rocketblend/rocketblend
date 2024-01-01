package blendfile

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/renderoptions"
)

type (
	outputTemplate struct {
		Project string
	}
)

func (s *service) Render(blendFile *BlendFile, opts ...renderoptions.Option) error {
	return s.RenderWithContext(context.Background(), blendFile, opts...)
}

func (s *service) RenderWithContext(ctx context.Context, blendFile *BlendFile, opts ...renderoptions.Option) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	options := &renderoptions.Options{
		Background: true,
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := options.Validate(); err != nil {
		return s.logAndReturnError("invalid options", err)
	}

	output, err := parseOutputTemplate(options.Output, outputTemplate{
		Project: blendFile.ProjectName,
	})
	if err != nil {
		return s.logAndReturnError("failed to parse output template", err)
	}

	args, err := s.getRuntimeArguments(blendFile, options.Background, []string{
		"--frame-start", fmt.Sprint(options.FrameStart),
		"--frame-end", fmt.Sprint(options.FrameEnd),
		"--frame-jump", fmt.Sprint(options.FrameStep),
		"--render-output", output,
		"--render-format", options.Format,
		"-x", "1",
		"-a",
	}...)

	if err != nil {
		return s.logAndReturnError("failed to get runtime arguments", err)
	}

	if err := s.runCommand(ctx, blendFile.Build.FilePath, args...); err != nil {
		return s.logAndReturnError("error running command", err)
	}

	return nil
}
