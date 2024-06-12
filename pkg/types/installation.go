package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Installation struct {
		Path    string          `json:"path" validate:"omitempty,filepath"`
		Type    PackageType     `json:"type" validate:"required,oneof=build addon"`
		Name    string          `json:"name,omitempty"` // This is required for addons.
		Version *semver.Version `json:"version,omitempty"`
	}

	GetInstallationsOpts struct {
		Dependencies []*Dependency `json:"dependencies"`
		Fetch        bool          `json:"fetch"`
		// Progress     chan<- Progress `json:"-"`
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
