package jot

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/jot/downloader"
	"github.com/rocketblend/rocketblend/pkg/jot/extractor"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

type (
	Storage interface {
		Read(reference reference.Reference, resource string) ([]byte, error)
		Write(reference reference.Reference, resource string, downloadUrl string) error
		WriteWithContext(ctx context.Context, reference reference.Reference, resource string, downloadUrl string) error
		DeleteAll(reference reference.Reference) error
	}

	Driver struct {
		mutex       sync.Mutex
		mutexes     map[string]*sync.Mutex
		storagePath string
		logger      logger.Logger
		downloader  downloader.Downloader
		extractor   extractor.Extractor
	}

	Options struct {
		StoragePath string
		Logger      logger.Logger
		Downloader  downloader.Downloader
		Extractor   extractor.Extractor
	}

	Option func(*Options)
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithDownloader(downloader downloader.Downloader) Option {
	return func(o *Options) {
		o.Downloader = downloader
	}
}

func WithExtractor(extractor extractor.Extractor) Option {
	return func(o *Options) {
		o.Extractor = extractor
	}
}

func WithStorageDir(storagePath string) Option {
	return func(o *Options) {
		o.StoragePath = storagePath
	}
}

func New(opts ...Option) (Storage, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Downloader == nil {
		options.Downloader = downloader.New(downloader.WithLogger(options.Logger))
	}

	if options.Extractor == nil {
		options.Extractor = extractor.New(
			extractor.WithLogger(options.Logger),
			extractor.WithCleanup())
	}

	// create storage dir
	err := os.MkdirAll(options.StoragePath, 0755)
	if err != nil {
		return nil, err
	}

	// create driver
	return &Driver{
		storagePath: filepath.Clean(options.StoragePath),
		mutexes:     make(map[string]*sync.Mutex),
		logger:      options.Logger,
		downloader:  options.Downloader,
		extractor:   options.Extractor,
	}, nil
}
