package types

import (
	"context"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Installation struct {
		Reference reference.Reference
		Path      string
		Type      PackageType
		Name      string
		Version   *semver.Version
	}

	GetInstallationOpts struct {
		References []reference.Reference
		Fetch      bool
	}

	GetInstallationResult struct {
		Installations []*Installation
	}

	RemoveInstallationOpts struct {
		References []reference.Reference
	}

	InstallationRepository interface {
		GetInstallation(ctx context.Context, opts *GetInstallationOpts) (*GetInstallationOpts, error)
		RemoveInstallation(ctx context.Context, opts *RemoveInstallationOpts) error
	}
)
