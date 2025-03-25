package command

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rocketblend/rocketblend/internal/cli/ui"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type insertPackageOpts struct {
	commandOpts
	Reference    string
	ProgressChan chan<- ui.ProgressEvent
}

func newInsertCommand(opts commandOpts) *cobra.Command {
	cc := &cobra.Command{
		Use:   "insert",
		Short: "Inserts a package into your local library",
		Long:  `Inserts a package into your local library.`,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasPrefix(args[0], "local/") {
				return fmt.Errorf("local package reference must start with local/")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			insOpts := insertPackageOpts{
				commandOpts: opts,
				Reference:   args[0],
			}
			if opts.Global.Verbose {
				return insertPackage(cmd.Context(), insOpts)
			}
			return insertWithUI(cmd.Context(), insOpts)
		},
	}

	return cc
}

// insertWithUI runs insertPackage asynchronously and displays progress using the generic progress UI.
func insertWithUI(ctx context.Context, opts insertPackageOpts) error {
	eventChan := make(chan ui.ProgressEvent, 10)
	opts.ProgressChan = eventChan
	insCtx, cancelInsert := context.WithCancel(ctx)
	defer cancelInsert()

	go func() {
		defer close(eventChan)
		if err := insertPackage(insCtx, opts); err != nil {
			eventChan <- ui.ErrorEvent{Message: err.Error()}
		}
	}()

	m := ui.NewProgressModel(eventChan, cancelInsert)
	program := tea.NewProgram(&m, tea.WithContext(ctx))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}

	return nil
}

// insertPackage performs the package insertion and emits progress events for each step.
func insertPackage(ctx context.Context, opts insertPackageOpts) error {
	emit := func(ev ui.ProgressEvent) {
		if opts.ProgressChan != nil {
			opts.ProgressChan <- ev
		}
	}

	emit(ui.StepEvent{Message: "Initialising..."})
	ref, err := reference.Parse(opts.Reference)
	if err != nil {
		return err
	}

	container, err := getContainer(containerOpts{
		AppName:     opts.AppName,
		Development: opts.Development,
		Level:       opts.Global.Level,
		Verbose:     opts.Global.Verbose,
	})
	if err != nil {
		return err
	}

	validator, err := container.GetValidator()
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Loading package file..."})
	packagePath := filepath.Join(opts.Global.WorkingDirectory, types.PackageFileName)
	pack, err := helpers.Load[types.Package](validator, packagePath)
	if err != nil {
		return err
	}

	repository, err := container.GetRepository()
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Inserting package into local library..."})
	if err := repository.InsertPackages(ctx, &types.InsertPackagesOpts{
		Packs: map[reference.Reference]*types.Package{
			ref: pack,
		},
	}); err != nil {
		return err
	}

	emit(ui.CompletionEvent{Message: "Package inserted!"})
	return nil
}
