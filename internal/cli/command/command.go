package command

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

type (
	RootCommandOpts struct {
		Name    string
		Version string
	}

	persistentFlags struct {
		workingDirectory string
		verbose          bool
	}
)

func NewRootCommand(opts *RootCommandOpts) *cobra.Command {
	persistentFlags := persistentFlags{}

	cc := &cobra.Command{
		Version: opts.Name,
		Use:     opts.Version,
		Short:   "RocketBlend is a build and addon manager for Blender projects.",
		Long: `RocketBlend is a CLI tool that streamlines the process of managing
builds and addons for Blender projects.

Documentation is available at https://docs.rocketblend.io/`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			path, err := validatePath(persistentFlags.workingDirectory)
			if err != nil {
				return err
			}

			persistentFlags.workingDirectory = path

			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cc.SetVersionTemplate("{{.Version}}\n")

	cc.AddCommand(
		newConfigCommand(),
		// c.newNewCommand(),
		// c.newInstallCommand(),
		// c.newUninstallCommand(),
		// c.newRunCommand(),
		// c.newRenderCommand(),
		// c.newResolveCommand(),
		// c.newDescribeCommand(),
		// c.newInsertCommand(),
	)

	cc.PersistentFlags().StringVarP(&persistentFlags.workingDirectory, "directory", "d", ".", "working directory for the command")
	cc.PersistentFlags().BoolVarP(&persistentFlags.verbose, "verbose", "v", false, "enable verbose logging")

	return cc
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
