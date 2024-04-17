package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newConfigCommand creates a new cobra command that manages the configuration for RocketBlend.
// It either sets a new configuration value if the 'set' flag is used, or retrieves a value for the provided key.
func newConfigCommand(opts commandOpts) *cobra.Command {
	var value string

	cc := &cobra.Command{
		Use:   "config [key]",
		Short: "Manage the configuration for RocketBlend",
		Long:  `Fetches or sets a configuration value for RocketBlend. Provide a key to retrieve its value, or use the 'set' flag to set a new value for the key.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			container, err := getContainer(containerOpts{
				AppName:     opts.AppName,
				Development: opts.Development,
				Level:       opts.Global.Level,
				Verbose:     opts.Global.Verbose,
			})
			if err != nil {
				fmt.Println(err)
			}

			configurator, err := container.GetConfigurator()
			if err != nil {
				return err
			}

			// If the 'set' flag is used, update the configuration value for the key.
			if value != "" {
				err := configurator.SetValueByString(key, value)
				if err != nil {
					return fmt.Errorf("failed to set value: %w", err)
				}

				return nil
			}

			// If the 'set' flag is not used, print the current value for the key.
			cmd.Println(configurator.GetValueByString(key))

			return nil
		},
	}

	cc.Flags().StringVarP(&value, "set", "s", "", "set a value in the config file")

	return cc
}
