package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

const (
	PackageFileName = "jetpack.yaml"

	PackageBuild PackageType = "build"
	PackageAddon PackageType = "addon"
)

type (
	PackageType string

	Source struct {
		Resource string   `json:"resource,omitempty"`
		URI      *URI     `json:"uri,omitempty"`
		Platform Platform `json:"platform,omitempty" validate:"omitempty,oneof=any windows linux macos/intel macos/apple"`
	}

	Package struct {
		Spec    *semver.Version `json:"spec,omitempty"`
		Type    PackageType     `json:"type" validate:"required oneof=build addon"`
		Name    string          `json:"name,omitempty"`
		Version *semver.Version `json:"version,omitempty"`
		Sources []*Source       `json:"sources" validate:"required"`
	}

	GetPackagesOpts struct {
		References []reference.Reference `json:"references" validate:"required"`
		Update     bool                  `json:"update"`
	}

	GetPackagesResult struct {
		Packs map[reference.Reference]*Package `json:"packs"`
	}

	RemovePackagesOpts struct {
		References []reference.Reference `json:"references" validate:"required"`
	}

	InsertPackagesOpts struct {
		Packs map[reference.Reference]*Package `json:"packs" validate:"required"`
	}

	PackageRepository interface {
		GetPackages(ctx context.Context, opts *GetPackagesOpts) (*GetPackagesResult, error)
		RemovePackages(ctx context.Context, opts *RemovePackagesOpts) error
		InsertPackages(ctx context.Context, opts *InsertPackagesOpts) error
	}
)

func (r *Package) Source(platform Platform) *Source {
	var defaultSource *Source

	for _, source := range r.Sources {
		switch {
		case source.Platform == platform:
			return source
		case (source.Platform == "" || source.Platform == "any") && defaultSource == nil:
			defaultSource = source
		}
	}

	return defaultSource
}

func (r *Package) Bundled() bool {
	for _, s := range r.Sources {
		if s.URI != nil {
			return false
		}
	}

	return true
}
