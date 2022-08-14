package remote

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewAddCommand(client *client.Client) *cobra.Command {
	var name string
	var url string

	c := &cobra.Command{
		Use:   "add",
		Short: "Adds a new remote repository",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.AddRemote(name, url); err != nil {
				fmt.Println(err)
			}
		},
	}

	c.Flags().StringVarP(&name, "name", "n", "", "The name of the remote")
	if err := c.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}

	c.Flags().StringVarP(&url, "url", "u", "", "The url of the remote")
	if err := c.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}

	return c
}
