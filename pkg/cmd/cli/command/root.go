package command

import (
	"github.com/rocketblend/rocketblend/pkg/client"

	"github.com/spf13/cobra"
)

func NewCommand(srv *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "rocketblend",
		Short: "Version and package manager for blender.",
		Long:  `RocketBlend is a tool for managing packages and versions of Blender.`,
	}

	c.SetVersionTemplate("{{.Version}}\n")
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd := NewInitCommand(srv)
	openCmd := NewOpenCommand(srv)
	installCmd := NewInstallCommand(srv)
	getCmd := NewGetCommand(srv)
	// addCmd := NewAddCommand(srv)

	c.AddCommand(
		initCmd,
		openCmd,
		installCmd,
		getCmd,
	)

	return c
}
