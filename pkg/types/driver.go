package types

import (
	"context"
)

type (
	LoadProfilesOpts struct {
		Paths []string `json:"paths"`
	}

	LoadProfilesResult struct {
		Profiles map[string]*Profile `json:"profiles"`
	}

	ResolveProfilesOpts struct {
		Profiles []*Profile `json:"profiles"`
	}

	ResolveProfilesResult struct {
		Installations [][]*Installation `json:"installations"`
	}

	TidyProfilesOpts struct {
		Profiles []*Profile `json:"profiles"`
		Fetch    bool       `json:"fetch"`
	}

	InstallProfilesOpts struct {
		Profiles []*Profile `json:"profiles"`
	}

	SaveProfilesOpts struct {
		Profiles map[string]*Profile `json:"profiles"`
	}

	Driver interface {
		LoadProfiles(ctx context.Context, opts *LoadProfilesOpts) (*LoadProfilesResult, error)
		ResolveProfiles(ctx context.Context, opts *ResolveProfilesOpts) (*ResolveProfilesResult, error)
		TidyProfiles(ctx context.Context, opts *TidyProfilesOpts) error
		InstallProfiles(ctx context.Context, opts *InstallProfilesOpts) error
		SaveProfiles(ctx context.Context, opts *SaveProfilesOpts) error
	}
)
