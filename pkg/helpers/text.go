package helpers

import (
	"path/filepath"
	"strings"
)

func ExtractName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}
