package types

import (
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

const (
	ProfileDirName  = ".rocketblend"
	ProfileFileName = "profile.json"
)

type (
	Dependency struct {
		Reference reference.Reference `json:"reference" validate:"required"`
		Type      PackageType         `json:"type,omitempty" validate:"omitempty,oneof=build addon"`
	}

	Profile struct {
		Spec         semver.Version `json:"spec,omitempty"`
		Dependencies []*Dependency  `json:"dependencies,omitempty" validate:"omitempty,dive,required"`
		Strict       bool           `json:"strict,omitempty"`
		// ARGS         []string       `json:"args,omitempty"`
	}
)

func (p *Profile) FindAll(packageType PackageType) []*Dependency {
	if p.Dependencies == nil {
		return nil
	}

	var dependencies []*Dependency
	for _, d := range p.Dependencies {
		if d.Type == packageType {
			dependencies = append(dependencies, d)
		}
	}

	return dependencies
}

func (p *Profile) AddDependencies(deps ...*Dependency) {
	// Add the new dependencies to the beginning of the list so that they override any existing build dependencies.
	p.Dependencies = append(deps, p.Dependencies...)
}

func (p *Profile) RemoveDependencies(deps ...*Dependency) {
	for _, dep := range deps {
		for i, d := range p.Dependencies {
			if d.Reference == dep.Reference {
				p.Dependencies = append(p.Dependencies[:i], p.Dependencies[i+1:]...)
				break
			}
		}
	}
}
