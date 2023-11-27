package rocketpack

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/driver/helpers"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
	"sigs.k8s.io/yaml"
)

type PackType string

const (
	FileName = "rocketpack.yaml"
)

type (
	RocketPack struct {
		Build *Build `json:"build,omitempty"`
		Addon *Addon `json:"addon,omitempty"`
	}
)

func (r *RocketPack) IsBuild() bool {
	return r.Build != nil
}

func (r *RocketPack) IsAddon() bool {
	return r.Addon != nil
}

func (r *RocketPack) IsPreInstalled() bool {
	return r.IsAddon() && r.Addon.IsPreInstalled()
}

func (r *RocketPack) GetDependencies() []reference.Reference {
	if r.IsBuild() {
		return r.Build.Addons
	}

	return nil
}

func (r *RocketPack) GetSources() Sources {
	if r.IsAddon() {
		sources := make(Sources)
		sources[runtime.Undefined] = &Source{
			Resource: r.Addon.Source.Resource,
			URI:      r.Addon.Source.URI,
		}

		return sources
	}

	return r.Build.Sources
}

func (r *RocketPack) GetVersion() semver.Version {
	if r.IsAddon() {
		return *r.Addon.Version
	}

	return *r.Build.Version
}

func Load(filePath string) (*RocketPack, error) {
	if err := helpers.ValidateFilePath(filePath, FileName); err != nil {
		return nil, fmt.Errorf("failed to validate file path: %s", err)
	}

	if err := helpers.FileExists(filePath); err != nil {
		return nil, fmt.Errorf("failed to find rocketpack: %s", err)
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

func Save(filePath string, rocketPack *RocketPack) error {
	if err := Validate(rocketPack); err != nil {
		return fmt.Errorf("failed to validate rocketfile: %w", err)
	}

	f, err := yaml.Marshal(rocketPack)
	if err != nil {
		return fmt.Errorf("failed to marshal rocketfile: %s", err)
	}

	if err := os.WriteFile(filePath, f, 0644); err != nil {
		return fmt.Errorf("failed to write rocketfile: %s", err)
	}

	return nil
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
