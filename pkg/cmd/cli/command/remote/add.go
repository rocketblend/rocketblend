package remote

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli/client"
	"github.com/spf13/cobra"
)

func NewAddCommand(client *client.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "Adds a new remote repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("add called")
		},
	}

	return c
}
