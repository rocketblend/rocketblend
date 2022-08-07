package remote

import (
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
