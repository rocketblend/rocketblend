package command

import "github.com/spf13/cobra"

func (srv *Service) newStartCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start project",
		Long:  `Start project`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("start called")
		},
	}

	return c
}
