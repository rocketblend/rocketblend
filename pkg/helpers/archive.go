package helpers

import (
	"path/filepath"
	"strings"
)

// SupportedArchiveExtensions is a list of supported archive extensions
var SupportedArchiveExtensions = []string{".zip", ".tar", ".tar.gz", ".tar.bz2", "tar.xz", ".dmg"}

// IsSupportedArchive returns true if the file extension is a supported archive extension
func IsSupportedArchive(filePath string) bool {
	// get the file name from the path
	_, fileName := filepath.Split(filePath)

	for _, extension := range SupportedArchiveExtensions {
		// check if the file name ends with the extension
		if strings.HasSuffix(fileName, extension) {
			return true
		}
	}

	return false
}
