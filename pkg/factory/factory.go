package factory

import (
	"github.com/flowshot-io/x/pkg/logger"
)

type (
	Options struct {
		Logger logger.Logger
	}

	Option func(*Options)

	factory struct {
		logger logger.Logger
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func New(opts ...Option) (*factory, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return &factory{
		logger: options.Logger,
	}, nil
}

func (f *factory) GetLogger() (logger.Logger, error) {
	return f.logger, nil
}
