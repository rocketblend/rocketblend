package installation

import (
	"path/filepath"
)

// SupportedArchiveExtensions is a list of supported archive extensions
var SupportedArchiveExtensions = []string{".zip", ".tar", ".tar.gz", ".tar.bz2", ".dmg"}

// isArchive returns true if the file extension is a supported archive extension
func isArchive(filePath string) bool {
	ext := filepath.Ext(filePath)
	for _, extension := range SupportedArchiveExtensions {
		if ext == extension {
			return true
		}
	}

	return false
}
