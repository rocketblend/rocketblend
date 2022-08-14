package cli

import (
	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command"
)

func Execute() error {
	conf, err := client.LoadConfig()
	if err != nil {
		return err
	}

	client, err := client.NewClient(*conf)
	if err != nil {
		return err
	}

	rootCmd := command.NewCommand(client)
	return rootCmd.Execute()
}
