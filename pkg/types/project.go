package types

import "path/filepath"

type (
	Project struct {
		BlendFilePath string      `json:"blendFilePath" validate:"required,filepath,blendfile"`
		RocketFile    *RocketFile `json:"rocketFile" validate:"required"`
	}
)

func (p *Project) Dir() string {
	return filepath.Dir(p.BlendFilePath)
}
