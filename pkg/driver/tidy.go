package driver

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (d *Driver) TidyProfiles(ctx context.Context, opts *types.TidyProfilesOpts) error {
	if err := d.validator.Validate(opts); err != nil {
		return err
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	tasks := make([]taskrunner.Task[[]*types.Dependency], len(opts.Profiles))
	for i, profile := range opts.Profiles {
		tasks[i] = func(ctx context.Context) ([]*types.Dependency, error) {
			dependencies, err := d.tidyDependencies(ctx, profile.Dependencies, opts.Fetch)
			if err != nil {
				return nil, err
			}

			return dependencies, nil
		}
	}

	results, err := taskrunner.Run(ctx, &taskrunner.RunOpts[[]*types.Dependency]{
		Tasks:          tasks,
		Mode:           d.executionMode,
		MaxConcurrency: d.maxConcurrency,
	})
	if err != nil {
		return err
	}

	for i, profile := range opts.Profiles {
		profile.Dependencies = results[i]
	}

	return nil
}

func (d *Driver) tidyDependencies(ctx context.Context, dependencies []*types.Dependency, update bool) ([]*types.Dependency, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	references := make([]reference.Reference, 0, len(dependencies))
	for _, dep := range dependencies {
		references = append(references, dep.Reference)
	}

	results, err := d.repository.GetPackages(ctx, &types.GetPackagesOpts{
		References: references,
		Update:     update,
	})
	if err != nil {
		return nil, err
	}

	tidied := dependencies
	for ref, pack := range results.Packs {
		dependencies = append(dependencies, &types.Dependency{
			Reference: ref,
			Type:      pack.Type,
		})
	}

	return tidied, nil
}
