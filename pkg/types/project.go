package types

import (
	"path/filepath"
	"strings"
)

type (
	Project struct {
		BlendFilePath string      `json:"blendFilePath" validate:"required,filepath,blendfile"`
		RocketFile    *RocketFile `json:"rocketFile" validate:"required"`
	}
)

func (p *Project) Dir() string {
	return filepath.Dir(p.BlendFilePath)
}

func (p *Project) Name() string {
	fileName := filepath.Base(p.BlendFilePath)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (p *Project) Requires() []*Dependency {
	if p.RocketFile == nil {
		return nil
	}

	return p.RocketFile.Requires()
}
