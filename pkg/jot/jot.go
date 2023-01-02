package jot

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/rocketblend/rocketblend/pkg/jot/downloader"
	"github.com/rocketblend/rocketblend/pkg/jot/extractor"
	"github.com/sirupsen/logrus"
)

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	// create default options
	opts := Options{}

	// if options are passed in, use those
	if options != nil {
		opts = *options
	}

	// if no logger is provided, create a default
	if opts.Logger == nil {
		l := logrus.New()
		l.Level = logrus.InfoLevel
		opts.Logger = l
	}

	// if no downloader is provided, create a default
	if opts.Downloader == nil {
		opts.Downloader = downloader.New()
	}

	if opts.Extractor == nil {
		opts.Extractor = extractor.New(nil)
	}

	// create driver
	driver := Driver{
		dir:        dir,
		mutexes:    make(map[string]*sync.Mutex),
		log:        opts.Logger,
		downloader: opts.Downloader,
		extractor:  opts.Extractor,
	}

	// if the database already exists, just use it
	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (store already exists)\n", dir)
		return &driver, nil
	}

	// if the database doesn't exist create it
	opts.Logger.Debug("Creating jot storage at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}
