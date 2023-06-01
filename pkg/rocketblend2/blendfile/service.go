package blendfile

import (
	"context"

	"github.com/flowshot-io/x/pkg/logger"
)

type (
	Service interface {
		Render(ctx context.Context, blendFile *BlendFile) error
		Run(ctx context.Context, blendFile *BlendFile) error
		Create(ctx context.Context, blendFile *BlendFile) error
	}

	Options struct {
		Logger logger.Logger
	}

	Option func(*Options)

	service struct {
		logger logger.Logger
	}
)

func NewService(opts ...Option) (Service, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return &service{
		logger: options.Logger,
	}, nil
}

func (s *service) Render(ctx context.Context, blendFile *BlendFile) error {
	return nil
}

func (s *service) Run(ctx context.Context, blendFile *BlendFile) error {
	return nil
}

func (s *service) Create(ctx context.Context, blendFile *BlendFile) error {
	return nil
}
