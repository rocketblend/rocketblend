package command

import (
	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command/remote"

	"github.com/spf13/cobra"
)

func NewCommand(srv *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "rocketblend-cli",
		Short: "version manager for blender.",
		Long:  `RocketBlend-Cli is a CLI tool to manage packages and installed versions of Blender.`,
	}

	c.SetVersionTemplate("{{.Version}}\n")
	c.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	removeCmd := NewRemoveCommand(srv)
	openCmd := NewOpenCommand(srv)
	createCmd := NewCreateCommand(srv)
	listCmd := NewListCommand(srv)
	exploreCmd := NewExploreCommand(srv)
	installCmd := NewInstallCommand(srv)

	remoteCmd := remote.NewCommand(srv)

	c.AddCommand(
		removeCmd,
		openCmd,
		createCmd,
		listCmd,
		installCmd,
		exploreCmd,
		remoteCmd,
	)

	return c
}
