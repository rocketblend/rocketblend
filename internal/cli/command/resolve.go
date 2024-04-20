package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	resolveProjectOpts struct {
		commandOpts
		// Format string
	}
)

// newResolveCommand creates a new cobra.Command that outputs resolved information about the project.
// This information includes dependencies and paths for the project on the local machine.
func newResolveCommand(opts commandOpts) *cobra.Command {
	//var format string

	cc := &cobra.Command{
		Use:   "resolve",
		Short: "Resolves and outputs project details",
		Long:  `Fetches and prints the resolved dependency paths for the project.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveProject(cmd.Context(), resolveProjectOpts{
				commandOpts: opts,
				// Format:      format,
			}); err != nil {
				return fmt.Errorf("failed to resolve project: %w", err)
			}

			return nil
		},
	}

	//cc.Flags().StringVarP(&format, "output", "o", "table", "output format (table, json)")

	return cc
}

func resolveProject(ctx context.Context, opts resolveProjectOpts) error {
	container, err := getContainer(containerOpts{
		AppName:     opts.AppName,
		Development: opts.Development,
		Level:       opts.Global.Level,
		Verbose:     opts.Global.Verbose,
	})
	if err != nil {
		return err
	}

	driver, err := container.GetDriver()
	if err != nil {
		return err
	}

	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{opts.Global.WorkingDirectory},
	})
	if err != nil {
		return err
	}

	resolve, err := driver.ResolveProfiles(ctx, &types.ResolveProfilesOpts{
		Profiles: profiles.Profiles,
	})
	if err != nil {
		return err
	}

	output, err := json.Marshal(resolve.Installations[0])
	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil
}
