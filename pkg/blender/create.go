package blender

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func (b *Blender) Create(ctx context.Context, opts *types.CreateOpts) error {
	if err := b.validator.Validate(opts); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(opts.BlendFile.Path), 0755); err != nil {
		return err
	}

	build := opts.BlendFile.Build()
	if build == nil {
		return errors.New("missing build")
	}

	script, err := createBlendFileScript(&createBlendFileData{
		FilePath: opts.BlendFile.Path,
	})
	if err != nil {
		return err
	}

	b.logger.Info("blender", map[string]interface{}{
		"message": "creating blend file",
		"path":    opts.BlendFile.Path,
	})

	if err := b.execute(ctx, build.Path, &arguments{
		Script:     script,
		Background: opts.Background,
	}, nil); err != nil {
		b.logger.Error("blender", map[string]interface{}{
			"message": "error creating blend file",
			"path":    opts.BlendFile.Path,
			"error":   err.Error(),
		})

		return err
	}

	b.logger.Info("blender", map[string]interface{}{
		"message": "blend file created",
		"path":    opts.BlendFile.Path,
	})

	return nil

}
