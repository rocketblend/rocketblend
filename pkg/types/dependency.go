package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	Dependency struct {
		Reference reference.Reference `json:"reference" validate:"required"`
		Type      PackageType         `json:"type,omitempty" validate:"omitempty oneof=build addon"`
	}

	Dependencies struct {
		Direct   []*Dependency `json:"direct,omitempty" validate:"omitempty,dive,required"`
		Indirect []*Dependency `json:"indirect,omitempty" validate:"omitempty,dive,required"`
	}

	ResolveDependenciesOpts struct {
		Dependencies []*Dependency `json:"dependencies" validate:"required,dive,required"`
	}

	ResolveDependenciesResult struct {
		Dependencies *Dependencies `json:"dependencies"`
	}

	DependencyRepository interface {
		ResolveDependencies(ctx context.Context, opts *ResolveDependenciesOpts) (*ResolveDependenciesResult, error)
	}
)
