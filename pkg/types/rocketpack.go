package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

const (
	PackageBuild PackageType = "build"
	PackageAddon PackageType = "addon"
)

type (
	PackageType string

	Source struct {
		Resource string    `json:"resource,omitempty"`
		URI      *URI      `json:"uri,omitempty"`
		Platform *Platform `json:"platform,omitempty" validate:"omitempty,oneof=any windows linux macos/intel macos/apple"`
	}

	RocketPack struct {
		Spec         *semver.Version `json:"spec,omitempty"`
		Type         PackageType     `json:"type" validate:"required oneof=build addon"`
		Name         string          `json:"name,omitempty"`
		Version      *semver.Version `json:"version,omitempty"`
		Sources      []*Source       `json:"sources" validate:"required,dive"`
		Dependencies []*Dependency   `json:"dependencies,omitempty" validate:"omitempty,dive"`
	}

	GetPackagesOpts struct {
		References []reference.Reference `validate:"required"`
		Update     bool                  `validate:"omitempty"`
		Depth      int                   `validate:"omitempty,gte=0" default:"0"` // 0 means no limit
	}

	GetPackagesResult struct {
		Packs map[reference.Reference]*RocketPack
	}

	RemovePackagesOpts struct {
		References []reference.Reference `validate:"required"`
	}

	InsertPackagesOpts struct {
		Packs map[reference.Reference]*RocketPack `validate:"required"`
	}

	PackageRepository interface {
		GetPackages(ctx context.Context, opts *GetPackagesOpts) (*GetPackagesResult, error)
		RemovePackages(ctx context.Context, opts *RemovePackagesOpts) error
		InsertPackages(ctx context.Context, opts *InsertPackagesOpts) error
	}
)

func (r *RocketPack) Bundled() bool {
	for _, s := range r.Sources {
		if s.URI != nil {
			return false
		}
	}

	return true
}
