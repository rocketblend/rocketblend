package command

import (
	"context"

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

			if err := srv.insertPackages(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack}); err != nil {
				return err
			}

			// Get Installations to trigger the installation process from the file system
			if err := srv.getInstallations(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack}); err != nil {
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

func (srv *Service) getInstallations(ctx context.Context, packs map[reference.Reference]*rocketpack.RocketPack) error {
	installationService, err := srv.factory.GetInstallationService()
	if err != nil {
		return err
	}

	_, err = installationService.GetInstallations(ctx, packs, false)
	return err
}

func (srv *Service) findPackage() (*rocketpack.RocketPack, error) {
	pack, err := rocketpack.Load(srv.flags.workingDirectory)
	if err != nil {
		return nil, err
	}

	return pack, nil
}
