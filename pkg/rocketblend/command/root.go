package command

import (
	"context"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/driver"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/build"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/config"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/factory"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/helpers"

	"github.com/spf13/cobra"
)

type (
	// persistentFlags holds the flags that are available across all the subcommands
	persistentFlags struct {
		workingDirectory string
		verbose          bool
	}

	// Service acts as a container for the CLI services,
	// holding instances of config service, driver, and persistent flags.
	Service struct {
		factory factory.Factory
		flags   *persistentFlags
	}
)

// NewService creates a new Service instance.
func NewService() (*Service, error) {
	factory, err := factory.New()
	if err != nil {
		return nil, err
	}

	return &Service{
		factory: factory,
		flags:   &persistentFlags{},
	}, nil
}

// NewRootCommand initializes a new cobra.Command object. This is the root command.
// All other commands are subcommands of this root command.
func (srv *Service) NewRootCommand() *cobra.Command {
	c := &cobra.Command{
		Version: build.Version,
		Use:     build.AppName,
		Short:   "RocketBlend is a build and addon manager for Blender projects.",
		Long: `RocketBlend is a CLI tool that streamlines the process of managing
builds and addons for Blender projects.

Documentation is available at https://docs.rocketblend.io/`,
		PersistentPreRunE: srv.persistentPreRun,
		SilenceUsage:      true,
		SilenceErrors:     true,
	}

	c.SetVersionTemplate("{{.Version}}\n")

	c.AddCommand(
		srv.newConfigCommand(),
		srv.newNewCommand(),
		srv.newInstallCommand(),
		srv.newUninstallCommand(),
		srv.newRunCommand(),
		srv.newRenderCommand(),
		srv.newResolveCommand(),
		srv.newDescribeCommand(),
		srv.newInsertCommand(),
	)

	c.PersistentFlags().StringVarP(&srv.flags.workingDirectory, "directory", "d", ".", "working directory for the command")
	c.PersistentFlags().BoolVarP(&srv.flags.verbose, "verbose", "v", false, "enable verbose logging")

	return c
}

func (srv *Service) createDriver(blendConfig *blendconfig.BlendConfig) (driver.Driver, error) {
	logger, err := srv.factory.GetLogger()
	if err != nil {
		return nil, err
	}

	rocketPackService, err := srv.factory.GetRocketPackService()
	if err != nil {
		return nil, err
	}

	installationService, err := srv.factory.GetInstallationService()
	if err != nil {
		return nil, err
	}

	blendFileService, err := srv.factory.GetBlendFileService()
	if err != nil {
		return nil, err
	}

	return driver.New(
		driver.WithLogger(logger),
		driver.WithRocketPackService(rocketPackService),
		driver.WithInstallationService(installationService),
		driver.WithBlendFileService(blendFileService),
		driver.WithBlendConfig(blendConfig),
	)
}

// persistentPreRun validates the working directory before running the command.
func (srv *Service) persistentPreRun(cmd *cobra.Command, args []string) error {
	var log logger.Logger
	if !srv.flags.verbose {
		log = logger.NoOp()
	}

	// Set logger per verbose flag
	err := srv.factory.SetLogger(log)
	if err != nil {
		return err
	}

	path, err := srv.validatePath(srv.flags.workingDirectory)
	if err != nil {
		return err
	}

	srv.flags.workingDirectory = path

	return nil
}

// validatePath checks if the path is valid and returns the absolute path.
func (srv *Service) validatePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	// get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func (srv *Service) getBlendConfig() (*blendconfig.BlendConfig, error) {
	blendFilePath, err := helpers.FindFilePathForExt(srv.flags.workingDirectory, blendconfig.BlenderFileExtension)
	if err != nil {
		return nil, err
	}

	return blendconfig.Load(blendFilePath, filepath.Join(filepath.Dir(blendFilePath), rocketfile.FileName))
}

func (srv *Service) getDriver() (driver.Driver, error) {
	blendConfig, err := srv.getBlendConfig()
	if err != nil {
		return nil, err
	}

	return srv.createDriver(blendConfig)
}

func (srv *Service) getConfig() (*config.Config, error) {
	configSrv, err := srv.factory.GetConfigService()
	if err != nil {
		return nil, err
	}

	return configSrv.Get()
}

func (srv *Service) runWithSpinner(ctx context.Context, f func(context.Context) error, options *helpers.SpinnerOptions) error {
	if options == nil {
		options = &helpers.SpinnerOptions{}
	}

	if srv.flags.verbose {
		return f(ctx)
	}

	return helpers.RunWithSpinner(ctx, f, options)
}
