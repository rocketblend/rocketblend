package command

import (
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newInstallCommand() *cobra.Command {
	var global bool
	var force bool

	c := &cobra.Command{
		Use:   "install [reference]",
		Short: "Install project dependencies",
		Long:  `Adds dependencies to the current project and installs them.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var ref *reference.Reference
			if len(args) > 0 {
				r, err := reference.Parse(args[0])
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				ref = &r
			}

			if !global {
				err := srv.driver.InstallDependencies(srv.flags.workingDirectory, ref, force)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
			} else if ref != nil {
				err := srv.driver.InstallPackByReference(*ref, force)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
			}

			cmd.Println("Dependencies installed!")
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "install dependencies globally")
	c.Flags().BoolVarP(&force, "force", "f", false, "force install dependencies (even if they are already installed)")

	return c
}

// func (srv *Service) installPack(refString string, force bool) error {
// 	reference, err := reference.Parse(refString)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if already installed.
// 	pack, _ := srv.driver.FindPackByReference(reference)

// 	if pack == nil || force {
// 		err = srv.driver.FetchPackByReference(reference)
// 		if err != nil {
// 			return err
// 		}

// 		err = srv.driver.PullPackByReference(reference)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
