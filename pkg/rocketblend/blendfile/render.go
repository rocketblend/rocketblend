package blendfile

import (
	"bufio"
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/rocketblend/blenderparser"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/blendfile/renderoptions"
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
		return fmt.Errorf("invalid render options: %w", err)
	}

	output, err := parseOutputTemplate(options.Output, outputTemplate{
		Project: blendFile.ProjectName,
	})
	if err != nil {
		return fmt.Errorf("failed to parse output template: %w", err)
	}

	args := []string{
		"--frame-start",
		fmt.Sprint(options.FrameStart),
		"--frame-end",
		fmt.Sprint(options.FrameEnd),
		"--frame-jump",
		fmt.Sprint(options.FrameStep),
		"--render-output",
		output,
		"--render-format",
		options.Format,
		"-x", // Set option to add the file extension to the end of the file.
		"1",
		"-a", // Render frames from start to end
	}

	cmd, err := s.getCommand(ctx, blendFile, options.Background, args...)
	if err != nil {
		return fmt.Errorf("failed to get command: %w", err)
	}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			info, err := blenderparser.ParseRenderOutput(scanner.Text())
			if err != nil {
				s.logger.Debug("Blender", map[string]interface{}{"Message": scanner.Text()})
				continue
			} else {
				s.logger.Info("Rendering", map[string]interface{}{"Frame": info.FrameNumber, "Memory": info.Memory, "PeakMemory": info.PeakMemory, "Time": info.Time, "Operation": info.Operation})
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("waiting for command: %w", err)
	}

	return nil
}
