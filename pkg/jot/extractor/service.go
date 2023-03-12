package extractor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
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
	defer cancel()

	// Create a channel to signal completion of the progress bar
	done := make(chan struct{})

	// Create a goroutine to display progress
	go func(ctx context.Context, done chan<- struct{}) {
		bar := progressbar.NewOptions(-1,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetDescription("Extracting"),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionClearOnFinish(),
		)

		for {
			select {
			case <-time.After(100 * time.Millisecond):
				bar.Add(1)
			case <-ctx.Done():
				bar.Finish()
				done <- struct{}{}
				return
			}
		}
	}(ctx, done)

	// mholt/archiver doesn't support .dmg files, so we need to handle them separately.
	// This isn't a 100% golang solution, but it works for now.
	var err error
	switch strings.ToLower(filepath.Ext(path)) {
	case ".dmg":
		err = extractDMG(path, extractPath)
	default:
		err = archiver.Unarchive(path, extractPath)
	}
	if err != nil {
		return err
	}

	if s.cleanup {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	// If we get here, the file has been cleaned up, so we can complete the progress bar goroutine
	cancel()
	<-done

	return nil
}
