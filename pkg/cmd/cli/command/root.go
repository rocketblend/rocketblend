package command

import (
	"github.com/rocketblend/rocketblend/pkg/core"

	"github.com/spf13/cobra"
)

type (
	Service struct {
		driver *core.Driver
	}
)

func NewService(driver *core.Driver) *Service {
	return &Service{
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

	openCMD := srv.newOpenCommand()
	fetchCMD := srv.newFetchCommand()
	pullCMD := srv.newPullCommand()
	findCMD := srv.newFindCommand()
	getCMD := srv.newGetCommand()
	createCMD := srv.newCreateCommand()

	c.AddCommand(
		openCMD,
		fetchCMD,
		pullCMD,
		findCMD,
		getCMD,
		createCMD,
	)

	return c
}
