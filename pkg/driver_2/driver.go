package driver2

import (
	"context"
	"errors"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Options struct {
		Logger    logger.Logger
		Validator types.Validator

		MaxConcurrency int
		ExecutionMode  taskrunner.ExecutionMode

		Projects []*types.Project

		Repository types.Repository
		Blender    types.Blender
	}

	Option func(*Options)

	driver struct {
		logger    logger.Logger
		validator types.Validator

		maxConcurrency int
		executionMode  taskrunner.ExecutionMode

		projects []*types.Project

		repository types.Repository
		blender    types.Blender

		mutex sync.Mutex
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

func WithExecutionMode(mode taskrunner.ExecutionMode, maxConcurrency int) Option {
	return func(o *Options) {
		o.ExecutionMode = mode
		o.MaxConcurrency = maxConcurrency
	}
}

func WithProject(projects ...*types.Project) Option {
	return func(o *Options) {
		o.Projects = projects
	}
}

func WithRepository(repository types.Repository) Option {
	return func(o *Options) {
		o.Repository = repository
	}
}

func WithBlender(blender types.Blender) Option {
	return func(o *Options) {
		o.Blender = blender
	}
}

func New(opts ...Option) (*driver, error) {
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

	if options.Blender == nil {
		return nil, errors.New("blender is required")
	}

	if len(options.Projects) == 0 {
		return nil, errors.New("projects are required")
	}

	return &driver{
		logger:     options.Logger,
		validator:  options.Validator,
		projects:   options.Projects,
		repository: options.Repository,
		blender:    options.Blender,
	}, nil
}

func (d *driver) AddDependencies(ctx context.Context, opts *types.AddDependenciesOpts) error {
	if err := d.validator.Validate(opts); err != nil {
		return err
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[struct{}], len(d.projects))
	for i, project := range d.projects {
		tasks[i] = func(ctx context.Context) (struct{}, error) {
			return struct{}{}, d.addDependencies(ctx, project, opts.References)
		}
	}

	_, err := taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) RemoveDependencies(ctx context.Context, opts *types.RemoveDependenciesOpts) error {
	if err := d.validator.Validate(opts); err != nil {
		return err
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[struct{}], len(d.projects))
	for i, project := range d.projects {
		tasks[i] = func(ctx context.Context) (struct{}, error) {
			return struct{}{}, d.removeDependencies(ctx, project, opts.References)
		}
	}

	_, err := taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) InstallDependencies(ctx context.Context) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[struct{}], len(d.projects))
	for i, project := range d.projects {
		tasks[i] = func(ctx context.Context) (struct{}, error) {
			return struct{}{}, d.installDependencies(ctx, project)
		}
	}

	_, err := taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) Resolve(ctx context.Context) (*types.BlendFile, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	blendFile, err := d.resolve(ctx, d.projects[0]) // TOOD: Handle multiple projects
	if err != nil {
		return nil, err
	}

	return blendFile, nil
}

func (d *driver) Save(ctx context.Context) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[struct{}], len(d.projects))
	for i, project := range d.projects {
		tasks[i] = func(ctx context.Context) (struct{}, error) {
			return struct{}{}, d.save(ctx, project)
		}
	}

	_, err := taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) addDependencies(ctx context.Context, project *types.Project, references []reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	result, err := d.repository.GetPackages(ctx, &types.GetPackagesOpts{
		References: references,
	})
	if err != nil {
		return err
	}

	dependencies := project.Profile.Dependencies
	for ref, pack := range result.Packs {
		dependencies = append(dependencies, &types.Dependency{
			Reference: ref,
			Type:      pack.Type,
		})
	}

	project.Profile.Dependencies = dependencies

	return nil
}

func (d *driver) removeDependencies(ctx context.Context, project *types.Project, references []reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dependencies := project.Profile.Dependencies
	for _, ref := range references {
		for index, dep := range dependencies {
			if dep.Reference == ref {
				dependencies = append(dependencies[:index], dependencies[index+1:]...)
				break
			}
		}
	}

	project.Profile.Dependencies = dependencies

	return nil
}

func (d *driver) installDependencies(ctx context.Context, project *types.Project) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := d.getInstallations(ctx, project.Requires(), true)
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) resolve(ctx context.Context, project *types.Project) (*types.BlendFile, error) {
	installations, err := d.getInstallations(ctx, project.Requires(), false)
	if err != nil {
		return nil, err
	}

	dependencies := make([]*types.Installation, 0, len(installations))
	for _, installation := range installations {
		dependencies = append(dependencies, installation)
	}

	return &types.BlendFile{
		Name:         project.Name(),
		Path:         project.BlendFilePath,
		ARGS:         project.Profile.ARGS,
		Dependencies: dependencies,
	}, nil
}

func (d *driver) getInstallations(ctx context.Context, dependencies []*types.Dependency, fetch bool) (map[reference.Reference]*types.Installation, error) {
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

func (d *driver) save(ctx context.Context, project *types.Project) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := helpers.Save(d.validator, filepath.Join(project.Dir(), types.ProjectConfigFileName), project); err != nil {
		return err
	}

	return nil
}
