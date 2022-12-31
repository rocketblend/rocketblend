package library

import (
	"github.com/Masterminds/semver/v3"
)

type (
	PackageSource struct {
		File string `json:"file"`
		URL  string `json:"url"`
	}

	Package struct {
		Reference    string         `json:"reference"`
		Name         string         `json:"name"`
		AddonVersion semver.Version `json:"addonVersion"`
		Source       PackageSource  `json:"source"`
	}
)
