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
		ARGS         []string       `json:"args,omitempty"`
		Dependencies *Dependencies  `json:"dependencies,omitempty" validate:"omitempty"`
	}
)

func (r *RocketFile) FindAll(packageType PackageType) []*Dependency {
	if r.Dependencies == nil {
		return nil
	}

	var dependencies []*Dependency
	for _, d := range r.Requires() {
		if d.Type == packageType {
			dependencies = append(dependencies, d)
		}
	}

	return dependencies
}

func (r *RocketFile) Requires() []*Dependency {
	if r.Dependencies == nil {
		return nil
	}

	return append(r.Dependencies.Direct, r.Dependencies.Indirect...)
}
