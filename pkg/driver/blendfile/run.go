package blendfile

import (
	"bufio"
	"context"
	"encoding/json"
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

	args, err := s.getRuntimeArguments(blendFile, options.Background)
	if err != nil {
		return s.logAndReturnError("failed to get runtime arguments", err)
	}

	if err := s.runCommand(ctx, blendFile.Build.FilePath, args...); err != nil {
		return s.logAndReturnError("error running command", err)
	}

	return nil
}

func (s *service) getRuntimeArguments(blendFile *BlendFile, background bool, postArgs ...string) ([]string, error) {
	preArgs := []string{}
	if background {
		preArgs = append(preArgs, "-b")
	}

	if blendFile.FilePath != "" {
		preArgs = append(preArgs, []string{blendFile.FilePath}...)
	}

	if s.addonsEnabled {
		json, err := json.Marshal(blendFile.Addons)
		if err != nil {
			return nil, s.logAndReturnError("failed to marshal addons", err)
		}

		postArgs = append([]string{
			"--python-expr",
			s.addonScript,
		}, postArgs...)

		postArgs = append(postArgs, []string{
			"--",
			"-a",
			string(json),
		}...)
	}

	// Blender requires arguments to be in a specific order
	return append(preArgs, postArgs...), nil
}

func (s *service) runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)

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

	s.logger.Info("Running command", map[string]interface{}{"command": cmd.String()})

	if err := cmd.Start(); err != nil {
		return s.logAndReturnError("failed to start command", err)
	}

	if err := cmd.Wait(); err != nil {
		return s.logAndReturnError("failed to wait for command", err)
	}

	return nil
}
