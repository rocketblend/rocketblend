package jot

import (
	"net/url"
	"path"
	"path/filepath"
)

// SupportedArchiveExtensions is a list of supported archive extensions
var SupportedArchiveExtensions = []string{".zip", ".tar", ".tar.gz", ".tar.bz2", ".dmg"}

// GetFilenameFromURL returns the filename from a URL
func GetFilenameFromURL(downloadURL string) string {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return ""
	}
	return path.Base(u.Path)
}

// IsArchive returns true if the file extension is a supported archive extension
func isArchive(filePath string) bool {
	ext := filepath.Ext(filePath)
	for _, extension := range SupportedArchiveExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}
