package blender

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"

	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	preArguments struct {
		background    bool
		blendFilePath string
	}

	postArguments struct {
		render *renderArguments
		addons []*types.Installation
		script string
	}

	arguments struct {
		// TOOD: Do I need these separate structs?
		preArguments  *preArguments
		postArguments *postArguments
	}

	executable struct {
		executable string
		arguments  *arguments
		output     chan string
	}
)

func (a *preArguments) ARGS() []string {
	args := []string{}
	if a.background {
		args = append(args, "-b")
	}

	if a.blendFilePath != "" {
		args = append(args, a.blendFilePath)
	}

	return args
}

func (a *postArguments) ARGS() []string {
	if a.script == "" {
		return nil
	}

	args := []string{
		"--python-expr",
		a.script,
	}

	if a.addons != nil {
		json, err := json.Marshal(a.addons)
		if err != nil {
			return nil
		}

		args = append(args, []string{
			"--",
			"-a",
			string(json),
		}...)
	}

	return args
}

func (a *arguments) ARGS() []string {
	args := []string{}
	if a.preArguments != nil {
		args = append(args, a.preArguments.ARGS()...)
	}

	if a.postArguments != nil {
		args = append(args, a.postArguments.ARGS()...)
	}

	return args
}

func (e *executable) ARGS() []string {
	args := []string{}
	if e.arguments != nil {
		args = append(args, e.arguments.ARGS()...)
	}

	return args
}

func (e *executable) Name() string {
	return e.executable
}

func (e *executable) OutputChannel() chan string {
	return e.output
}

func (b *blender) execute(ctx context.Context, name string, arguments *arguments, outputChannel chan string) error {
	b.logger.Info("executing", map[string]interface{}{
		"executable": name,
		"arguments":  arguments,
	})

	if err := Execute(ctx, &executable{
		executable: name,
		arguments:  arguments,
		output:     outputChannel,
	}); err != nil {
		return err
	}

	b.logger.Debug("execution completed")

	return nil
}

// Execute runs the given executable with output sent to the executable's output channel.
func Execute(ctx context.Context, executable types.Executable) error {
	cmd := exec.CommandContext(ctx, executable.Name(), executable.ARGS()...)

	outputChannel := executable.OutputChannel()
	if outputChannel != nil {
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		defer cmdReader.Close()

		cmd.Stderr = cmd.Stdout

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				outputChannel <- scanner.Text()
			}
			close(outputChannel)
		}()
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

// ProcessChannel reads strings from a channel and applies a processing function to each string.
func ProcessChannel(outputChannel <-chan string, processFunc func(string)) {
	for data := range outputChannel {
		processFunc(data)
	}
}
