package rocketpack

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
	"sigs.k8s.io/yaml"
)

type (
	AddonSource struct {
		File string `json:"file" validate:"required"`
		URL  string `json:"url,omitempty" validate:"url"`
	}

	Addon struct {
		Name    string          `json:"name" validate:"required"`
		Version *semver.Version `json:"version,omitempty"`
		Source  *AddonSource    `json:"source,omitempty"`
	}

	BuildSource struct {
		Platform   runtime.Platform `json:"platform" validate:"required"`
		Executable string           `json:"executable" validate:"required"`
		URL        string           `json:"url" validate:"required,url"`
	}

	Build struct {
		Args    string          `json:"args,omitempty"`
		Version *semver.Version `json:"version,omitempty"`
		Sources []*BuildSource  `json:"sources" validate:"required"`
		Addons  []string        `json:"addons,omitempty"`
	}

	RocketPack struct {
		Build *Build `json:"build,omitempty"`
		Addon *Addon `json:"addon,omitempty"`
	}
)

func (i *Build) GetSourceForPlatform(platform runtime.Platform) *BuildSource {
	if i.Sources == nil {
		return nil
	}

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

func load(bytes []byte) (*RocketPack, error) {
	pack := &RocketPack{}
	if err := yaml.Unmarshal(bytes, pack); err != nil {
		return nil, err
	}

	if pack.Build != nil && pack.Build.Version == nil {
		pack.Build.Version = &semver.Version{}
	}

	if pack.Addon != nil && pack.Addon.Version == nil {
		pack.Addon.Version = &semver.Version{}
	}

	err := validate(pack)
	if err != nil {
		return nil, err
	}

	return pack, nil
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
