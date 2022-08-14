package remote

import (
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/client"
	"github.com/spf13/cobra"
)

// remoteCmd represents the remote command
var RemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Remote is a palette that contains remote based commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

}

func NewCommand(client *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "remote",
		Short: "Remote is a palette that contains remote based commands",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	listCmd := NewListCommand(client)
	addCmd := NewAddCommand(client)
	removeCmd := NewRemoveCommand(client)

	c.AddCommand(
		listCmd,
		addCmd,
		removeCmd,
	)

	return c
}
