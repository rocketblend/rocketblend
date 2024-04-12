package driver2

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/types"
)

const (
	MaxDependencyDepth = 10
)

type (
	Options struct {
		Logger    logger.Logger
		Validator types.Validator

		Project *types.Project // Make list?

		Repository types.Repository
		Blender    types.Blender
	}

	Option func(*Options)

	driver struct {
		logger    logger.Logger
		validator types.Validator

		project *types.Project

		repository types.Repository
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

func WithProject(project *types.Project) Option {
	return func(o *Options) {
		o.Project = project
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
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Validator == nil {
		return nil, errors.New("validator is required")
	}

	if options.Project == nil {
		return nil, errors.New("project is required")
	}

	if options.Repository == nil {
		return nil, errors.New("repository is required")
	}

	if options.Blender == nil {
		return nil, errors.New("blender is required")
	}

	return &driver{
		logger:     options.Logger,
		validator:  options.Validator,
		project:    options.Project,
		repository: options.Repository,
		blender:    options.Blender,
	}, nil
}

func (d *driver) Render(ctx context.Context, opts *types.RenderOpts) error {
	if err := d.validator.Validate(opts); err != nil {
		return err
	}

	blendFile, err := d.resolve(ctx, d.project)
	if err != nil {
		return err
	}

	if err := d.blender.RenderBlendFile(ctx, &types.RenderBlendFileOpts{
		BlendFile:  blendFile,
		RenderOpts: *opts,
	}); err != nil {
		return err
	}

	return nil
}

func (d *driver) Run(ctx context.Context, opts *types.RunOpts) error {
	if err := d.validator.Validate(opts); err != nil {
		return err
	}

	blendFile, err := d.resolve(ctx, d.project)
	if err != nil {
		return err
	}

	if err := d.blender.RunBlendFile(ctx, &types.RunBlendFileOpts{
		BlendFile: blendFile,
		RunOpts:   *opts,
	}); err != nil {
		return err
	}

	return nil
}

func (d *driver) Create(ctx context.Context) error {
	blendFile, err := d.resolve(ctx, d.project)
	if err != nil {
		return err
	}

	if err := d.blender.CreateBlendFile(ctx, blendFile); err != nil {
		return err
	}

	if err := d.save(ctx, d.project); err != nil {
		return err
	}

	return nil
}

func (d *driver) AddDependencies(ctx context.Context, opts *types.AddDependenciesOpts) error {
	if err := d.validator.Validate(opts); err != nil {
		return err
	}

	if err := d.addDependencies(ctx, d.project, opts.References); err != nil {
		return err
	}

	return nil
}

func (d *driver) RemoveDependencies(ctx context.Context, opts *types.RemoveDependenciesOpts) error {
	if err := d.validator.Validate(opts); err != nil {
		return err
	}

	if err := d.removeDependencies(ctx, d.project, opts.References); err != nil {
		return err
	}

	return nil
}

func (d *driver) InstallDependencies(ctx context.Context) error {
	if err := d.installDependencies(ctx, d.project); err != nil {
		return err
	}

	return nil
}

func (d *driver) Tidy(ctx context.Context) error {
	if err := d.tidy(ctx, d.project); err != nil {
		return err
	}

	return nil
}

func (d *driver) Resolve(ctx context.Context) (*types.BlendFile, error) {
	blendFile, err := d.resolve(ctx, d.project)
	if err != nil {
		return nil, err
	}

	return blendFile, nil
}

func (d *driver) addDependencies(ctx context.Context, project *types.Project, references []reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dependencies := project.RocketFile.Dependencies.Direct
	for _, ref := range references {
		dependencies = append(dependencies, &types.Dependency{
			Reference: ref,
		})
	}

	tidied, err := d.tidyDependencies(ctx, dependencies)
	if err != nil {
		return err
	}

	project.RocketFile.Dependencies = tidied
	if err = d.save(ctx, project); err != nil {
		return err
	}

	return nil
}

func (d *driver) removeDependencies(ctx context.Context, project *types.Project, references []reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dependencies := project.RocketFile.Dependencies.Direct
	for _, ref := range references {
		for index, dep := range dependencies {
			if dep.Reference == ref {
				dependencies = append(dependencies[:index], dependencies[index+1:]...)
				break
			}
		}
	}

	tidied, err := d.tidyDependencies(ctx, dependencies)
	if err != nil {
		return err
	}

	project.RocketFile.Dependencies = tidied
	if err = d.save(ctx, project); err != nil {
		return err
	}

	return nil
}

func (d *driver) tidy(ctx context.Context, project *types.Project) error {
	tidied, err := d.tidyDependencies(ctx, project.RocketFile.Dependencies.Direct)
	if err != nil {
		return err
	}

	project.RocketFile.Dependencies = tidied
	if err = d.save(ctx, project); err != nil {
		return err
	}

	return nil
}

func (d *driver) tidyDependencies(ctx context.Context, dependencies []*types.Dependency) (*types.Dependencies, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	references := make([]reference.Reference, 0, len(dependencies))
	for _, dep := range dependencies {
		references = append(references, dep.Reference)
	}

	result, err := d.repository.GetPackages(ctx, &types.GetPackagesOpts{
		References: references,
		Depth:      MaxDependencyDepth,
	})
	if err != nil {
		return nil, err
	}

	direct := make([]*types.Dependency, len(dependencies))
	indirect := make([]*types.Dependency, len(result.Packs)-len(dependencies))
	for ref := range result.Packs {
		found := false
		for _, dep := range dependencies {
			if dep.Reference == ref {
				found = true
				direct = append(direct, &types.Dependency{
					Reference: ref,
					Type:      result.Packs[ref].Type,
				})
				break
			}
		}

		if !found {
			indirect = append(indirect, &types.Dependency{
				Reference: ref,
				Type:      result.Packs[ref].Type,
			})
		}
	}

	return &types.Dependencies{
		Direct:   direct,
		Indirect: indirect,
	}, nil
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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
		ARGS:         project.RocketFile.ARGS,
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

	if err := helpers.Save(d.validator, filepath.Join(project.Dir(), types.RocketFileName), project); err != nil {
		return err
	}

	return nil
}
