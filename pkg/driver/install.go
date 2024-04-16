package driver2

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (d *Driver) InstallProfiles(ctx context.Context, opts *types.InstallProfilesOpts) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[struct{}], len(opts.Profiles))
	for i, profile := range opts.Profiles {
		tasks[i] = func(ctx context.Context) (struct{}, error) {
			if err := d.installDependencies(ctx, profile); err != nil {
				return struct{}{}, err
			}

			return struct{}{}, nil
		}
	}

	_, err := taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) installDependencies(ctx context.Context, profile *types.Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := d.getInstallations(ctx, profile.Dependencies, true)
	if err != nil {
		return err
	}

	return nil
}
