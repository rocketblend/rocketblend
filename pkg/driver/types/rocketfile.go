package types

import (
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Dependencies struct {
		Direct   []*Dependency `json:"direct,omitempty" validate:"omitempty,dive"`
		Indirect []*Dependency `json:"indirect,omitempty" validate:"omitempty,dive"`
	}

	RocketFile struct {
		Spec         semver.Version `json:"spec,omitempty"`
		Args         []string       `json:"args,omitempty"`
		Dependencies *Dependencies  `json:"dependencies,omitempty" validate:"omitempty"`
	}
)
