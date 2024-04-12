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

	RenderBlendFileOpts struct {
		BlendFile *BlendFile `json:"blendFile" validate:"required"`
		RenderOpts
	}

	RunBlendFileOpts struct {
		BlendFile *BlendFile `json:"blendFile" validate:"required"`
		RunOpts
	}

	CreateBlendFileOpts struct {
		BlendFile *BlendFile `json:"blendFile" validate:"required"`
	}

	Blender interface {
		RenderBlendFile(ctx context.Context, opts *RenderBlendFileOpts) error
		RunBlendFile(ctx context.Context, opts *RunBlendFileOpts) error
		CreateBlendFile(ctx context.Context, blendFile *BlendFile) error
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
