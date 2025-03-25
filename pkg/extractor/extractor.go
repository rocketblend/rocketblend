package extractor

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/rocketblend/rocketblend/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Options struct {
		Cleanup bool
		Logger  logger.Logger
	}

	Option func(*Options)

	Extractor struct {
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

func New(opts ...Option) (*Extractor, error) {
	options := &Options{
		Cleanup: false,
		Logger:  logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	options.Logger.Debug("initializing Extractor", map[string]interface{}{"cleanup": options.Cleanup})

	return &Extractor{
		cleanup: options.Cleanup,
		logger:  options.Logger,
	}, nil
}

func (e *Extractor) Extract(ctx context.Context, opts *types.ExtractOpts) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	logContext := map[string]interface{}{
		"path":       opts.Path,
		"outputPath": opts.OutputPath,
	}

	e.logger.Info("extracting", logContext)

	// mholt/archiver doesn't support .dmg files, so we need to handle them separately.
	// This isn't a 100% golang solution, but it works for now.
	var err error
	switch strings.ToLower(filepath.Ext(opts.Path)) {
	case ".dmg":
		e.logger.Debug("extracting DMG file", logContext)
		err = e.extractDMG(ctx, opts.Path, opts.OutputPath)
	default:
		e.logger.Debug("extracting archive", logContext)
		err = archiver.Unarchive(opts.Path, opts.OutputPath)
	}
	if err != nil {
		logContext["error"] = err.Error()
		e.logger.Error("extraction error", logContext)
		return err
	}

	if e.cleanup {
		e.logger.Debug("cleaning up source file", logContext)
		err = os.Remove(opts.Path)
		if err != nil {
			logContext["error"] = err.Error()
			e.logger.Error("cleanup error", logContext)
			return err
		}
	}

	e.logger.Debug("extraction complete", logContext)

	return nil
}
