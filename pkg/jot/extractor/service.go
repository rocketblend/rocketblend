package extractor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
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
