package installation

import (
	"context"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketpack"
)

type (
	Service interface {
		GetInstallations(ctx context.Context, rocketPacks []*rocketpack.RocketPack, readOnly bool) ([]*Installation, error)
		RemoveInstallations(ctx context.Context, rocketPacks []*rocketpack.RocketPack) error
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

func (s *service) GetInstallations(ctx context.Context, rocketPacks []*rocketpack.RocketPack, readOnly bool) ([]*Installation, error) {
	return nil, nil
}

func (s *service) RemoveInstallations(ctx context.Context, rocketPacks []*rocketpack.RocketPack) error {
	return nil
}
