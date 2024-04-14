package types

import (
	"context"
)

type (
	Executable interface {
		Name() string
		ARGS() []string
		OutputChannel() chan string
	}

	BlendFile struct {
		Name         string          `json:"name" validate:"required"`
		Path         string          `json:"path" validate:"required,filepath,blendfile"`
		Dependencies []*Installation `json:"dependencies" validate:"required,onebuild,dive,required"`
		ARGS         []string        `json:"args"`
	}

	BlenderOpts struct {
		Background bool       `json:"background"`
		BlendFile  *BlendFile `json:"blendFile,omitempty" validate:"omitempty,dive,required"`
	}

	RenderOpts struct {
		Start         int
		End           int
		Step          int
		Output        string
		Format        RenderFormat
		CyclesDevices []CyclesDevice
		Threads       int
		BlenderOpts
	}

	RunOpts struct {
		BlenderOpts
	}

	CreateOpts struct {
		BlenderOpts
	}

	Blender interface {
		Render(ctx context.Context, opts *RenderOpts) error
		Run(ctx context.Context, opts *RunOpts) error
		Create(ctx context.Context, opts *CreateOpts) error
	}
)

func (b *BlendFile) Build() *Installation {
	builds := b.find(PackageBuild)
	if len(builds) > 0 {
		return builds[0]
	}

	return nil
}

func (b *BlendFile) Addons() []*Installation {
	return b.find(PackageAddon)
}

func (b *BlendFile) find(packageType PackageType) []*Installation {
	if b.Dependencies == nil {
		return nil
	}

	var dependencies []*Installation
	for _, d := range b.Dependencies {
		if d.Type == packageType {
			dependencies = append(dependencies, d)
		}
	}

	return dependencies
}
