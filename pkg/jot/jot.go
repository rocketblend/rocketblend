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
		mutex      sync.Mutex
		mutexes    map[string]*sync.Mutex
		storageDir string
		logger     logger.Logger
		downloader downloader.Downloader
		extractor  extractor.Extractor
	}

	Options struct {
		StorageDir string
		Logger     logger.Logger
		Downloader downloader.Downloader
		Extractor  extractor.Extractor
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

func WithStorageDir(storageDir string) Option {
	return func(o *Options) {
		o.StorageDir = storageDir
	}
}

func New(opts ...Option) (Storage, error) {
	options := &Options{
		Logger:     logger.NoOp(),
		Extractor:  extractor.New(nil),
		Downloader: downloader.New(),
	}

	for _, opt := range opts {
		opt(options)
	}

	err := os.MkdirAll(options.StorageDir, 0755)
	if err != nil {
		return nil, err
	}

	// create driver
	return &Driver{
		storageDir: filepath.Clean(options.StorageDir),
		mutexes:    make(map[string]*sync.Mutex),
		logger:     options.Logger,
		downloader: options.Downloader,
		extractor:  options.Extractor,
	}, nil
}
