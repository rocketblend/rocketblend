package blender

import (
	"context"

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

func (b *blender) Render(ctx context.Context, opts *types.RenderOpts) error {
	return nil
}
