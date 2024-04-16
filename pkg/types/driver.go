package types

import (
	"context"
)

type (
	LoadProfilesOpts struct {
		Paths []string `json:"paths" validate:"required,dive,dir"`
	}

	LoadProfilesResult struct {
		Profiles []*Profile `json:"profiles"`
	}

	ResolveProfilesOpts struct {
		Profiles []*Profile `json:"profiles" validate:"required,dive,required"`
	}

	ResolveProfilesResult struct {
		Installations [][]*Installation `json:"installations"`
	}

	TidyProfilesOpts struct {
		Profiles []*Profile `json:"profiles" validate:"required,dive,required"`
		Fetch    bool       `json:"fetch" validate:"required"`
	}

	InstallProfilesOpts struct {
		Profiles []*Profile `json:"profiles" validate:"required,dive,required"`
	}

	SaveProfilesOpts struct {
		Profiles map[string]*Profile `json:"profiles" validate:"required,dive,required"`
	}

	Driver interface {
		LoadProfiles(ctx context.Context, opts *LoadProfilesOpts) (*LoadProfilesResult, error)
		ResolveProfiles(ctx context.Context, opts *ResolveProfilesOpts) (*ResolveProfilesResult, error)
		TidyProfiles(ctx context.Context, opts *TidyProfilesOpts) error
		InstallProfiles(ctx context.Context, opts *InstallProfilesOpts) error
		SaveProfiles(ctx context.Context, opts *SaveProfilesOpts) error
	}
)
