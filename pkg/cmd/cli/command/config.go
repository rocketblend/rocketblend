package command

import (
	"github.com/spf13/cobra"
)

func (srv *Service) newConfigCommand() *cobra.Command {
	var value string

	c := &cobra.Command{
		Use:   "config [key]",
		Short: "Manage the configuration for rocketblend",
		Long:  `Manage the configuration for rocketblend`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			if value != "" {
				err := srv.config.SetValueByString(key, value)
				if err != nil {
					cmd.PrintErr(err)
					return
				}
			} else {
				cmd.Println(srv.config.GetValueByString(key))
			}
		},
	}

	c.Flags().StringVarP(&value, "set", "s", "", "set a value in the config file")

	return c
}
