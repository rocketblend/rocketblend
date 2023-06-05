package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newConfigCommand creates a new cobra command that manages the configuration for RocketBlend.
// It either sets a new configuration value if the 'set' flag is used, or retrieves a value for the provided key.
func (srv *Service) newConfigCommand() *cobra.Command {
	var value string

	c := &cobra.Command{
		Use:   "config [key]",
		Short: "Manage the configuration for RocketBlend",
		Long:  `Fetches or sets a configuration value for RocketBlend. Provide a key to retrieve its value, or use the 'set' flag to set a new value for the key.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			config, err := srv.factory.GetConfigService()
			if err != nil {
				return fmt.Errorf("failed to create config service: %w", err)
			}

			// If the 'set' flag is used, update the configuration value for the key.
			if value != "" {
				err := config.SetValueByString(key, value)
				if err != nil {
					return fmt.Errorf("failed to set value: %w", err)
				}

				return nil
			}

			// If the 'set' flag is not used, print the current value for the key.
			cmd.Println(config.GetValueByString(key))

			return nil
		},
	}

	c.Flags().StringVarP(&value, "set", "s", "", "set a value in the config file")

	return c
}
