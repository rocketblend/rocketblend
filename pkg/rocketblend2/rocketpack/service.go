package rocketpack

import (
	"context"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

type (
	Service interface {
		GetPackages(ctx context.Context, references ...reference.Reference) ([]*RocketPack, error)
		RemovePackages(ctx context.Context, references ...reference.Reference) error
	}

	Options struct {
		Logger logger.Logger
	}

	Option func(*Options)

	service struct {
		logger logger.Logger
	}
)

func NewService(opts ...Option) Service {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return &service{
		logger: options.Logger,
	}
}

func (s *service) GetPackages(ctx context.Context, references ...reference.Reference) ([]*RocketPack, error) {
	return nil, nil
}

func (s *service) RemovePackages(ctx context.Context, references ...reference.Reference) error {
	return nil
}
