package command

import "github.com/spf13/cobra"

func (srv *Service) newNewCommand() *cobra.Command {
	var dir string
	var skipInstall bool

	c := &cobra.Command{
		Use:   "new [name]",
		Short: "create a new project",
		Long:  `Create a new project`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("new called")
		},
	}

	c.Flags().StringVarP(&dir, "directory", "d", "", "creates project in the specified directory (default: current directory)")
	c.Flags().BoolVarP(&skipInstall, "skip-install", "s", false, "skip installing dependencies")

	return c
}
