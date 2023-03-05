package extractor

import (
	"context"
	"fmt"
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

	// Create a goroutine to display progress
	go displayProgress(ctx)

	// Defer the cancel function
	defer cancel()

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

	return nil
}

func displayProgress(ctx context.Context) {
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Extracting"),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(7),
	)

	for {
		select {
		case <-time.After(500 * time.Millisecond):
			bar.Add(1)
		case <-ctx.Done():
			fmt.Println("Goroutine cancelled")
			return
		}
	}
}
