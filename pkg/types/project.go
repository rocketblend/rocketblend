package types

import (
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/semver"
)

const ProjectConfigFileName = "rocketblend.json"

type (
	ProjectConfig struct {
		Spec         semver.Version `json:"spec,omitempty"`
		ARGS         []string       `json:"args,omitempty"`
		Dependencies *Dependencies  `json:"dependencies,omitempty" validate:"omitempty,dive,required"`
	}

	Project struct {
		BlendFilePath string         `json:"blendFilePath" validate:"required,filepath,blendfile"`
		Config        *ProjectConfig `json:"config" validate:"required"`
	}
)

func (r *ProjectConfig) FindAll(packageType PackageType) []*Dependency {
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

func (r *ProjectConfig) Requires() []*Dependency {
	if r.Dependencies == nil {
		return nil
	}

	return append(r.Dependencies.Direct, r.Dependencies.Indirect...)
}

func (p *Project) Dir() string {
	return filepath.Dir(p.BlendFilePath)
}

func (p *Project) Name() string {
	fileName := filepath.Base(p.BlendFilePath)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (p *Project) Requires() []*Dependency {
	if p.Config == nil {
		return nil
	}

	return p.Config.Requires()
}
