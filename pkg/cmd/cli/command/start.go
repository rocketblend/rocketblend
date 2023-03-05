package command

import (
	"github.com/spf13/cobra"
)

func (srv *Service) newStartCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Start project",
		Long:  `Start project`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				cmd.Println(err)
				return
			}

			err = srv.driver.Start(blend, []string{})
			if err != nil {
				cmd.Println(err)
				return
			}
		},
	}

	return c
}
