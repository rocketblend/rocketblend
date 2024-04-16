package driver2

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (d *Driver) LoadProfiles(ctx context.Context, opts *types.LoadProfilesOpts) (*types.LoadProfilesResult, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[*getProfileResult], len(opts.Paths))
	for i, path := range opts.Paths {
		tasks[i] = func(ctx context.Context) (*getProfileResult, error) {
			project, err := d.load(ctx, path)
			if err != nil {
				return nil, err
			}

			return &getProfileResult{
				path:    path,
				profile: project,
			}, nil
		}
	}

	results, err := taskrunner.Run(ctx, &taskrunner.RunOpts[*getProfileResult]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return nil, err
	}

	profiles := make(map[string]*types.Profile, len(results))
	for _, result := range results {
		profiles[result.path] = result.profile
	}

	return &types.LoadProfilesResult{
		Profiles: profiles,
	}, nil
}

func (d *Driver) load(ctx context.Context, path string) (*types.Profile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	profile, err := helpers.Load[types.Profile](d.validator, profileFilePath(path))
	if err != nil {
		return nil, err
	}

	return profile, nil
}
