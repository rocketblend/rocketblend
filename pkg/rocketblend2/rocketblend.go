package rocketblend2

import (
	"context"
	"fmt"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/installation"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketpack"
)

type (
	RocketPackService interface {
		GetPackages(ctx context.Context, references ...reference.Reference) ([]*rocketpack.RocketPack, error)
		RemovePackages(ctx context.Context, references ...reference.Reference) error
	}

	InstallationService interface {
		GetInstallations(ctx context.Context, references ...reference.Reference) ([]*installation.Installation, error)
		RemoveInstallations(ctx context.Context, references ...reference.Reference) error
	}

	BlendConfigService interface {
		InstallDependencies(ctx context.Context, config *blendconfig.BlendConfig) error
		ResolveBlendFile(ctx context.Context, config *blendconfig.BlendConfig) (*blendfile.BlendFile, error)
	}

	BlendFileService interface {
		Render(ctx context.Context, blendFile *blendfile.BlendFile) error
		Run(ctx context.Context, blendFile *blendfile.BlendFile) error
		Create(ctx context.Context, blendFile *blendfile.BlendFile) error
	}

	Driver interface {
		Render(ctx context.Context) error
		Run(ctx context.Context) error

		InstallDependencies(ctx context.Context) error

		AddDependencies(ctx context.Context, references ...reference.Reference) error
		RemoveDependencies(ctx context.Context, references ...reference.Reference) error
	}

	Options struct {
		logger      logger.Logger
		blendConfig *blendconfig.BlendConfig
	}

	Option func(*Options)

	driver struct {
		logger logger.Logger

		rocketPackService  RocketPackService
		blendConfigService BlendConfigService
		blendFileService   BlendFileService

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
	blendFile, err := d.resolveBlendFile()
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.Render(ctx, blendFile); err != nil {
		return fmt.Errorf("failed to render blend file: %w", err)
	}

	return nil
}

func (d *driver) Run(ctx context.Context) error {
	blendFile, err := d.resolveBlendFile()
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.Run(ctx, blendFile); err != nil {
		return fmt.Errorf("failed to run blend file: %w", err)
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
			d.blendConfig.RocketFile.Build = references[index]
		}

		if packType == rocketpack.TypeAddon {
			d.blendConfig.RocketFile.Addons = append(d.blendConfig.RocketFile.Addons, references[index])
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
			d.blendConfig.RocketFile.Build = ""
		}

		if packType == rocketpack.TypeAddon {
			for i, addon := range d.blendConfig.RocketFile.Addons {
				if addon == references[index] {
					d.blendConfig.RocketFile.Addons = append(d.blendConfig.RocketFile.Addons[:i], d.blendConfig.RocketFile.Addons[i+1:]...)
				}
			}
		}
	}

	if err = d.save(ctx); err != nil {
		return fmt.Errorf("failed to save blend config: %w", err)
	}

	return nil
}

func (d *driver) InstallDependencies(ctx context.Context) error {
	return d.blendConfigService.InstallDependencies(ctx, d.blendConfig)
}

func (d *driver) save(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (d *driver) resolveBlendFile() (*blendfile.BlendFile, error) {
	return d.blendConfigService.ResolveBlendFile(context.Background(), d.blendConfig)
}
