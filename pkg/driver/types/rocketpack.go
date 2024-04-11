package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/downloader"
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
		Resource string          `json:"resource,omitempty"`
		URI      *downloader.URI `json:"uri,omitempty"`
		Platform *Platform       `json:"platform,omitempty" validate:"omitempty,oneof=any windows linux macos/intel macos/apple"`
	}

	RocketPack struct {
		Spec         *semver.Version `json:"spec,omitempty"`
		Type         PackageType     `json:"type" validate:"required oneof=build addon"`
		Name         string          `json:"name,omitempty"`
		Version      *semver.Version `json:"version,omitempty"`
		Sources      []*Source       `json:"sources" validate:"required,dive"`
		Dependencies []*Dependency   `json:"dependencies,omitempty" validate:"omitempty,dive"`
	}

	GetPackageOpts struct {
		References  []reference.Reference `validate:"required"`
		ForceUpdate bool
	}

	GetPackageResult struct {
		Packs map[reference.Reference]*RocketPack
	}

	RemovePackageOpts struct {
		References []reference.Reference `validate:"required"`
	}

	InsertPackageOpts struct {
		Packs map[reference.Reference]*RocketPack `validate:"required"`
	}

	PackageRepository interface {
		GetPackage(ctx context.Context, opts *GetPackageOpts) (*GetPackageResult, error)
		RemovePackage(ctx context.Context, opts *RemovePackageOpts) error
		InsertPackage(ctx context.Context, opts *InsertPackageOpts) error
	}
)
