package install

import (
	"github.com/rocketblend/rocketblend/pkg/core/library"
)

type (
	Install struct {
		Id       string `json:"id"`
		Build    string `json:"build"`
		Path     string `json:"path"`
		CheckSum string `json:"checksum"`
		//Build *library.Build `json:"build"`
	}

	Pack struct {
		Id      string           `json:"id"`
		Path    string           `json:"path"`
		Package *library.Package `json:"package"`
	}
)

// func (i *Install) GetExecutableForPlatform(platform string) string {
// 	return filepath.Join(i.Path, i.Build.GetSourceForPlatform(platform).Executable)
// }
