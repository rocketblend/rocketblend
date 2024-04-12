package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	BlenderOpts struct {
		Background bool `json:"background"`
	}

	RenderOpts struct {
		FrameStart int    `json:"frameStart" validate:"gte=0"`
		FrameEnd   int    `json:"frameEnd" validate:"gtfield=FrameStart"`
		FrameStep  int    `json:"frameStep" validate:"gte=1"`
		Output     string `json:"output"`
		Format     string `json:"format"`
		BlenderOpts
	}

	RunOpts struct {
		BlenderOpts
	}

	AddDependenciesOpts struct {
		References []reference.Reference `json:"references"`
		Fetch      bool                  `json:"fetch"`
	}

	RemoveDependenciesOpts struct {
		References []reference.Reference `json:"references"`
	}

	Driver interface {
		Render(ctx context.Context, opts *RenderOpts) error
		Run(ctx context.Context, opts *RunOpts) error
		Create(ctx context.Context) error

		InstallDependencies(ctx context.Context) error

		AddDependencies(ctx context.Context, opts *AddDependenciesOpts) error
		RemoveDependencies(ctx context.Context, opts *RemoveDependenciesOpts) error

		Resolve(ctx context.Context) (*BlendFile, error)
	}
)
