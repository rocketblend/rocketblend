package supervisor

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"slices"

	"github.com/rocketblend/rocketblend/pkg/types"
)

var ignoredDirs = []string{".git", "node_modules", "vendor", ".idea", ".vscode"}

func Discover(ctx context.Context, root string) (map[string]string, error) {
	projectFiles := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if d.IsDir() {
			if slices.Contains(ignoredDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(d.Name(), types.BlendFileExtension) {
			dir := filepath.Dir(path)
			if _, exists := projectFiles[dir]; !exists {
				projectFiles[dir] = path
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return projectFiles, nil
}
