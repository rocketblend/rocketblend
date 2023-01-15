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

	initCMD := NewInitCommand(srv)
	openCMD := NewOpenCommand(srv)
	fetchCMD := NewFetchCommand(srv)
	pullCMD := NewPullCommand(srv)
	findCMD := NewFindCommand(srv)
	getCMD := NewGetCommand(srv)

	c.AddCommand(
		initCMD,
		openCMD,
		fetchCMD,
		pullCMD,
		findCMD,
		getCMD,
	)

	return c
}
