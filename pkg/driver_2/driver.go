package driver2

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Options struct {
		Logger    logger.Logger
		Validator types.Validator

		Project *types.Project // Make list?

		Repoistory types.Repository
		Blender    types.Blender
	}

	Option func(*Options)

	driver struct {
		logger    logger.Logger
		validator types.Validator

		project *types.Project

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

func WithProject(project *types.Project) Option {
	return func(o *Options) {
		o.Project = project
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

	if options.Project == nil {
		return nil, errors.New("project is required")
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
		project:    options.Project,
		repoistory: options.Repoistory,
		blender:    options.Blender,
	}, nil
}

func (d *driver) Render(ctx context.Context, opts *types.RenderOpts) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	blendFile, err := d.Resolve(ctx)
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
	if err := ctx.Err(); err != nil {
		return err
	}

	blendFile, err := d.ResolveBlendFile(ctx)
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
	if err := ctx.Err(); err != nil {
		return err
	}

	installations, err := d.Get(ctx, false)
	if err != nil {
		return err
	}

	blendFile := d.resolveBlendFile(installations)
	if err := d.blender.CreateBlendFile(ctx, blendFile); err != nil {
		return err
	}

	if err := d.save(ctx); err != nil {
		return err
	}

	return nil
}

func (d *driver) AddDependencies(ctx context.Context, opts *types.AddDependenciesOpts) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// This will also include the dependencies of the dependencies
	packs, err := d.repoistory.Get(ctx, forceUpdate, references...)
	if err != nil {
		return err
	}

	// Add dependencies to blend config using passed in references
	for _, ref := range references {
		pack := packs[ref]
		if pack.IsBuild() {
			d.blendConfig.RocketFile.SetBuild(ref)
		}

		if pack.IsAddon() {
			d.blendConfig.RocketFile.AddAddons(ref)
		}
	}

	// Install new dependencies
	_, err = d.installationService.Get(ctx, packs, false)
	if err != nil {
		return err
	}

	// Save blend config
	if err = d.save(ctx); err != nil {
		return err
	}

	return nil
}

func (d *driver) RemoveDependencies(ctx context.Context, opts *types.RemoveDependenciesOpts) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	packs, err := d.repoistory.Get(ctx, false, references...)
	if err != nil {
		return fmt.Errorf("failed to get rocket packs: %w", err)
	}

	for index, pack := range packs {
		if pack.IsBuild() {
			d.blendConfig.RocketFile.SetBuild("")
		}

		if pack.IsAddon() {
			d.blendConfig.RocketFile.RemoveAddons(index)
		}
	}

	if err = d.save(ctx); err != nil {
		return fmt.Errorf("failed to save blend config: %w", err)
	}

	return nil
}

func (d *driver) InstallDependencies(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := d.get(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

func (d *driver) Resolve(ctx context.Context) (*types.BlendFile, error) {
	return d.resolve(ctx)
}

func (d *driver) resolve(ctx context.Context, project *types.Project) (*types.BlendFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	installations, err := d.get(ctx, project, false)
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

func (d *driver) getInstallations(ctx context.Context, references []reference.Reference, fetch bool) (map[reference.Reference]*types.Installation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	result, err := d.repoistory.GetInstallations(ctx, &types.GetInstallationsOpts{
		References: references,
		Fetch:      fetch,
	})
	if err != nil {
		return nil, err
	}

	return result.Installations, nil
}

func (d *driver) getDependencies(ctx context.Context) (map[reference.Reference]*types.RocketPack, error) {
	return d.rocketPackService.Get(ctx, false, d.blendConfig.RocketFile.GetDependencies()...)
}

func (d *driver) save(ctx context.Context) error {
	return rocketfile.Save(filepath.Join(d.blendConfig.ProjectPath, rocketfile.FileName), d.blendConfig.RocketFile)
}
