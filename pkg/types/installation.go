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

	GetInstallationOpts struct {
		References []reference.Reference `json:"references"`
		Fetch      bool                  `json:"fetch"`
	}

	GetInstallationResult struct {
		Installations map[reference.Reference]*Installation `json:"installations"`
	}

	RemoveInstallationOpts struct {
		References []reference.Reference `json:"references"`
	}

	InstallationRepository interface {
		GetInstallation(ctx context.Context, opts *GetInstallationOpts) (*GetInstallationResult, error)
		RemoveInstallation(ctx context.Context, opts *RemoveInstallationOpts) error
	}
)
