package types

import (
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

const ProjectConfigFileName = "rocketblend.yaml"

type (
	Dependency struct {
		Reference reference.Reference `json:"reference" validate:"required"`
		Type      PackageType         `json:"type,omitempty" validate:"omitempty oneof=build addon"`
	}

	Profile struct {
		Spec         semver.Version `json:"spec,omitempty"`
		Dependencies []*Dependency  `json:"dependencies,omitempty" validate:"omitempty,dive,required"`
		// ARGS         []string       `json:"args,omitempty"`
	}

	Project struct {
		BlendFilePath string   `json:"blendFilePath" validate:"required,filepath,blendfile"`
		Profile       *Profile `json:"config" validate:"required"`
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

func (p *Project) Dir() string {
	return filepath.Dir(p.BlendFilePath)
}

func (p *Project) Name() string {
	fileName := filepath.Base(p.BlendFilePath)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (p *Project) Requires() []*Dependency {
	if p.Profile == nil {
		return nil
	}

	return p.Profile.Dependencies
}
