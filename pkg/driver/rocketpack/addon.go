package rocketpack

import (
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Addon struct {
		Name    string          `json:"name" validate:"required"`
		Version *semver.Version `json:"version,omitempty"`

		Source *Source `json:"source,omitempty"`
	}
)

func (a *Addon) IsPreInstalled() bool {
	return a.Source == nil
}
