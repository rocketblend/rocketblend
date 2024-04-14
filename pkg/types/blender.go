package types

import (
	"context"
)

type (
	BlendFile struct {
		Name         string          `json:"name" validate:"required"`
		Path         string          `json:"path" validate:"required,filepath,blendfile"`
		Dependencies []*Installation `json:"dependencies" validate:"required,onebuild,dive,required"`
		ARGS         []string        `json:"args"`
	}

	BlenderOpts struct {
		Background   bool `json:"background"`
		ModifyAddons bool `json:"modifyAddons"`
	}

	RenderOpts struct {
		BlendFile  *BlendFile `json:"blendFile" validate:"required"`
		FrameStart int        `json:"frameStart" validate:"gte=0"`
		FrameEnd   int        `json:"frameEnd" validate:"gtfield=FrameStart"`
		FrameStep  int        `json:"frameStep" validate:"gte=1"`
		Output     string     `json:"output"`
		Format     string     `json:"format"`
		BlenderOpts
	}

	RunOpts struct {
		BlendFile *BlendFile `json:"blendFile,omitempty" validate:"omitempty,dive,required"`
		BlenderOpts
	}

	CreateOpts struct {
		BlendFile *BlendFile `json:"blendFile" validate:"required"`
	}

	Blender interface {
		Render(ctx context.Context, opts *RenderOpts) error
		Run(ctx context.Context, opts *RunOpts) error
		Create(ctx context.Context, opts *CreateOpts) error
	}
)

func (b *BlendFile) FindAll(packageType PackageType) []*Installation {
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
