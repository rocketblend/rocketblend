package command

import "github.com/spf13/cobra"

func (srv *Service) newRunCommand() *cobra.Command {
	var background bool

	c := &cobra.Command{
		Use:   "run",
		Short: "Run a project",
		Long:  `Run a project`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("run called")
		},
	}

	c.Flags().BoolVarP(&background, "background", "b", false, "run the project in the background")

	return c
}
