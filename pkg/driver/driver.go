package driver

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/renderoptions"
	"github.com/rocketblend/rocketblend/pkg/driver/blendfile/runoptions"
	"github.com/rocketblend/rocketblend/pkg/driver/installation"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
)

type (
	Driver interface {
		Render(ctx context.Context, opts ...renderoptions.Option) error
		Run(ctx context.Context, opts ...runoptions.Option) error
		Create(ctx context.Context) error

		InstallDependencies(ctx context.Context) error

		AddDependencies(ctx context.Context, forceUpdate bool, references ...reference.Reference) error
		RemoveDependencies(ctx context.Context, references ...reference.Reference) error

		ResolveBlendFile(ctx context.Context) (*blendfile.BlendFile, error)
	}

	Options struct {
		Logger      logger.Logger
		BlendConfig *blendconfig.BlendConfig

		InstallationService installation.Service
		RocketPackService   rocketpack.Service
		BlendFileService    blendfile.Service
	}

	Option func(*Options)

	driver struct {
		logger logger.Logger

		installationService installation.Service
		rocketPackService   rocketpack.Service
		blendFileService    blendfile.Service

		blendConfig *blendconfig.BlendConfig
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithBlendConfig(blendConfig *blendconfig.BlendConfig) Option {
	return func(o *Options) {
		o.BlendConfig = blendConfig
	}
}

func WithInstallationService(installationService installation.Service) Option {
	return func(o *Options) {
		o.InstallationService = installationService
	}
}

func WithRocketPackService(rocketPackService rocketpack.Service) Option {
	return func(o *Options) {
		o.RocketPackService = rocketPackService
	}
}

func WithBlendFileService(blendFileService blendfile.Service) Option {
	return func(o *Options) {
		o.BlendFileService = blendFileService
	}
}

func New(opts ...Option) (Driver, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.InstallationService == nil {
		isrv, err := installation.NewService(
			installation.WithLogger(options.Logger),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create default installation service: %w", err)
		}

		options.InstallationService = isrv
	}

	if options.RocketPackService == nil {
		rpsrv, err := rocketpack.NewService(
			rocketpack.WithLogger(options.Logger),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create default rocket pack service: %w", err)
		}

		options.RocketPackService = rpsrv
	}

	if options.BlendFileService == nil {
		bfsrv, err := blendfile.NewService(
			blendfile.WithLogger(options.Logger),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create default blend file service: %w", err)
		}

		options.BlendFileService = bfsrv
	}

	if options.BlendConfig == nil {
		return nil, fmt.Errorf("blend config is required")
	}

	if err := blendconfig.Validate(options.BlendConfig); err != nil {
		return nil, fmt.Errorf("invalid blend config: %w", err)
	}

	options.Logger.Debug("initializing rocketblend driver", map[string]interface{}{
		"ProjectPath":   options.BlendConfig.ProjectPath,
		"BlendFileName": options.BlendConfig.BlendFileName,
	})

	return &driver{
		logger:              options.Logger,
		installationService: options.InstallationService,
		rocketPackService:   options.RocketPackService,
		blendFileService:    options.BlendFileService,
		blendConfig:         options.BlendConfig,
	}, nil
}

func (d *driver) Render(ctx context.Context, opts ...renderoptions.Option) error {
	d.logger.Debug("rendering blend file")

	blendFile, err := d.ResolveBlendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.RenderWithContext(ctx, blendFile, opts...); err != nil {
		return fmt.Errorf("failed to render blend file: %w", err)
	}

	return nil
}

func (d *driver) Run(ctx context.Context, opts ...runoptions.Option) error {
	d.logger.Debug("running blend file")

	blendFile, err := d.ResolveBlendFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.RunWithContext(ctx, blendFile, opts...); err != nil {
		return fmt.Errorf("failed to run blend file: %w", err)
	}

	return nil
}

func (d *driver) Create(ctx context.Context) error {
	d.logger.Debug("creating blend file")

	installations, err := d.Get(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to get installations: %w", err)
	}

	blendFile, err := d.resolveBlendFile(ctx, installations)
	if err != nil {
		return fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := d.blendFileService.CreateWithContext(ctx, blendFile); err != nil {
		return fmt.Errorf("failed to create blend file: %w", err)
	}

	if err := d.save(ctx); err != nil {
		return fmt.Errorf("failed to save blend config: %w", err)
	}

	return nil
}

func (d *driver) AddDependencies(ctx context.Context, forceUpdate bool, references ...reference.Reference) error {
	d.logger.Debug("adding dependencies", map[string]interface{}{"References": references})

	// This will also include the dependencies of the dependencies
	packs, err := d.rocketPackService.GetWithContext(ctx, forceUpdate, references...)
	if err != nil {
		return fmt.Errorf("failed to get rocket packs: %w", err)
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
	_, err = d.installationService.GetWithContext(ctx, packs, false)
	if err != nil {
		return fmt.Errorf("failed to get installations: %w", err)
	}

	// Save blend config
	if err = d.save(ctx); err != nil {
		return fmt.Errorf("failed to save blend config: %w", err)
	}

	return nil
}

func (d *driver) RemoveDependencies(ctx context.Context, references ...reference.Reference) error {
	d.logger.Debug("removing dependencies", map[string]interface{}{"References": references})

	packs, err := d.rocketPackService.GetWithContext(ctx, false, references...)
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
	d.logger.Debug("installing dependencies")

	_, err := d.Get(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

func (d *driver) ResolveBlendFile(ctx context.Context) (*blendfile.BlendFile, error) {
	d.logger.Debug("resolving blend file")

	installations, err := d.Get(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get installations: %w", err)
	}

	blendFile, err := d.resolveBlendFile(ctx, installations)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve blend file: %w", err)
	}

	if err := blendfile.Validate(blendFile); err != nil {
		return nil, fmt.Errorf("invalid blend file: %w", err)
	}

	return blendFile, nil
}

func (d *driver) Get(ctx context.Context, readOnly bool) (map[reference.Reference]*installation.Installation, error) {
	packs, err := d.getDependencies(ctx)
	if err != nil {
		return nil, err
	}

	return d.installationService.GetWithContext(ctx, packs, readOnly)
}

func (d *driver) resolveBlendFile(ctx context.Context, installations map[reference.Reference]*installation.Installation) (*blendfile.BlendFile, error) {
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

	return blendFile, nil
}

func (d *driver) getDependencies(ctx context.Context) (map[reference.Reference]*rocketpack.RocketPack, error) {
	return d.rocketPackService.GetWithContext(ctx, false, d.blendConfig.RocketFile.GetDependencies()...)
}

func (d *driver) save(ctx context.Context) error {
	return rocketfile.Save(filepath.Join(d.blendConfig.ProjectPath, rocketfile.FileName), d.blendConfig.RocketFile)
}
