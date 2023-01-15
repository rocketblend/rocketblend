package rocketpack

import (
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
		Args    string          `json:"args" validate:"required"`
		Version *semver.Version `json:"version"`
		Source  []*BuildSource  `json:"source" validate:"required"`
		Addons  []string        `json:"addons"`
	}

	RocketPack struct {
		Version *semver.Version `json:"version"`
		Build   *Build          `json:"build" validate:"excluded_unless=Addon"`
		Addon   *Addon          `json:"addon" validate:"excluded_unless=Build"`
	}
)

func (i *Build) GetSourceForPlatform(platform runtime.Platform) *BuildSource {
	for _, s := range i.Source {
		if s.Platform == platform {
			return s
		}
	}

	return nil
}

func validate(rp *RocketPack) error {
	validate := validator.New()
	err := validate.Struct(rp)
	if err != nil {
		return err
	}
	return nil
}
