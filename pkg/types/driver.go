package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	AddDependenciesOpts struct {
		References []reference.Reference `json:"references"`
		Fetch      bool                  `json:"fetch"`
	}

	RemoveDependenciesOpts struct {
		References []reference.Reference `json:"references"`
	}

	Driver interface {
		InstallDependencies(ctx context.Context) error

		AddDependencies(ctx context.Context, opts *AddDependenciesOpts) error
		RemoveDependencies(ctx context.Context, opts *RemoveDependenciesOpts) error

		Resolve(ctx context.Context) (*BlendFile, error)
		Save(ctx context.Context) error
	}
)
