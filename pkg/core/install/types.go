package install

import (
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/library"
)

type (
	Install struct {
		Id    string         `json:"id"`
		Path  string         `json:"path"`
		Build *library.Build `json:"build"`
	}

	Pack struct {
		Id      string           `json:"id"`
		Path    string           `json:"path"`
		Package *library.Package `json:"package"`
	}
)

func (i *Install) GetExecutablePath() string {
	// Replace with platform specific path
	return filepath.Join(i.Path, "")
}
