package blender

import (
	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/rocketblend/rocketblend/pkg/validator"
)

type (
	Options struct {
		Logger    types.Logger
		Validator types.Validator
	}

	Option func(*Options)

	Blender struct {
		logger    types.Logger
		validator types.Validator
	}
)

func WithLogger(logger types.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithValidator(validator types.Validator) Option {
	return func(o *Options) {
		o.Validator = validator
	}
}

func New(opts ...Option) (*Blender, error) {
	options := &Options{
		Logger:    logger.NoOp(),
		Validator: validator.New(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return &Blender{
		logger:    options.Logger,
		validator: options.Validator,
	}, nil
}
