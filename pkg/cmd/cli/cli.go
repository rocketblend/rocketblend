package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli/client"
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command"
)

const (
	verison = "0.0.1"
	app     = "rocketblend"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot find user home directory: %v", err)
	}

	dir := filepath.Join(home, fmt.Sprintf(".%s", app))
	client, err := client.NewClient(client.Config{
		InstallationDir: dir,
	})

	if err != nil {
		return fmt.Errorf("cannot create client: %v", err)
	}

	rootCmd := command.NewCommand(client)
	return rootCmd.Execute()
}
