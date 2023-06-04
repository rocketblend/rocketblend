package extractor

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/mholt/archiver/v3"
)

type (
	Extractor interface {
		Extract(path string, extractPath string) error
		ExtractWithContext(ctx context.Context, path string, extractPath string) error
	}

	Options struct {
		Cleanup bool
		Logger  logger.Logger
	}

	Option func(*Options)

	extractor struct {
		cleanup bool
		logger  logger.Logger
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithCleanup() Option {
	return func(o *Options) {
		o.Cleanup = true
	}
}

func New(opts ...Option) Extractor {
	options := &Options{
		Cleanup: false,
		Logger:  logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	options.Logger.Debug("Initializing extractor", map[string]interface{}{"cleanup": options.Cleanup})

	return &extractor{
		cleanup: options.Cleanup,
		logger:  options.Logger,
	}
}

func (e *extractor) Extract(path string, extractPath string) error {
	return e.ExtractWithContext(context.Background(), path, extractPath)
}

func (e *extractor) ExtractWithContext(ctx context.Context, path string, extractPath string) error {
	logContext := map[string]interface{}{
		"path":        path,
		"extractPath": extractPath,
	}

	e.logger.Info("Extracting", logContext)

	// mholt/archiver doesn't support .dmg files, so we need to handle them separately.
	// This isn't a 100% golang solution, but it works for now.
	var err error
	switch strings.ToLower(filepath.Ext(path)) {
	case ".dmg":
		e.logger.Debug("Extracting DMG file", logContext)
		err = e.extractDMGWithContext(ctx, path, extractPath)
	default:
		e.logger.Debug("Extracting archive", logContext)
		err = archiver.Unarchive(path, extractPath)
	}
	if err != nil {
		logContext["error"] = err.Error()
		e.logger.Error("Extraction error", logContext)
		return err
	}

	if e.cleanup {
		e.logger.Debug("Cleaning up source file", logContext)
		err = os.Remove(path)
		if err != nil {
			logContext["error"] = err.Error()
			e.logger.Error("Cleanup error", logContext)
			return err
		}
	}

	e.logger.Debug("Extraction complete", logContext)

	return nil
}
