package blendfile

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func (s *service) Create(ctx context.Context, blendFile *BlendFile) error {
	if err := os.MkdirAll(filepath.Dir(blendFile.FilePath), 0755); err != nil {
		return errors.Wrap(err, "error creating directories")
	}

	script, err := parseOutputTemplate(s.createScript, map[string]string{"path": blendFile.FilePath})
	if err != nil {
		return s.logAndReturnError("failed to parse output template", err)
	}

	cmd := exec.CommandContext(ctx, blendFile.Build.FilePath, "-b", "--python-expr", script)

	if err := s.runCommand(cmd); err != nil {
		return s.logAndReturnError("error running command", err)
	}

	if err := Validate(blendFile); err != nil {
		return s.logAndReturnError("invalid blend file", err)
	}

	return nil
}
