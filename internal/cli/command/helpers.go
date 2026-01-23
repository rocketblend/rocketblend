package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/internal/cli/ui"
	"github.com/rocketblend/rocketblend/pkg/types"
)

// findFilePathForExt searches the specified directory for a file with the given extension
func findFilePathForExt(dir string, ext string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		return "", fmt.Errorf("failed to list files in current directory: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in directory")
	}

	return files[0], nil
}

// runWithProgressUI is a helper that runs the provided work function with a progress UI
// when verbose is false, otherwise it runs the work function directly.
func runWithProgressUI(ctx context.Context, verbose bool, work func(ctx context.Context, eventChan chan<- ui.ProgressEvent) error) error {
	if verbose {
		return work(ctx, nil)
	}

	return ui.Run(ctx, work)
}

// displayJSON returns a pretty-printed JSON representation of the provided value.
func displayJSON(v any) (string, error) {
	display, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}

	return string(display), nil
}

// runWithConfirmation runs the provided work function. If it returns an ErrUntrustedDomain error,
// it prompts the user for confirmation before retrying with bypass enabled.
func runWithConfirmation(ctx context.Context, verbose bool, work func(ctx context.Context, eventChan chan<- ui.ProgressEvent, bypass bool) error) error {
	err := runWithProgressUI(ctx, verbose, func(ctx context.Context, eventChan chan<- ui.ProgressEvent) error {
		return work(ctx, eventChan, false)
	})

	if err != nil {
		if ctx.Err() != nil {
			return err
		}

		if errors.Is(err, types.ErrUntrustedDomain) {
			confirmed, confirmErr := ui.Confirm("Warning: One or more dependencies require downloading from untrusted domains.")
			if confirmErr != nil {
				return confirmErr
			}
			if !confirmed {
				return nil
			}

			return runWithProgressUI(ctx, verbose, func(ctx context.Context, eventChan chan<- ui.ProgressEvent) error {
				return work(ctx, eventChan, true)
			})
		}
	}

	return err
}
