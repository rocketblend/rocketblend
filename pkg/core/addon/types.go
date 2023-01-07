package addon

import "github.com/rocketblend/rocketblend/pkg/semver"

type (
	AddonSource struct {
		File string `json:"file"`
		URL  string `json:"url"`
	}

	Addon struct {
		Reference    string         `json:"reference"`
		Name         string         `json:"name"`
		AddonVersion semver.Version `json:"addonVersion"`
		Source       AddonSource    `json:"source"`
	}
)
