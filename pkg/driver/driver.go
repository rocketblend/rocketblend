package driver

import (
	"context"
	"errors"
	"sync"

	"github.com/rocketblend/rocketblend/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Options struct {
		Logger    types.Logger
		Validator types.Validator

		MaxConcurrency int
		ExecutionMode  taskrunner.ExecutionMode

		Repository types.Repository
		Blender    types.Blender
	}

	Option func(*Options)

	Driver struct {
		logger    types.Logger
		validator types.Validator

		maxConcurrency int
		executionMode  taskrunner.ExecutionMode

		repository types.Repository

		mutex sync.Mutex
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

func WithRepository(repository types.Repository) Option {
	return func(o *Options) {
		o.Repository = repository
	}
}

func WithExecutionMode(mode taskrunner.ExecutionMode, maxConcurrency int) Option {
	return func(o *Options) {
		o.ExecutionMode = mode
		o.MaxConcurrency = maxConcurrency
	}
}

func New(opts ...Option) (*Driver, error) {
	options := &Options{
		Logger:         logger.NoOp(),
		ExecutionMode:  taskrunner.Concurrent,
		MaxConcurrency: 5,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Validator == nil {
		return nil, errors.New("validator is required")
	}

	if options.Repository == nil {
		return nil, errors.New("repository is required")
	}

	return &Driver{
		logger:     options.Logger,
		validator:  options.Validator,
		repository: options.Repository,
	}, nil
}

func (d *Driver) getInstallations(ctx context.Context, dependencies []*types.Dependency, fetch bool) (map[reference.Reference]*types.Installation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	result, err := d.repository.GetInstallations(ctx, &types.GetInstallationsOpts{
		Dependencies: dependencies,
		Fetch:        fetch,
	})
	if err != nil {
		return nil, err
	}

	return result.Installations, nil
}
