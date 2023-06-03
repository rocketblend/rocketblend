package rocketpack

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/helpers"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
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

// IsBuild returns true if the rocket pack is a build.
func (r *RocketPack) IsBuild() bool {
	return r.Build != nil
}

// IsAddon returns true if the rocket pack is an addon.
func (r *RocketPack) IsAddon() bool {
	return r.Addon != nil
}

// IsLocalOnly returns true if the build is local only, meaning it is pre-installed or has no download URL.
func (r *RocketPack) IsLocalOnly(platform runtime.Platform) (bool, error) {
	if r.IsBuild() {
		return r.Build.IsLocalOnly(platform)
	}

	if r.IsAddon() {
		return r.Addon.IsLocalOnly(), nil
	}

	return false, fmt.Errorf("invalid rocket pack: neither build nor addon are defined")
}

// IsPreInstalled returns true if the package is pre-installed.
func (r *RocketPack) IsPreInstalled() bool {
	return r.IsAddon() && r.Addon.IsPreInstalled()
}

// GetDependencies returns the dependencies for the rocket pack.
func (r *RocketPack) GetDependencies() []reference.Reference {
	if r.IsBuild() {
		return r.Build.Addons
	}

	if r.IsAddon() {
		return nil
	}

	return nil
}

// GetDownloadUrl returns the download URL for the rocket pack.
func (r *RocketPack) GetDownloadUrl(platform runtime.Platform) (string, error) {
	if r.IsBuild() {
		return r.Build.GetDownloadUrl(platform)
	}

	if r.IsAddon() {
		return r.Addon.GetDownloadUrl()
	}

	return "", fmt.Errorf("invalid rocket pack: neither build nor addon are defined")
}

// GetExecutableName returns the executable name for the rocket pack.
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
