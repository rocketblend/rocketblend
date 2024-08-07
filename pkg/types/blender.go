package types

import (
	"context"
)

const BlendFileExtension = ".blend"

type (
	Executable interface {
		Name() string
		ARGS() []string
		OutputChannel() chan string
	}

	BlendFile struct {
		Path         string          `json:"path" validate:"required,filepath,blendfile"`
		Dependencies []*Installation `json:"dependencies" validate:"required,onebuild,dive,required"`
		Strict       bool            `json:"strict"`
	}

	BlenderOpts struct {
		Background bool       `json:"background"`
		BlendFile  *BlendFile `json:"blendFile,omitempty" validate:"omitempty"`
	}

	RenderOpts struct {
		Start       int          `json:"start"`
		End         int          `json:"end"`
		Step        int          `json:"step"`
		Regions     int          `json:"regions" validate:"omitempty,oneof=1 2 4 8 16 32 64 128"`
		RegionIndex int          `json:"regionIndex" validate:"omitempty,min=0"`
		Output      string       `json:"output"`
		Format      string       `json:"format"`
		Engine      RenderEngine `json:"engine" validate:"omitempty,oneof=cycles eevee workbench"`
		Threads     int          `json:"threads" validate:"omitempty,gte=0,lte=1024"`
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
