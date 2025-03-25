package command

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/internal/cli/ui"
)

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
