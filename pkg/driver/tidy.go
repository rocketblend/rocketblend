package driver

import (
	"context"
	"sort"

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

	d.logger.Debug("tidying dependencies", map[string]interface{}{
		"dependencies": dependencies,
		"update":       update,
	})

	references := make([]reference.Reference, 0, len(dependencies))
	seen := make(map[reference.Reference]struct{})
	for _, dep := range dependencies {
		if _, exists := seen[dep.Reference]; !exists {
			references = append(references, dep.Reference)
			seen[dep.Reference] = struct{}{}
		}
	}

	results, err := d.repository.GetPackages(ctx, &types.GetPackagesOpts{
		References: references,
		Update:     update,
	})
	if err != nil {
		return nil, err
	}

	foundBuild := false
	tidied := make([]*types.Dependency, 0, len(results.Packs))

	// Iterate over the references to keep the order.
	for _, ref := range references {
		if pack, exists := results.Packs[ref]; exists {
			if pack.Type == types.PackageBuild {
				if foundBuild {
					continue
				}

				foundBuild = true
			}

			tidied = append(tidied, &types.Dependency{
				Reference: ref,
				Type:      pack.Type,
			})
		}
	}

	// Sort the dependencies by reference to make it cleaner.
	sort.Slice(tidied, func(i, j int) bool {
		return tidied[i].Reference < tidied[j].Reference
	})

	return tidied, nil
}
