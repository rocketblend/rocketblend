package command

import "github.com/spf13/cobra"

func (srv *Service) newListCommand() *cobra.Command {
	var global bool

	c := &cobra.Command{
		Use:   "list",
		Short: "Lists all the dependencies for a project",
		Long:  `Lists all the dependencies for a project`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("list called")
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "list dependencies globally")

	return c
}
