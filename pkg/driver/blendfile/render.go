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

func (s *service) Render(ctx context.Context, blendFile *BlendFile, opts ...renderoptions.Option) error {
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

	args := []string{
		"--frame-start", fmt.Sprint(options.FrameStart),
		"--frame-end", fmt.Sprint(options.FrameEnd),
		"--frame-jump", fmt.Sprint(options.FrameStep),
		"--render-output", output,
		"--render-format", options.Format,
		"-x", "1",
		"-a",
	}

	cmd, err := s.getCommand(ctx, blendFile, options.Background, args...)
	if err != nil {
		return s.logAndReturnError("error getting command", err)
	}

	if err := s.runCommand(cmd); err != nil {
		return s.logAndReturnError("error running command", err)
	}

	return nil
}
