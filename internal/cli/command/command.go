package command

import (
	"fmt"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/container"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	global struct {
		WorkingDirectory string
		Verbose          bool
	}

	commandOpts struct {
		AppName     string
		Development bool
		Global      *global
	}

	RootCommandOpts struct {
		Name    string
		Version string
	}
)

func NewRootCommand(opts *RootCommandOpts) *cobra.Command {
	global := global{}
	commandOpts := commandOpts{
		AppName:     opts.Name,
		Development: false,
		Global:      &global,
	}

	cc := &cobra.Command{
		Version: opts.Version,
		Use:     opts.Name,
		Short:   "RocketBlend is a build and addon manager for Blender projects.",
		Long: `RocketBlend is a CLI tool that streamlines the process of managing
builds and addons for Blender projects.

Documentation is available at https://docs.rocketblend.io/`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			path, err := validatePath(global.WorkingDirectory)
			if err != nil {
				return err
			}

			global.WorkingDirectory = path

			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cc.SetVersionTemplate("{{.Version}}\n")

	cc.AddCommand(
		newConfigCommand(commandOpts),
		newNewCommand(commandOpts),
		// c.newInstallCommand(),
		// c.newUninstallCommand(),
		newRunCommand(commandOpts),
		// c.newRenderCommand(),
		// c.newResolveCommand(),
		// c.newDescribeCommand(),
		// c.newInsertCommand(),
	)

	cc.PersistentFlags().StringVarP(&global.WorkingDirectory, "directory", "d", ".", "working directory for the command")
	cc.PersistentFlags().BoolVarP(&global.Verbose, "verbose", "v", false, "enable verbose logging")

	return cc
}

func getContainer(name string, development bool, verbose bool) (types.Container, error) {
	logLevel := "info"
	if verbose {
		logLevel = "debug"
	}

	logger := logger.New(
		logger.WithLogLevel(logLevel),
		logger.WithWriters(logger.PrettyWriter()),
	)

	container, err := container.New(
		container.WithLogger(logger),
		container.WithApplicationName(name),
		container.WithDevelopmentMode(development),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	return container, nil
}

// func createDriver(blendConfig *blendconfig.BlendConfig) (driver.Driver, error) {
// 	logger, err := c.factory.GetLogger()
// 	if err != nil {
// 		return nil, err
// 	}

// 	rocketPackService, err := c.factory.GetRocketPackService()
// 	if err != nil {
// 		return nil, err
// 	}

// 	installationService, err := c.factory.GetInstallationService()
// 	if err != nil {
// 		return nil, err
// 	}

// 	blendFileService, err := c.factory.GetBlendFileService()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return driver.New(
// 		driver.WithLogger(logger),
// 		driver.WithRocketPackService(rocketPackService),
// 		driver.WithInstallationService(installationService),
// 		driver.WithBlendFileService(blendFileService),
// 		driver.WithBlendConfig(blendConfig),
// 	)
// }

// validatePath checks if the path is valid and returns the absolute path.
func validatePath(path string) (string, error) {
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

// func getBlendConfig() (*blendconfig.BlendConfig, error) {
// 	blendFilePath, err := findFilePathForExt(c.flags.workingDirectory, blendconfig.BlenderFileExtension)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return blendconfig.Load(blendFilePath, filepath.Join(filepath.Dir(blendFilePath), rocketfile.FileName))
// }

// func getDriver() (driver.Driver, error) {
// 	blendConfig, err := c.getBlendConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c.createDriver(blendConfig)
// }

// func getConfig() (*config.Config, error) {
// 	configSrv, err := c.factory.GetConfigService()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return configSrv.Get()
// }
