package cli

import (
	"context"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/rocketblend/rocketblend/internal/cli/build"
	"github.com/rocketblend/rocketblend/internal/cli/config"
	"github.com/rocketblend/rocketblend/internal/cli/factory"
	"github.com/rocketblend/rocketblend/pkg/driver"
	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
	"github.com/spf13/cobra"
)

type (
	// TOOD: Pretty sure I don't need to do this
	persistentFlags struct {
		workingDirectory string
		verbose          bool
	}

	cli struct {
		factory factory.Factory
		flags   *persistentFlags
	}
)

func New() (*cobra.Command, error) {
	cmd, err := new()
	if err != nil {
		return nil, err
	}

	rootCMD := cmd.NewRootCommand()

	// Configure help template colours
	cc.Init(&cc.Config{
		RootCmd:         rootCMD,
		Headings:        cc.Cyan + cc.Bold + cc.Underline,
		Commands:        cc.Bold,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		Aliases:         cc.Magenta,
		Example:         cc.Green + cc.Italic,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	return rootCMD, nil
}

func new() (*cli, error) {
	factory, err := factory.New()
	if err != nil {
		return nil, err
	}

	return &cli{
		factory: factory,
		flags:   &persistentFlags{},
	}, nil
}

// NewRootCommand initializes a new cobra.Command object. This is the root command.
// All other commands are subcommands of this root command.
func (c *cli) NewRootCommand() *cobra.Command {
	cc := &cobra.Command{
		Version: build.Version,
		Use:     build.AppName,
		Short:   "RocketBlend is a build and addon manager for Blender projects.",
		Long: `RocketBlend is a CLI tool that streamlines the process of managing
builds and addons for Blender projects.

Documentation is available at https://docs.rocketblend.io/`,
		PersistentPreRunE: c.persistentPreRun,
		SilenceUsage:      true,
		SilenceErrors:     true,
	}

	cc.SetVersionTemplate("{{.Version}}\n")

	cc.AddCommand(
		c.newConfigCommand(),
		c.newNewCommand(),
		c.newInstallCommand(),
		c.newUninstallCommand(),
		c.newRunCommand(),
		c.newRenderCommand(),
		c.newResolveCommand(),
		c.newDescribeCommand(),
		c.newInsertCommand(),
	)

	cc.PersistentFlags().StringVarP(&c.flags.workingDirectory, "directory", "d", ".", "working directory for the command")
	cc.PersistentFlags().BoolVarP(&c.flags.verbose, "verbose", "v", false, "enable verbose logging")

	return cc
}

func (c *cli) createDriver(blendConfig *blendconfig.BlendConfig) (driver.Driver, error) {
	logger, err := c.factory.GetLogger()
	if err != nil {
		return nil, err
	}

	rocketPackService, err := c.factory.GetRocketPackService()
	if err != nil {
		return nil, err
	}

	installationService, err := c.factory.GetInstallationService()
	if err != nil {
		return nil, err
	}

	blendFileService, err := c.factory.GetBlendFileService()
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
func (c *cli) persistentPreRun(cmd *cobra.Command, args []string) error {
	var log logger.Logger
	if !c.flags.verbose {
		log = logger.NoOp()
	}

	// Set logger per verbose flag
	err := c.factory.SetLogger(log)
	if err != nil {
		return err
	}

	path, err := c.validatePath(c.flags.workingDirectory)
	if err != nil {
		return err
	}

	c.flags.workingDirectory = path

	return nil
}

// validatePath checks if the path is valid and returns the absolute path.
func (c *cli) validatePath(path string) (string, error) {
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

func (c *cli) getBlendConfig() (*blendconfig.BlendConfig, error) {
	blendFilePath, err := findFilePathForExt(c.flags.workingDirectory, blendconfig.BlenderFileExtension)
	if err != nil {
		return nil, err
	}

	return blendconfig.Load(blendFilePath, filepath.Join(filepath.Dir(blendFilePath), rocketfile.FileName))
}

func (c *cli) getDriver() (driver.Driver, error) {
	blendConfig, err := c.getBlendConfig()
	if err != nil {
		return nil, err
	}

	return c.createDriver(blendConfig)
}

func (c *cli) getConfig() (*config.Config, error) {
	configSrv, err := c.factory.GetConfigService()
	if err != nil {
		return nil, err
	}

	return configSrv.Get()
}

func (c *cli) runWithSpinner(ctx context.Context, f func(context.Context) error, options *spinnerOptions) error {
	if options == nil {
		options = &spinnerOptions{}
	}

	if c.flags.verbose {
		return f(ctx)
	}

	return runWithSpinner(ctx, f, options)
}
