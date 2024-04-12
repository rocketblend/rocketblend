package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Installation struct {
		Reference reference.Reference `json:"reference"` // TODO: Do I need this here?
		Path      string              `json:"path"`
		Type      PackageType         `json:"type"`
		Name      string              `json:"name"`
		Version   *semver.Version     `json:"version"`
	}

	GetInstallationOpts struct {
		References []reference.Reference `json:"references"`
		Fetch      bool                  `json:"fetch"`
	}

	GetInstallationResult struct {
		Installations []*Installation `json:"installations"`
	}

	RemoveInstallationOpts struct {
		References []reference.Reference `json:"references"`
	}

	InstallationRepository interface {
		GetInstallation(ctx context.Context, opts *GetInstallationOpts) (*GetInstallationResult, error)
		RemoveInstallation(ctx context.Context, opts *RemoveInstallationOpts) error
	}
)
