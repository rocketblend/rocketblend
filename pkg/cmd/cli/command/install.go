package command

import "github.com/spf13/cobra"

func (srv *Service) newInstallCommand() *cobra.Command {
	var global bool

	c := &cobra.Command{
		Use:   "install [reference]",
		Short: "Install a project dependencies",
		Long:  `Adds dependencies to the current project and installs them.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("install called")
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "install dependencies globally")

	return c
}
