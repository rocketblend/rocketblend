package rocketblend2

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/installation"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketpack"
)

type (
	Driver interface {
		Render(ctx context.Context) error
		Run(ctx context.Context) error
		Create(ctx context.Context) error

		InstallDependencies(ctx context.Context) error

		AddDependencies(ctx context.Context, references ...reference.Reference) error
		RemoveDependencies(ctx context.Context, references ...reference.Reference) error

		ResolveBlendFile(ctx context.Context) (*blendfile.BlendFile, error)
	}

	Options struct {
		logger      logger.Logger
		blendConfig *blendconfig.BlendConfig
	}

	Option func(*Options)

	driver struct {
		logger logger.Logger

		InstallationService installation.Service
		rocketPackService   rocketpack.Service
		blendFileService    blendfile.Service

		blendConfig *blendconfig.BlendConfig
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func WithBlendConfig(blendConfig *blendconfig.BlendConfig) Option {
	return func(o *Options) {
		o.blendConfig = blendConfig
	}
}

func New(opts ...Option) (Driver, error) {
	options := &Options{
		logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.blendConfig == nil {
		return nil, fmt.Errorf("blend config is required")
	}

	if err := blendconfig.Validate(options.blendConfig); err != nil {
		return nil, fmt.Errorf("invalid blend config: %w", err)
	}

	return &driver{
		logger: options.logger,

		blendConfig: options.blendConfig,
	}, nil
}

func (d *driver) Render(ctx context.Context) error {
	blendFile, err := d.ResolveBlendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.Render(ctx, blendFile); err != nil {
		return fmt.Errorf("failed to render blend file: %w", err)
	}

	return nil
}

func (d *driver) Run(ctx context.Context) error {
	blendFile, err := d.ResolveBlendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.Run(ctx, blendFile); err != nil {
		return fmt.Errorf("failed to run blend file: %w", err)
	}

	return nil
}

func (d *driver) Create(ctx context.Context) error {
	blendFile, err := d.ResolveBlendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.Create(ctx, blendFile); err != nil {
		return fmt.Errorf("failed to create blend file: %w", err)
	}

	if err := d.save(ctx); err != nil {
		return fmt.Errorf("failed to save blend config: %w", err)
	}

	return nil
}

func (d *driver) AddDependencies(ctx context.Context, references ...reference.Reference) error {
	packs, err := d.rocketPackService.GetPackages(ctx, references...)
	if err != nil {
		return fmt.Errorf("failed to get rocket packs: %w", err)
	}

	for index, pack := range packs {
		packType, err := pack.GetType()
		if err != nil {
			return fmt.Errorf("failed to get rocket pack type: %w", err)
		}

		if packType == rocketpack.TypeBuild {
			d.blendConfig.RocketFile.SetBuild(references[index])
		}

		if packType == rocketpack.TypeAddon {
			d.blendConfig.RocketFile.AddAddons(references[index])
		}
	}

	if err = d.save(ctx); err != nil {
		return fmt.Errorf("failed to save blend config: %w", err)
	}

	return nil
}

func (d *driver) RemoveDependencies(ctx context.Context, references ...reference.Reference) error {
	packs, err := d.rocketPackService.GetPackages(ctx, references...)
	if err != nil {
		return fmt.Errorf("failed to get rocket packs: %w", err)
	}

	for index, pack := range packs {
		packType, err := pack.GetType()
		if err != nil {
			return fmt.Errorf("failed to get rocket pack type: %w", err)
		}

		if packType == rocketpack.TypeBuild {
			d.blendConfig.RocketFile.SetBuild("")
		}

		if packType == rocketpack.TypeAddon {
			d.blendConfig.RocketFile.RemoveAddons(references[index])
		}
	}

	if err = d.save(ctx); err != nil {
		return fmt.Errorf("failed to save blend config: %w", err)
	}

	return nil
}

func (d *driver) InstallDependencies(ctx context.Context) error {
	packs, err := d.getDependencies(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rocket packs: %w", err)
	}

	_, err = d.InstallationService.GetInstallations(ctx, packs, false)
	if err != nil {
		return fmt.Errorf("failed to get installations: %w", err)
	}

	return nil
}

func (d *driver) ResolveBlendFile(ctx context.Context) (*blendfile.BlendFile, error) {
	packs, err := d.getDependencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rocket packs: %w", err)
	}

	installations, err := d.InstallationService.GetInstallations(ctx, packs, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get installations: %w", err)
	}

	var build *blendfile.Build
	var addons []*blendfile.Addon
	for _, installation := range installations {
		packType, err := installation.RocketPack.GetType()
		if err != nil {
			return nil, fmt.Errorf("failed to get rocket pack type: %w", err)
		}

		if packType == rocketpack.TypeBuild {
			build = &blendfile.Build{
				FilePath: installation.FilePath,
				ARGS:     installation.RocketPack.Build.Args,
			}
		}

		if packType == rocketpack.TypeAddon {
			addons = append(addons, &blendfile.Addon{
				FilePath: installation.FilePath,
				Name:     installation.RocketPack.Addon.Name,
				Version:  installation.RocketPack.Addon.Version,
			})
		}
	}

	blendFile := &blendfile.BlendFile{
		FilePath: d.blendConfig.BlendFilePath(),
		Build:    build,
		Addons:   addons,
		ARGS:     d.blendConfig.RocketFile.GetArgs(),
	}

	if err := blendfile.Validate(blendFile); err != nil {
		return nil, fmt.Errorf("invalid blend file: %w", err)
	}

	return blendFile, nil
}

func (d *driver) getDependencies(ctx context.Context) ([]*rocketpack.RocketPack, error) {
	// TODO: make sure GetPackages returns dependencies from builds as separate packages
	return d.rocketPackService.GetPackages(ctx, d.blendConfig.RocketFile.GetDependencies()...)
}

func (d *driver) save(ctx context.Context) error {
	return rocketfile.Save(filepath.Join(d.blendConfig.ProjectPath, rocketfile.FileName), d.blendConfig.RocketFile)
}
