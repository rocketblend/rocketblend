package driver

import (
	"context"
	"errors"

	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (d *Driver) LoadProfiles(ctx context.Context, opts *types.LoadProfilesOpts) (*types.LoadProfilesResult, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[*types.Profile], len(opts.Paths))
	for i, path := range opts.Paths {
		tasks[i] = func(ctx context.Context) (*types.Profile, error) {
			project, err := d.load(ctx, path, opts.Default)
			if err != nil {
				return nil, err
			}

			return project, nil
		}
	}

	profiles, err := taskrunner.Run(ctx, &taskrunner.RunOpts[*types.Profile]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return nil, err
	}

	return &types.LoadProfilesResult{
		Profiles: profiles,
	}, nil
}

func (d *Driver) load(ctx context.Context, path string, defaultProfile *types.Profile) (*types.Profile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	d.logger.Debug("loading profile", map[string]interface{}{
		"path": path,
	})

	profile, err := helpers.Load[types.Profile](d.validator, profileFilePath(path))
	if err != nil {
		if errors.Is(err, types.ErrFileNotFound) && defaultProfile != nil {
			return defaultProfile, nil
		}

		return nil, err
	}

	return profile, nil
}
