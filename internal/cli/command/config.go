package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type configOpts struct {
	commandOpts
	Key   string
	Value string
}

// newConfigCommand creates a new cobra command that manages the configuration for the cli.
func newConfigCommand(opts commandOpts) *cobra.Command {
	var value string

	cc := &cobra.Command{
		Use:   "config [key]",
		Short: "Manage the configuration for RocketBlend",
		Long: `Manages RocketBlend's configuration settings.

Retrieve a config value by providing its key, or update it using --set <value>.
Use "::" for nested keys (e.g., "example::key").`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := ""
			if len(args) > 0 {
				key = args[0]
			}

			return manageConfig(configOpts{
				commandOpts: opts,
				Key:         key,
				Value:       value,
			})
		},
	}

	cc.Flags().StringVarP(&value, "set", "s", "", "set a value in the config file")

	return cc
}

// manageConfig performs the configuration management: it sets a new value if provided,
// or retrieves the current value(s) if no value is provided.
func manageConfig(opts configOpts) error {
	container, err := getContainer(containerOpts{
		AppName:     opts.AppName,
		Development: opts.Development,
		Level:       opts.Global.Level,
		Verbose:     opts.Global.Verbose,
	})
	if err != nil {
		return err
	}

	configurator, err := container.GetConfigurator()
	if err != nil {
		return err
	}

	if opts.Value != "" {
		if err := configurator.SetValueByString(opts.Key, opts.Value); err != nil {
			return fmt.Errorf("failed to set value: %w", err)
		}

		return nil
	}

	configValue, err := getConfigValue(configurator, opts.Key)
	if err != nil {
		return err
	}

	fmt.Println(configValue)

	return nil
}

func getConfigValue(configurator types.Configurator, value string) (string, error) {
	if value == "" {
		return displayJSON(map[string]interface{}{
			"path":   configurator.Path(),
			"values": configurator.GetAllValues(),
		})
	}

	return configurator.GetValueByString(value), nil
}
