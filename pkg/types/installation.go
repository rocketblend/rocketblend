package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Installation struct {
		Path    string          `json:"path" validate:"required,filepath"`
		Type    PackageType     `json:"type,omitempty" validate:"omitempty,oneof=build addon"`
		Name    string          `json:"name" validate:"required"`
		Version *semver.Version `json:"version,omitempty"`
	}

	GetInstallationsOpts struct {
		Dependencies []*Dependency `json:"dependencies"`
		Fetch        bool          `json:"fetch"`
	}

	GetInstallationsResult struct {
		Installations map[reference.Reference]*Installation `json:"installations"`
	}

	RemoveInstallationsOpts struct {
		References []reference.Reference `json:"references"`
	}

	InstallationRepository interface {
		GetInstallations(ctx context.Context, opts *GetInstallationsOpts) (*GetInstallationsResult, error)
		RemoveInstallations(ctx context.Context, opts *RemoveInstallationsOpts) error
	}
)
