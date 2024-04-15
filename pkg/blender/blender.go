package blender

import (
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
)

type (
	Options struct {
		Logger    logger.Logger
		Validator types.Validator
	}

	Option func(*Options)

	blender struct {
		logger    logger.Logger
		validator types.Validator
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func New(opts ...Option) (*blender, error) {
	options := &Options{
		Logger:    logger.NoOp(),
		Validator: validator.New(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return &blender{
		logger:    options.Logger,
		validator: options.Validator,
	}, nil
}
