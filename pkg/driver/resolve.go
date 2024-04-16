package driver

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (d *Driver) ResolveProfiles(ctx context.Context, opts *types.ResolveProfilesOpts) (*types.ResolveProfilesResult, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[[]*types.Installation], len(opts.Profiles))
	for i, profile := range opts.Profiles {
		tasks[i] = func(ctx context.Context) ([]*types.Installation, error) {
			installations, err := d.resolve(ctx, profile)
			if err != nil {
				return nil, err
			}

			return installations, nil
		}
	}

	installations, err := taskrunner.Run(ctx, &taskrunner.RunOpts[[]*types.Installation]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return nil, err
	}

	return &types.ResolveProfilesResult{
		Installations: installations,
	}, nil
}

func (d *Driver) resolve(ctx context.Context, profile *types.Profile) ([]*types.Installation, error) {
	installations, err := d.getInstallations(ctx, profile.Dependencies, false)
	if err != nil {
		return nil, err
	}

	d.logger.Debug("resolving profile", map[string]interface{}{
		"profile":       profile,
		"installations": installations,
	})

	dependencies := make([]*types.Installation, 0, len(installations))
	for _, installation := range installations {
		dependencies = append(dependencies, installation)
	}

	return dependencies, nil
}
