package executable

import "github.com/Masterminds/semver/v3"

type (
	Addon struct {
		Name    string         `json:"name"`
		Version semver.Version `json:"version"`
		Path    string         `json:"path"`
	}

	Executable struct {
		Path   string   `json:"path"`
		Addons *[]Addon `json:"addons"`
		ARGS   string   `json:"args"`
	}
)
