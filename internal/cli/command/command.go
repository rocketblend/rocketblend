package command

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/container"
	"github.com/rocketblend/rocketblend/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	global struct {
		WorkingDirectory string
		Verbose          bool
		Level            string
	}

	commandOpts struct {
		AppName     string
		Development bool
		Global      *global
	}

	containerOpts struct {
		AppName     string
		Development bool
		Level       string
		Verbose     bool
	}

	RootCommandOpts struct {
		Name    string
		Version string
	}
)

func NewRootCommand(opts *RootCommandOpts) *cobra.Command {
	global := global{}
	commandOpts := commandOpts{
		AppName:     opts.Name,
		Development: false,
		Global:      &global,
	}

	if opts.Version == "dev" {
		commandOpts.Development = true

	}

	cc := &cobra.Command{
		Version: opts.Version,
		Use:     opts.Name,
		Short:   "An improved command-line for Blender.",
		Long: `An improved command-line for Blender. Manage versions,
add-ons, renders and more.

Documentation is available at https://docs.rocketblend.io/`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			path, err := validatePath(global.WorkingDirectory)
			if err != nil {
				return err
			}

			global.WorkingDirectory = path

			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cc.SetVersionTemplate("{{.Version}}\n")

	cc.AddCommand(
		newConfigCommand(commandOpts),
		newNewCommand(commandOpts),
		newInstallCommand(commandOpts),
		newUninstallCommand(commandOpts),
		newRunCommand(commandOpts),
		newRenderCommand(commandOpts),
		newResolveCommand(commandOpts),
		newDescribeCommand(commandOpts),
		newInsertCommand(commandOpts),
	)

	cc.PersistentFlags().StringVarP(&global.WorkingDirectory, "directory", "d", ".", "working directory for the command")
	cc.PersistentFlags().BoolVarP(&global.Verbose, "verbose", "v", false, "enable verbose logging")
	cc.PersistentFlags().StringVarP(&global.Level, "log-level", "l", "info", "log level (debug, info, warn, error)")

	return cc
}

func getContainer(opts containerOpts) (types.Container, error) {
	container, err := container.New(
		container.WithLogger(getLogger(opts.Level, opts.Verbose)),
		container.WithApplicationName(opts.AppName),
		container.WithDevelopmentMode(opts.Development),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	return container, nil
}

func getLogger(level string, verbose bool) types.Logger {
	if verbose {
		return logger.New(
			logger.WithLogLevel(level),
			logger.WithWriters(logger.PrettyWriter()),
		)
	}

	return logger.NoOp()
}

// validatePath checks if the path is valid and returns the absolute path.
func validatePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	// get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// askForConfirmation asks the user for confirmation. If autoConfirm is true, it will return true.
func askForConfirmation(ctx context.Context, prompt string, autoConfirm bool) bool {
	if autoConfirm {
		return true
	}

	fmt.Print(prompt + " (yes/no): ")

	responseChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			responseChan <- "error"
			return
		}
		responseChan <- response
	}()

	select {
	case <-ctx.Done():
		return false
	case response := <-responseChan:
		response = strings.TrimSpace(response)
		return response == "y" || response == "yes"
	}
}
