package extractor

import (
	"os"

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
	opts := Options{}

	// if options are passed in, use those
	if options != nil {
		opts = *options

		// if no cleanup is provided, create a default
		opts.Cleanup = true
	}

	return &Service{
		cleanup: opts.Cleanup,
	}
}

func (s *Service) Extract(path string, extractPath string) error {
	err := archiver.Unarchive(path, extractPath)
	if err != nil {
		return err
	}

	if s.cleanup {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
