package jot

import (
	"net/url"
	"path"
)

func GetFilenameFromURL(downloadURL string) string {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return ""
	}
	return path.Base(u.Path)
}
