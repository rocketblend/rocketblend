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
		Short: "Version and addon manager for blender.",
		Long:  `RocketBlend is a tool for managing addons and versions of Blender.`,
	}

	c.SetVersionTemplate("{{.Version}}\n")
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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
