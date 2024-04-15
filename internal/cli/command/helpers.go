package command

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
)

type spinnerOptions struct {
	Spinner          []string
	Suffix           string
	CompletedMessage string
	CanceledMessage  string
}

// func removeDuplicateStr(strs []string) []string {
// 	sort.Strings(strs)
// 	for i := len(strs) - 1; i > 0; i-- {
// 		if strs[i] == strs[i-1] {
// 			strs = append(strs[:i], strs[i+1:]...)
// 		}
// 	}

// 	return strs
// }

func findFilePathForExt(dir string, ext string) (string, error) {
	// Get a list of all files in the current directory.

	files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
	if err != nil {
		return "", fmt.Errorf("failed to list files in current directory: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in directory")
	}

	return files[0], nil
}

func runWithSpinner(ctx context.Context, f func(context.Context) error, options *spinnerOptions) error {
	if options == nil {
		options = &spinnerOptions{}
	}

	if options.Spinner == nil {
		options.Spinner = spinner.CharSets[11]
	}

	if options.Suffix == "" {
		options.Suffix = "Loading..."
	}

	if options.CompletedMessage == "" {
		options.CompletedMessage = "Complete!"
	}

	if options.CanceledMessage == "" {
		options.CanceledMessage = "Canceled!"
	}

	// create a new spinner with the given set of spinner characters
	s := spinner.New(options.Spinner, 100*time.Millisecond,
		spinner.WithHiddenCursor(true),
		spinner.WithSuffix(" "+options.Suffix),
	)

	// Start the spinner
	s.Start()

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		err := f(ctx)
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.Stop()
		fmt.Println(options.CanceledMessage)
		return ctx.Err()
	case err := <-errChan:
		s.Stop()
		if err != nil {
			return err
		}
		fmt.Println(options.CompletedMessage)
		return nil
	}
}
