package command

import (
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/cli/common"
	"github.com/rocketblend/rocketblend/pkg/cli/config"
	"github.com/rocketblend/rocketblend/pkg/cli/helpers"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"

	"github.com/spf13/cobra"
)

type (
	persistentFlags struct {
		workingDirectory string
	}

	Service struct {
		config *config.Service
		driver *rocketblend.Driver
		flags  *persistentFlags
	}
)

func NewService(config *config.Service, driver *rocketblend.Driver) *Service {
	return &Service{
		config: config,
		driver: driver,
		flags:  &persistentFlags{},
	}
}

func (srv *Service) NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     rocketblend.Name,
		Aliases: common.Aliases,
		Short:   "RocketBlend is a build and addon manager for Blender projects.",
		Long: `RocketBlend is a CLI tool that streamlines the process of managing
builds and addons for Blender projects.

Documentation is available at https://docs.rocketblend.io/`,
		PersistentPreRun: srv.persistentPreRun,
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
	// listCMD := srv.newListCommand()
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
		// listCMD,
		describeCMD,
	)

	c.PersistentFlags().StringVarP(&srv.flags.workingDirectory, "directory", "d", ".", "working directory for the command")
	// TODO: add PersistentPreRunE to validate the working directory.

	return c
}

func (srv *Service) persistentPreRun(cmd *cobra.Command, args []string) {
	path, err := srv.validatePath(srv.flags.workingDirectory)
	if err != nil {
		cmd.Println(err)
		return
	}

	srv.flags.workingDirectory = path
}

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
