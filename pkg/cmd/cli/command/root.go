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
		Short: "Version and addon manager for blender.",
		Long:  `RocketBlend is a tool for managing addons and versions of Blender.`,
	}

	c.SetVersionTemplate("{{.Version}}\n")
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	openCMD := srv.newOpenCommand()
	fetchCMD := srv.newFetchCommand()
	pullCMD := srv.newPullCommand()
	findCMD := srv.newFindCommand()
	getCMD := srv.newGetCommand()

	c.AddCommand(
		openCMD,
		fetchCMD,
		pullCMD,
		findCMD,
		getCMD,
	)

	return c
}
