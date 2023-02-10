package command

import (
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/config"
	"github.com/rocketblend/rocketblend/pkg/core"

	"github.com/spf13/cobra"
)

type (
	Service struct {
		config *config.Service
		driver *core.Driver
	}
)

func NewService(config *config.Service, driver *core.Driver) *Service {
	return &Service{
		config: config,
		driver: driver,
	}
}

func (srv *Service) NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "rocketblend",
		Short: "RocketBlend is a build and add-ons manager for Blender.",
		Long: `RocketBlend is a powerful CLI tool that streamlines the process of managing
builds and add-ons for Blender, making installation and maintenance easier.

Documentation is available at https://docs.rocketblend.io/`,
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
	listCMD := srv.newListCommand()

	c.AddCommand(
		configCMD,
		newCMD,
		installCMD,
		uninstallCMD,
		runCMD,
		startCMD,
		renderCMD,
		resolveCMD,
		listCMD,
	)

	return c
}
