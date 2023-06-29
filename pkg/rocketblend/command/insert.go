package command

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/driver/source"
	"github.com/spf13/cobra"
)

func (srv *Service) newInsertCommand() *cobra.Command {
	var cleanup bool

	c := &cobra.Command{
		Use:   "insert",
		Short: "Inserts a package into your local library",
		Long:  `Inserts a package into your local library.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := reference.Parse(args[0])
			if err != nil {
				return err
			}

			pack, err := srv.findPackage()
			if err != nil {
				return err
			}

			// Try this first, because it's the most likely to fail.
			if err := srv.insertInstallations(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack}); err != nil {
				return err
			}

			if err := srv.insertPackages(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack}); err != nil {
				return err
			}

			return nil
		},
	}

	c.Flags().BoolVar(&cleanup, "cleanup", true, "Cleans up the package after inserting it")

	return c
}

func (srv *Service) insertPackages(ctx context.Context, packs map[reference.Reference]*rocketpack.RocketPack) error {
	packageService, err := srv.factory.GetRocketPackService()
	if err != nil {
		return err
	}

	return packageService.InsertPackages(ctx, packs)
}

func (srv *Service) insertInstallations(ctx context.Context, packs map[reference.Reference]*rocketpack.RocketPack) error {
	installationService, err := srv.factory.GetInstallationService()
	if err != nil {
		return err
	}

	sources := make(map[reference.Reference]*source.Source, len(packs))
	for ref, pack := range packs {
		if !pack.IsAddon() {
			return fmt.Errorf("inserting builds is not supported yet")
		}

		isLocalOnly, err := pack.IsLocalOnly(runtime.Undefined)
		if err != nil {
			return err
		}

		if !isLocalOnly {
			return fmt.Errorf("inserting remote addons is not supported yet")
		}

		// Addons don't have a source, so we can just use the undefined runtime.
		s, err := pack.GetSource(runtime.Undefined)
		if err != nil {
			return err
		}

		sources[ref] = s
	}

	return installationService.InsertInstallations(ctx, sources)
}

func (srv *Service) findPackage() (*rocketpack.RocketPack, error) {
	pack, err := rocketpack.Load(srv.flags.workingDirectory)
	if err != nil {
		return nil, err
	}

	return pack, nil
}
