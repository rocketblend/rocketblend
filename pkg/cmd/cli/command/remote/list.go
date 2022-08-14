package remote

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli/client"
	"github.com/spf13/cobra"
)

func NewListCommand(client *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "Lists all remote repositories",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			remotes, err := client.GetRemotes()
			if err != nil {
				fmt.Println(err)
			}

			for _, remote := range remotes {
				fmt.Println(remote)
			}
		},
	}

	return c
}
