package extractor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mholt/archiver/v3"
	"github.com/schollz/progressbar/v3"
)

type (
	Service struct {
		cleanup bool
	}

	Options struct {
		Cleanup bool
	}
)

func New(options *Options) *Service {
	// create default options
	opts := Options{
		Cleanup: true,
	}

	// if options are passed in, use those
	if options != nil {
		opts = *options
	}

	return &Service{
		cleanup: opts.Cleanup,
	}
}

func (s *Service) Extract(path string, extractPath string) error {
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	// Create a goroutine to display progress
	go func(ctx context.Context) {
		bar := progressbar.NewOptions(-1,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetDescription("Extracting"),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionClearOnFinish(),
			// progressbar.OptionOnCompletion(func() {
			// 	fmt.Fprint(os.Stderr, "\n")
			// }),
		)

		defer bar.Finish()
		defer wg.Done()

		for {
			select {
			case <-time.After(100 * time.Millisecond):
				bar.Add(1)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	// Defer the cancel function
	defer cancel()
	defer wg.Wait()

	// mholt/archiver doesn't support .dmg files, so we need to handle them separately.
	// This isn't a 100% golang solution, but it works for now.
	switch strings.ToLower(filepath.Ext(path)) {
	case ".dmg":
		err := extractDMG(path, extractPath)
		if err != nil {
			return err
		}
	default:
		err := archiver.Unarchive(path, extractPath)
		if err != nil {
			return err
		}
	}

	if s.cleanup {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	// Call cancel function explicitly and wait for it to complete.
	cancel()
	wg.Wait()

	return nil
}
