package command

import "github.com/spf13/cobra"

func (srv *Service) newUninstallCommand() *cobra.Command {
	var global bool

	c := &cobra.Command{
		Use:   "uninstall [reference]",
		Short: "Remove project dependencies",
		Long:  `Remove dependencies on the current project`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("uninstall called")
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "install dependencies globally")

	return c
}
