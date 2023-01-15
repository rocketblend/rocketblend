package command

import (
	"github.com/rocketblend/rocketblend/pkg/client"

	"github.com/spf13/cobra"
)

func NewCommand(srv *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "rocketblend",
		Short: "Version and addon manager for blender.",
		Long:  `RocketBlend is a tool for managing addons and versions of Blender.`,
	}

	c.SetVersionTemplate("{{.Version}}\n")
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd := NewInitCommand(srv)
	openCmd := NewOpenCommand(srv)
	installCmd := NewInstallCommand(srv)
	getCmd := NewGetCommand(srv)

	fetchCMD := NewFetchCommand(srv)
	pullCMD := NewPullCommand(srv)
	findCMD := NewFindCommand(srv)
	grabCMD := NewGrabCommand(srv)

	c.AddCommand(
		initCmd,
		openCmd,
		installCmd,
		getCmd,
		fetchCMD,
		pullCMD,
		findCMD,
		grabCMD,
	)

	return c
}
