package resource

import (
	_ "embed"
)

//go:embed resources/addonScript.gopy
var addonScript string

type Service struct {
	addonScript string
}

func NewService() *Service {
	return &Service{
		addonScript: addonScript,
	}
}

func (s *Service) GetAddonScript() string {
	return s.addonScript
}
