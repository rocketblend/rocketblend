package blender

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/driver/blenderparser"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	preArgumentOpts struct {
		Background    bool
		BlendFilePath string
	}

	postArgumentOpts struct {
		Addons []*types.Installation
		Script string
	}

	argumentOpts struct {
		PreArgumentOpts    *preArgumentOpts
		PostArgumentOpts   *postArgumentOpts
		IgnoreDependencies bool
	}

	executeOpts struct {
		Executable string
		Arguments  *argumentOpts
	}
)

func (a *preArgumentOpts) ARGS() []string {
	args := []string{}
	if a.Background {
		args = append(args, "-b")
	}

	if a.BlendFilePath != "" {
		args = append(args, a.BlendFilePath)
	}

	return args
}

func (a *postArgumentOpts) ARGS() []string {
	args := []string{}
	if a.Addons != nil {
		json, err := json.Marshal(a.Addons)
		if err != nil {
			return nil
		}

		args = append(args, []string{
			"--python-expr",
			a.Script,
			"--",
			"-a",
			string(json),
		}...)
	}

	return args
}

func (a *argumentOpts) ARGS() []string {
	args := []string{}
	if a.PreArgumentOpts != nil {
		args = append(args, a.PreArgumentOpts.ARGS()...)
	}

	if a.PostArgumentOpts != nil {
		args = append(args, a.PostArgumentOpts.ARGS()...)
	}

	return args
}

func (b *blender) Run(ctx context.Context, opts *types.RunOpts) error {
	if err := b.validator.Validate(opts); err != nil {
		return err
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

func (b *blender) execute(ctx context.Context, opts *executeOpts) error {

	return nil
}

func (b *blender) getRuntimeArguments(opts *getRuntimeCommandOpts) ([]string, error) {
	preArgs := []string{}
	if opts.Background {
		preArgs = append(preArgs, "-b")
	}

	if opts.BlendFile.Path != "" {
		preArgs = append(preArgs, []string{opts.BlendFile.Path}...)
	}

	if opts.ModifyAddons {
		addons := opts.BlendFile.FindAll(types.PackageAddon)
		json, err := json.Marshal(addons)
		if err != nil {
			return nil, err
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

func (b *blender) runCommand(ctx context.Context, name string, args ...string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

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
					s.logger.Debug("blender", map[string]interface{}{"Message": strings.TrimSpace(scanner.Text())})
				}
			} else {
				s.logger.Info("rendering", map[string]interface{}{"Frame": info.FrameNumber, "Memory": info.Memory, "PeakMemory": info.PeakMemory, "Time": info.Time, "Operation": info.Operation})
			}
		}
	}()

	s.logger.Info("running command", map[string]interface{}{"command": cmd.String()})

	if err := cmd.Start(); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		return s.logAndReturnError("failed to start command", err)
	}

	if err := cmd.Wait(); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		return s.logAndReturnError("failed to wait for command", err)
	}

	return nil
}
