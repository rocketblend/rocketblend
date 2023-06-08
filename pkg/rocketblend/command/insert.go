package command

import (
	"context"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/spf13/cobra"
)

func (srv *Service) newInsertCommand() *cobra.Command {
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

			if pack.IsBuild() {
				return fmt.Errorf("inserting builds is not supported yet")
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

	return installationService.InsertInstallations(ctx, packs, srv.flags.workingDirectory)
}

func (srv *Service) findPackage() (*rocketpack.RocketPack, error) {
	pack, err := rocketpack.Load(srv.flags.workingDirectory)
	if err != nil {
		return nil, err
	}

	return pack, nil
}
