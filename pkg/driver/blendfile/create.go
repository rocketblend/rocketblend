package blendfile

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (s *service) Create(blendFile *BlendFile) error {
	return s.CreateWithContext(context.Background(), blendFile)
}

func (s *service) CreateWithContext(ctx context.Context, blendFile *BlendFile) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(blendFile.FilePath), 0755); err != nil {
		return errors.Wrap(err, "error creating directories")
	}

	script, err := parseOutputTemplate(s.createScript, map[string]string{"path": blendFile.FilePath})
	if err != nil {
		return s.logAndReturnError("failed to parse output template", err)
	}

	args := []string{
		"-b",
		"--python-expr",
		script,
	}

	if err := s.runCommand(ctx, blendFile.Build.FilePath, args...); err != nil {
		return s.logAndReturnError("error running command", err)
	}

	if err := Validate(blendFile); err != nil {
		return s.logAndReturnError("invalid blend file", err)
	}

	return nil
}
