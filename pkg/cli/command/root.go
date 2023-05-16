package command

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/cli/config"
	"github.com/rocketblend/rocketblend/pkg/cli/helpers"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"

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
		config *config.Service
		driver *rocketblend.Driver
		flags  *persistentFlags
	}
)

// NewService creates a new Service instance.
func NewService(config *config.Service, driver *rocketblend.Driver) *Service {
	return &Service{
		config: config,
		driver: driver,
		flags:  &persistentFlags{},
	}
}

// NewCommand initializes a new cobra.Command object. This is the root command.
// All other commands are subcommands of this root command.
func (srv *Service) NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   rocketblend.Name,
		Short: "RocketBlend is a build and addon manager for Blender projects.",
		Long: `RocketBlend is a CLI tool that streamlines the process of managing
builds and addons for Blender projects.

Documentation is available at https://docs.rocketblend.io/`,
		PersistentPreRun: srv.persistentPreRun,
		SilenceUsage:     true,
	}

	c.SetVersionTemplate("{{.Version}}\n")

	configCMD := srv.newConfigCommand()
	newCMD := srv.newNewCommand()
	installCMD := srv.newInstallCommand()
	uninstallCMD := srv.newUninstallCommand()
	runCMD := srv.newRunCommand()
	startCMD := srv.newStartCommand()
	renderCMD := srv.newRenderCommand()
	resolveCMD := srv.newResolveCommand()
	describeCMD := srv.newDescribeCommand()

	c.AddCommand(
		configCMD,
		newCMD,
		installCMD,
		uninstallCMD,
		runCMD,
		startCMD,
		renderCMD,
		resolveCMD,
		describeCMD,
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

// findBlendFile searches for a blend file in the provided directory.
func (srv *Service) findBlendFile(dir string) (*rocketblend.BlendFile, error) {
	dir, err := srv.validatePath(dir)
	if err != nil {
		return nil, err
	}

	path, err := helpers.FindFilePathForExt(dir, rocketblend.BlenderFileExtension)
	if err != nil {
		return nil, err
	}

	blend, err := srv.driver.Load(path)
	if err != nil {
		return nil, err
	}

	return blend, nil
}

// parseReference converts a reference string into a reference struct.
func (srv *Service) parseReference(arg string) (*reference.Reference, error) {
	r, err := reference.Parse(arg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference: %w", err)
	}

	return &r, nil
}
