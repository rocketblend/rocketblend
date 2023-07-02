package command

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/helpers"
	"github.com/spf13/cobra"
)

func (srv *Service) newInsertCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "insert",
		Short: "Inserts a package into your local library",
		Long:  `Inserts a package into your local library.`,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasPrefix(args[0], "local/") {
				return fmt.Errorf("local package reference must start with local/")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse the reference
			ref, err := reference.Parse(args[0])
			if err != nil {
				return err
			}

			// Find the package
			pack, err := srv.findPackage()
			if err != nil {
				return err
			}

			// Insert the package into the library
			if err := srv.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				return srv.insertPackages(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack})
			}, &helpers.SpinnerOptions{Suffix: "Inserting package..."}); err != nil {
				return err
			}

			// Get Installations to trigger the installation process
			if err := srv.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				return srv.getInstallations(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack})
			}, &helpers.SpinnerOptions{Suffix: "Installing package..."}); err != nil {
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
	pack, err := rocketpack.Load(filepath.Join(srv.flags.workingDirectory, rocketpack.FileName))
	if err != nil {
		return nil, err
	}

	return pack, nil
}
