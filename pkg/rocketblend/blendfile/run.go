package blendfile

import (
	"bufio"
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/rocketblend/blendfile/runoptions"
)

func (s *service) Run(ctx context.Context, blendFile *BlendFile, opts ...runoptions.Option) error {
	options := &runoptions.Options{
		Background: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := options.Validate(); err != nil {
		return err
	}

	cmd, err := s.getCommand(ctx, blendFile, options.Background)
	if err != nil {
		return fmt.Errorf("failed to get command: %w", err)
	}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			s.logger.Debug("Blender", map[string]interface{}{"Message": scanner.Text()})
		}
	}()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for command: %w", err)
	}

	return nil
}
