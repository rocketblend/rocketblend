package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	describePackageOpts struct {
		commandOpts
		Reference string
	}
)

func newDescribeCommand(opts commandOpts) *cobra.Command {
	cc := &cobra.Command{
		Use:   "describe [reference]",
		Short: "Fetches a package definition",
		Long:  `Fetches the definition of a package by its reference.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := describePackage(cmd.Context(), describePackageOpts{
				commandOpts: opts,
				Reference:   args[0],
			}); err != nil {
				return fmt.Errorf("failed to describe package: %w", err)
			}

			return nil
		},
	}

	return cc
}

func describePackage(ctx context.Context, opts describePackageOpts) error {
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

	repository, err := container.GetRepository()
	if err != nil {
		return err
	}

	packages, err := repository.GetPackages(ctx, &types.GetPackagesOpts{
		References: []reference.Reference{ref},
	})
	if err != nil {
		return err
	}

	display, err := json.Marshal(packages.Packs[ref])
	if err != nil {
		return err
	}

	fmt.Println(string(display))

	return nil
}
