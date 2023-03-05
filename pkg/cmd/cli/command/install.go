package command

import (
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/helpers"
	"github.com/rocketblend/rocketblend/pkg/core/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/core/rocketpack"
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
			// Parse reference if specified
			var pack *rocketpack.RocketPack
			if len(args) > 0 {
				var err error
				reference, err := reference.Parse(args[0])
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				pack, err = srv.driver.DescribePackByReference(reference)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
			}

			var dependencies []string
			var rkt *rocketfile.RocketFile
			if !global {
				var err error
				rkt, err = rocketfile.Load(srv.flags.workingDirectory)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				// Add the specified pack to the rocketfile.
				if pack != nil {
					if pack.Addon != nil {
						rkt.Addons = helpers.RemoveDuplicateStr(append(rkt.Addons, args[0]))
					}

					if pack.Build != nil {
						rkt.Build = args[0]
					}
				}

				dependencies = append([]string{rkt.Build}, rkt.Addons...)
			}

			// Install a specific pack globally.
			if global && pack != nil {
				dependencies = append(dependencies, args[0])
			}

			cmd.Println("Installing dependencies...")

			for _, dep := range dependencies {
				cmd.Println("Installing", dep)

				err := srv.installPack(dep, force)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
			}

			// Save the rocketfile. This will only happen if the rocketfile was loaded and a pack was specified.
			if rkt != nil && pack != nil {
				cmd.Println("Updating rocketfile...")

				err := rocketfile.Save(srv.flags.workingDirectory, rkt)
				if err != nil {
					cmd.PrintErr(err)
					return
				}
			}

			cmd.Println("Done!")
		},
	}

	c.Flags().BoolVarP(&global, "global", "g", false, "install dependencies globally")
	c.Flags().BoolVarP(&force, "force", "f", false, "force install dependencies (even if they are already installed)")

	return c
}

func (srv *Service) installPack(refString string, force bool) error {
	reference, err := reference.Parse(refString)
	if err != nil {
		return err
	}

	// Check if already installed.
	pack, _ := srv.driver.FindPackByReference(reference)

	if pack == nil || force {
		err = srv.driver.FetchPackByReference(reference)
		if err != nil {
			return err
		}

		err = srv.driver.PullPackByReference(reference)
		if err != nil {
			return err
		}
	}

	return nil
}
