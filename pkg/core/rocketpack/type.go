package rocketpack

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	AddonSource struct {
		File string `json:"file" validate:"required"`
		URL  string `json:"url" validate:"required,url"`
	}

	Addon struct {
		Name    string         `json:"name" validate:"required"`
		Version semver.Version `json:"version"`
		Source  AddonSource    `json:"source" validate:"required"`
	}

	BuildSource struct {
		Platform   runtime.Platform `json:"platform" validate:"required"`
		Executable string           `json:"executable" validate:"required"`
		URL        string           `json:"url" validate:"required,url"`
	}

	Build struct {
		Args    string          `json:"args,omitempty"`
		Version *semver.Version `json:"version"`
		Sources []*BuildSource  `json:"sources" validate:"required"`
		Addons  []string        `json:"addons,omitempty"`
	}

	RocketPack struct {
		// Version *semver.Version `json:"version"`
		// PackVersion *semver.Version `json:"packVersion"`
		// Explorable bool   `json:"explorable"`
		Build *Build `json:"build,omitempty"`
		Addon *Addon `json:"addon,omitempty"`
	}
)

func (i *Build) GetSourceForPlatform(platform runtime.Platform) *BuildSource {
	for _, s := range i.Sources {
		if s.Platform == platform {
			return s
		}
	}

	return nil
}

func (p *RocketPack) ToString() string {
	j, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return ""
	}

	return string(j)
}

func validate(rp *RocketPack) error {
	// See issue: https://github.com/go-playground/validator/issues/938
	if rp.Build != nil && rp.Addon != nil {
		return fmt.Errorf("packs cannot contain both a build and an addon")
	}

	if rp.Build == nil && rp.Addon == nil {
		return fmt.Errorf("packs must contain either a build or an addon")
	}

	validate := validator.New()
	err := validate.Struct(rp)
	if err != nil {
		return err
	}
	return nil
}
