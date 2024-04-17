package command

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	insertPackageOpts struct {
		commandOpts
		Reference string
	}
)

func newInsertCommand(opts commandOpts) *cobra.Command {
	cc := &cobra.Command{
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
			return runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				if err := insertPackage(ctx, insertPackageOpts{
					commandOpts: opts,
					Reference:   args[0],
				}); err != nil {
					return fmt.Errorf("failed to insert package: %w", err)
				}

				return nil
			}, &spinnerOptions{
				Suffix:  "Inserting package...",
				Verbose: opts.Global.Verbose,
			})
		},
	}

	return cc
}

func insertPackage(ctx context.Context, opts insertPackageOpts) error {
	ref, err := reference.Parse(opts.Reference)
	if err != nil {
		return err
	}

	container, err := getContainer(containerOpts{
		AppName:     opts.AppName,
		Development: opts.Development,
		Level:       opts.Global.Level,
		Verbose:     opts.Global.Verbose,
	})
	if err != nil {
		return err
	}

	validator, err := container.GetValidator()
	if err != nil {
		return err
	}

	packagePath := filepath.Join(opts.Global.WorkingDirectory, types.PackageFileName)
	pack, err := helpers.Load[types.Package](validator, packagePath)
	if err != nil {
		return err
	}

	repository, err := container.GetRepository()
	if err != nil {
		return err
	}

	if err := repository.InsertPackages(ctx, &types.InsertPackagesOpts{
		Packs: map[reference.Reference]*types.Package{
			ref: pack,
		},
	}); err != nil {
		return err
	}

	return nil
}
