package driver2

import (
	"errors"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Options struct {
		Logger    logger.Logger
		Validator types.Validator

		Projects []*types.Project

		Repoistory types.Repository
		Blender    types.Blender
	}

	Option func(*Options)

	driver struct {
		logger    logger.Logger
		validator types.Validator

		projects []*types.Project

		repoistory types.Repository
		blender    types.Blender
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

func WithProjects(projects []*types.Project) Option {
	return func(o *Options) {
		o.Projects = projects
	}
}

func WithRepository(repoistory types.Repository) Option {
	return func(o *Options) {
		o.Repoistory = repoistory
	}
}

func WithBlender(blender types.Blender) Option {
	return func(o *Options) {
		o.Blender = blender
	}
}

func New(opts ...Option) (*driver, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Validator == nil {
		return nil, errors.New("validator is required")
	}

	if options.Projects == nil {
		return nil, errors.New("projects is required")
	}

	if options.Repoistory == nil {
		return nil, errors.New("repoistory is required")
	}

	if options.Blender == nil {
		return nil, errors.New("blender is required")
	}

	return &driver{
		logger:     options.Logger,
		validator:  options.Validator,
		projects:   options.Projects,
		repoistory: options.Repoistory,
		blender:    options.Blender,
	}, nil
}
