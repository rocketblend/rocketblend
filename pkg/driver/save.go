package driver

import (
	"context"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (d *Driver) SaveProfiles(ctx context.Context, opts *types.SaveProfilesOpts) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[struct{}], 0, len(opts.Profiles))
	for path, profile := range opts.Profiles {
		tasks = append(tasks, func(ctx context.Context) (struct{}, error) {
			return struct{}{}, d.save(ctx, path, profile)
		})
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

func (d *Driver) save(ctx context.Context, path string, profile *types.Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.logger.Debug("saving profile", map[string]interface{}{
		"path":    path,
		"profile": profile,
	})

	if err := helpers.Save(d.validator, profileFilePath(path), profile); err != nil {
		return err
	}

	return nil
}

func profileFilePath(path string) string {
	return filepath.Join(path, types.ProfileFileName)
}
