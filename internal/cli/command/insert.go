package command

// import (
// 	"context"
// 	"fmt"
// 	"path/filepath"
// 	"strings"

// 	"github.com/rocketblend/rocketblend/pkg/reference"
// 	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
// 	"github.com/spf13/cobra"
// )

// func newInsertCommand() *cobra.Command {
// 	cc := &cobra.Command{
// 		Use:   "insert",
// 		Short: "Inserts a package into your local library",
// 		Long:  `Inserts a package into your local library.`,
// 		Args:  cobra.ExactArgs(1),
// 		PreRunE: func(cmd *cobra.Command, args []string) error {
// 			if !strings.HasPrefix(args[0], "local/") {
// 				return fmt.Errorf("local package reference must start with local/")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			// Parse the reference
// 			ref, err := reference.Parse(args[0])
// 			if err != nil {
// 				return err
// 			}

// 			// Find the package
// 			pack, err := c.findPackage()
// 			if err != nil {
// 				return err
// 			}

// 			// Insert the package into the library
// 			if err := c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
// 				return c.Insert(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack})
// 			}, &spinnerOptions{Suffix: "Inserting package..."}); err != nil {
// 				return err
// 			}

// 			// Get Installations to trigger the installation process
// 			if err := c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
// 				return c.Get(cmd.Context(), map[reference.Reference]*rocketpack.RocketPack{ref: pack})
// 			}, &spinnerOptions{Suffix: "Installing package..."}); err != nil {
// 				return err
// 			}

// 			return nil
// 		},
// 	}

// 	return cc
// }

// func (c *cli) Insert(ctx context.Context, packs map[reference.Reference]*rocketpack.RocketPack) error {
// 	packageService, err := c.factory.GetRocketPackService()
// 	if err != nil {
// 		return err
// 	}

// 	return packageService.Insert(ctx, packs)
// }

// func (c *cli) Get(ctx context.Context, packs map[reference.Reference]*rocketpack.RocketPack) error {
// 	installationService, err := c.factory.GetInstallationService()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = installationService.Get(ctx, packs, false)
// 	return err
// }

// func (c *cli) findPackage() (*rocketpack.RocketPack, error) {
// 	pack, err := rocketpack.Load(filepath.Join(c.flags.workingDirectory, rocketpack.FileName))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pack, nil
// }
