package cli

import (
	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command"
)

func Execute() error {
	client, err := client.New()
	if err != nil {
		return err
	}

	rootCmd := command.NewCommand(client)
	return rootCmd.Execute()
}
