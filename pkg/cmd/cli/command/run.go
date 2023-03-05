package command

import (
	"github.com/spf13/cobra"
)

func (srv *Service) newRunCommand() *cobra.Command {
	var background bool

	c := &cobra.Command{
		Use:   "run",
		Short: "Run project",
		Long:  `Run project`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				cmd.Println(err)
				return
			}

			err = srv.driver.Run(blend, background, []string{})
			if err != nil {
				cmd.Println(err)
				return
			}
		},
	}

	c.Flags().BoolVarP(&background, "background", "b", false, "run the project in the background")

	return c
}
