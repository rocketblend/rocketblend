package remote

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli/client"
	"github.com/spf13/cobra"
)

func NewRemoveCommand(client *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "remove",
		Short: "Removes a remote repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("remove called")
		},
	}

	return c
}
