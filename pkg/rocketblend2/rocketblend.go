package rocketblend2

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendfile/renderoptions"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendfile/runoptions"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/installation"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketpack"
)

type (
	Driver interface {
		Render(ctx context.Context, opts ...renderoptions.Option) error
		Run(ctx context.Context, opts ...runoptions.Option) error
		Create(ctx context.Context) error

		InstallDependencies(ctx context.Context) error

		AddDependencies(ctx context.Context, references ...reference.Reference) error
		RemoveDependencies(ctx context.Context, references ...reference.Reference) error

		ResolveBlendFile(ctx context.Context) (*blendfile.BlendFile, error)
	}

	Options struct {
		logger      logger.Logger
		blendConfig *blendconfig.BlendConfig

		installationService installation.Service
		rocketPackService   rocketpack.Service
		blendFileService    blendfile.Service
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

func WithInstallationService(installationService installation.Service) Option {
	return func(o *Options) {
		o.installationService = installationService
	}
}

func WithRocketPackService(rocketPackService rocketpack.Service) Option {
	return func(o *Options) {
		o.rocketPackService = rocketPackService
	}
}

func WithBlendFileService(blendFileService blendfile.Service) Option {
	return func(o *Options) {
		o.blendFileService = blendFileService
	}
}

func New(opts ...Option) (Driver, error) {
	options := &Options{
		logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.installationService == nil {
		isrv, err := installation.NewService()
		if err != nil {
			return nil, fmt.Errorf("failed to create default installation service: %w", err)
		}

		options.installationService = isrv
	}

	if options.rocketPackService == nil {
		rpsrv, err := rocketpack.NewService()
		if err != nil {
			return nil, fmt.Errorf("failed to create default rocket pack service: %w", err)
		}

		options.rocketPackService = rpsrv
	}

	if options.blendFileService == nil {
		bfsrv, err := blendfile.NewService()
		if err != nil {
			return nil, fmt.Errorf("failed to create default blend file service: %w", err)
		}

		options.blendFileService = bfsrv
	}

	if options.blendConfig == nil {
		return nil, fmt.Errorf("blend config is required")
	}

	if err := blendconfig.Validate(options.blendConfig); err != nil {
		return nil, fmt.Errorf("invalid blend config: %w", err)
	}

	return &driver{
		logger:              options.logger,
		InstallationService: options.installationService,
		rocketPackService:   options.rocketPackService,
		blendFileService:    options.blendFileService,
		blendConfig:         options.blendConfig,
	}, nil
}

func (d *driver) Render(ctx context.Context, opts ...renderoptions.Option) error {
	blendFile, err := d.ResolveBlendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.Render(ctx, blendFile, opts...); err != nil {
		return fmt.Errorf("failed to render blend file: %w", err)
	}

	return nil
}

func (d *driver) Run(ctx context.Context, opts ...runoptions.Option) error {
	blendFile, err := d.ResolveBlendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.Run(ctx, blendFile, opts...); err != nil {
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
		if pack.IsBuild() {
			d.blendConfig.RocketFile.SetBuild(index)
		}

		if pack.IsAddon() {
			d.blendConfig.RocketFile.AddAddons(index)
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

	name := filepath.Base(d.blendConfig.BlendFilePath())
	blendFile := &blendfile.BlendFile{
		ProjectName: strings.TrimSuffix(name, filepath.Ext(name)),
		FilePath:    d.blendConfig.BlendFilePath(),
		ARGS:        d.blendConfig.RocketFile.GetArgs(),
	}

	for _, installation := range installations {
		if installation.IsBuild() {
			blendFile.Build = installation.Build
		}

		if installation.IsAddon() {
			blendFile.Addons = append(blendFile.Addons, installation.Addon)
		}
	}

	if err := blendfile.Validate(blendFile); err != nil {
		return nil, fmt.Errorf("invalid blend file: %w", err)
	}

	return blendFile, nil
}

func (d *driver) getDependencies(ctx context.Context) (map[reference.Reference]*rocketpack.RocketPack, error) {
	return d.rocketPackService.GetPackages(ctx, d.blendConfig.RocketFile.GetDependencies()...)
}

func (d *driver) save(ctx context.Context) error {
	return rocketfile.Save(filepath.Join(d.blendConfig.ProjectPath, rocketfile.FileName), d.blendConfig.RocketFile)
}
