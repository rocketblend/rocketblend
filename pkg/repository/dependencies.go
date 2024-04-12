package repository

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/types"
)

const MaxDependencyDepth = 10

func (r *repository) ResolveDependencies(ctx context.Context, opts *types.ResolveDependenciesOpts) (*types.ResolveDependenciesResult, error) {
	if err := r.validator.Validate(opts); err != nil {
		return nil, err
	}

	dependencies, err := r.resolveDependencies(ctx, opts.Dependencies)
	if err != nil {
		return nil, err
	}

	return &types.ResolveDependenciesResult{
		Dependencies: dependencies,
	}, nil
}

func (r *repository) resolveDependencies(ctx context.Context, dependencies []*types.Dependency) (*types.Dependencies, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	references := make([]reference.Reference, 0, len(dependencies))
	for _, dep := range dependencies {
		references = append(references, dep.Reference)
	}

	packs, err := r.getPackages(ctx, references, MaxDependencyDepth, false)
	if err != nil {
		return nil, err
	}

	direct := make([]*types.Dependency, len(dependencies))
	indirect := make([]*types.Dependency, len(packs)-len(dependencies))
	for ref := range packs {
		found := false
		for _, dep := range dependencies {
			if dep.Reference == ref {
				found = true
				direct = append(direct, &types.Dependency{
					Reference: ref,
					Type:      packs[ref].Type,
				})
				break
			}
		}

		if !found {
			indirect = append(indirect, &types.Dependency{
				Reference: ref,
				Type:      packs[ref].Type,
			})
		}
	}

	return &types.Dependencies{
		Direct:   direct,
		Indirect: indirect,
	}, nil
}

// func (r *repository) resolveDependencies(ctx context.Context, dependencies []*types.Dependency) ([]*types.Installation, error) {
// 	if err := ctx.Err(); err != nil {
// 		return nil, err
// 	}

// 	installations, err := r.getInstallations(ctx, dependencies, false)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result := make([]*types.Installation, 0, len(installations))
// 	for _, installation := range installations {
// 		result = append(result, installation)
// 	}

// 	return result, nil
// }
