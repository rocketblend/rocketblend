package command

import (
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newUninstallCommand() *cobra.Command {
	var global bool

	c := &cobra.Command{
		Use:   "uninstall [reference]",
		Short: "Remove project dependencies",
		Long:  `Remove dependencies on the current project`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ref, err := reference.Parse(args[0])
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if !global {
				err := srv.driver.UninstallDependencies(srv.flags.workingDirectory, ref)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
			} else {
				err := srv.driver.UninstallPackByReference(ref)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
			}
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "install dependencies globally")

	return c
}
