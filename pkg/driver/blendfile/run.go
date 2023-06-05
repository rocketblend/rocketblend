package blendfile

import (
	"bufio"
	"context"
	"os/exec"

	"github.com/rocketblend/rocketblend/pkg/driver/blenderparser"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/runoptions"
)

func (s *service) Run(ctx context.Context, blendFile *BlendFile, opts ...runoptions.Option) error {
	options := &runoptions.Options{
		Background: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := options.Validate(); err != nil {
		return s.logAndReturnError("invalid run options", err)
	}

	cmd, err := s.getCommand(ctx, blendFile, options.Background)
	if err != nil {
		return s.logAndReturnError("failed to get command", err)
	}

	if err := s.runCommand(ctx, cmd); err != nil {
		return s.logAndReturnError("error running command", err)
	}

	return nil
}

func (s *service) runCommand(ctx context.Context, cmd *exec.Cmd) error {
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return s.logAndReturnError("failed to get stdout pipe", err)
	}
	defer cmdReader.Close()

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			info, err := blenderparser.ParseRenderOutput(scanner.Text())
			if err != nil {
				if scanner.Text() != "" {
					s.logger.Debug("Blender", map[string]interface{}{"Message": scanner.Text()})
				}
			} else {
				s.logger.Info("Rendering", map[string]interface{}{"Frame": info.FrameNumber, "Memory": info.Memory, "PeakMemory": info.PeakMemory, "Time": info.Time, "Operation": info.Operation})
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		return s.logAndReturnError("failed to start command", err)
	}

	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.Canceled {
			s.logger.Debug("Context cancelled - command exited")
			return nil
		}

		return s.logAndReturnError("failed to wait for command", err)
	}

	return nil
}
