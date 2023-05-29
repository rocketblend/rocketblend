package rocketpack

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/helpers"
	"github.com/rocketblend/rocketblend/pkg/semver"
	"sigs.k8s.io/yaml"
)

const FileName = "rocketpack.yaml"

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
		Args    string                `json:"args,omitempty"`
		Version *semver.Version       `json:"version,omitempty"`
		Sources []*BuildSource        `json:"sources" validate:"required"`
		Addons  []reference.Reference `json:"addons,omitempty"`
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

func Load(filePath string) (*RocketPack, error) {
	if err := helpers.ValidateFilePath(filePath, FileName); err != nil {
		return nil, fmt.Errorf("failed to validate file path: %s", err)
	}

	if err := helpers.FileExists(filePath); err != nil {
		return nil, fmt.Errorf("failed to find blend file: %s", err)
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	var rocketPack RocketPack
	if err := yaml.Unmarshal(f, &rocketPack); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	if err := Validate(&rocketPack); err != nil {
		return nil, fmt.Errorf("failed to validate rocketfile: %w", err)
	}

	return &rocketPack, nil
}

func Validate(rp *RocketPack) error {
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
