package command

import (
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/cli/build"
	"github.com/rocketblend/rocketblend/pkg/cli/factory"
	"github.com/rocketblend/rocketblend/pkg/cli/helpers"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/blendconfig"

	"github.com/spf13/cobra"
)

type (
	// persistentFlags holds the flags that are available across all the subcommands
	persistentFlags struct {
		workingDirectory string
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
		PersistentPreRun: srv.persistentPreRun,
		SilenceUsage:     true,
	}

	c.SetVersionTemplate("{{.Version}}\n")

	c.AddCommand(
		srv.newConfigCommand(),
		srv.newNewCommand(),
		srv.newInstallCommand(),
		srv.newUninstallCommand(),
		srv.newRunCommand(),
		srv.newStartCommand(),
		srv.newRenderCommand(),
		srv.newResolveCommand(),
		srv.newDescribeCommand(),
	)

	c.PersistentFlags().StringVarP(&srv.flags.workingDirectory, "directory", "d", ".", "working directory for the command")

	return c
}

// persistentPreRun validates the working directory before running the command.
func (srv *Service) persistentPreRun(cmd *cobra.Command, args []string) {
	path, err := srv.validatePath(srv.flags.workingDirectory)
	if err != nil {
		cmd.Println(err)
		return
	}

	srv.flags.workingDirectory = path
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

// findBlendFile finds the first Blender file in the given directory.
func (srv *Service) findBlendFilePath(dir string) (string, error) {
	dir, err := srv.validatePath(dir)
	if err != nil {
		return "", err
	}

	path, err := helpers.FindFilePathForExt(dir, blendconfig.BlenderFileExtension)
	if err != nil {
		return "", err
	}

	return path, nil
}
