package remote

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewRemoveCommand(client *client.Client) *cobra.Command {
	var name string

	c := &cobra.Command{
		Use:   "remove",
		Short: "Removes a remote repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.RemoveRemote(name); err != nil {
				fmt.Println(err)
			}
		},
	}

	c.Flags().StringVarP(&name, "name", "n", "", "The name of the remote")
	if err := c.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}

	return c
}
