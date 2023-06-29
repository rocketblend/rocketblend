package rocketpack

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/driver/helpers"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/driver/source"
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

func (r *RocketPack) IsLocalOnly(platform runtime.Platform) (bool, error) {
	if r.IsBuild() {
		return r.Build.IsLocalOnly(platform)
	}

	if r.IsAddon() {
		return r.Addon.IsLocalOnly(), nil
	}

	return false, fmt.Errorf("invalid rocket pack: neither build nor addon are defined")
}

func (r *RocketPack) IsPreInstalled() bool {
	return r.IsAddon() && r.Addon.IsPreInstalled()
}

func (r *RocketPack) GetDependencies() []reference.Reference {
	if r.IsBuild() {
		return r.Build.Addons
	}

	if r.IsAddon() {
		return nil
	}

	return nil
}

// TODO: we should be using this not GetDownloadUrl, etc. Sources should be a map with the key being the platform and the value being the source.
func (r *RocketPack) GetSource(platform runtime.Platform) (*source.Source, error) {
	if r.IsBuild() {
		return r.Build.GetSource(platform)
	}

	if r.IsAddon() {
		return r.Addon.GetSource()
	}

	return nil, fmt.Errorf("invalid rocket pack: neither build nor addon are defined")

}

func (r *RocketPack) GetDownloadUrl(platform runtime.Platform) (string, error) {
	if r.IsBuild() {
		return r.Build.GetDownloadUrl(platform)
	}

	if r.IsAddon() {
		return r.Addon.GetDownloadUrl()
	}

	return "", fmt.Errorf("invalid rocket pack: neither build nor addon are defined")
}

func (r *RocketPack) GetExecutableName(platform runtime.Platform) (string, error) {
	if r.IsBuild() {
		return r.Build.GetExecutableName(platform)
	}

	if r.IsAddon() {
		return r.Addon.GetExecutableName()
	}

	return "", fmt.Errorf("invalid rocket pack: neither build nor addon are defined")
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