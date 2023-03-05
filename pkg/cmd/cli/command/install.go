package command

import (
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
			var deps []string
			var rkt *rocketfile.RocketFile

			if !global {
				r, err := rocketfile.Load(srv.flags.workingDirectory)
				if err != nil {
					cmd.PrintErr(err)
					return
				}

				rkt = r // Assigning to rkt directly gives null pointer outside if statement
			}

			cmd.Println("Installing dependencies...")

			if len(args) == 1 {
				reference, err := reference.Parse(args[0])
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				pack, err := srv.driver.DescribePackByReference(reference)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				deps = append(deps, reference.String())

				if rkt != nil {
					if pack.Addon != nil {
						rkt.Addons = append(rkt.Addons, reference.String())
					}

					if pack.Build != nil {
						rkt.Build = reference.String()
					}

					deps = append([]string{rkt.Build}, rkt.Addons...)
				}
			}

			if len(deps) == 0 {
				cmd.Println("No dependencies to install")
				return
			}

			var build string
			var addons []string
			for _, dep := range deps {
				cmd.Println("Installing", dep)

				pack, err := srv.installPack(dep, force)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				if pack.Addon != nil {
					addons = append(addons, dep)
				}

				if pack.Build != nil {
					build = dep
				}
			}

			if rkt != nil && len(args) == 1 {
				cmd.Println("Updating rocketfile...")

				rkt.Build = build
				rkt.Addons = addons
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

func (srv *Service) installPack(refString string, force bool) (*rocketpack.RocketPack, error) {
	reference, err := reference.Parse(refString)
	if err != nil {
		return nil, err
	}

	// Check if already installed.
	pack, _ := srv.driver.FindPackByReference(reference)

	if pack == nil || force {
		err = srv.driver.FetchPackByReference(reference)
		if err != nil {
			return nil, err
		}

		err = srv.driver.PullPackByReference(reference)
		if err != nil {
			return nil, err
		}

		pack, err = srv.driver.FindPackByReference(reference)
		if err != nil {
			return nil, err
		}
	}

	return pack, nil
}
