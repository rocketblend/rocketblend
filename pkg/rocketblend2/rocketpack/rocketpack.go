package rocketpack

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/helpers"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
	"sigs.k8s.io/yaml"
)

type PackType string

const (
	TypeBuild   PackType = "Build"
	TypeAddon   PackType = "Addon"
	TypeUnknown PackType = "Unknown"

	FileName = "rocketpack.yaml"
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

func (r *RocketPack) IsBuild() bool {
	return r.Build != nil
}

func (r *RocketPack) IsAddon() bool {
	return r.Addon != nil
}

func (r *RocketPack) GetDependencies() []reference.Reference {
	if r.Build != nil {
		return r.Build.Addons
	}

	if r.Addon != nil {
		return nil
	}

	return nil
}

func (r *RocketPack) GetDownloadUrl(platform runtime.Platform) (string, error) {
	if r.Build != nil {
		return r.Build.GetDownloadUrl(platform)
	}

	if r.Addon != nil {
		return r.Addon.GetDownloadUrl()
	}

	return "", fmt.Errorf("invalid rocket pack: neither build nor addon are defined")
}

func (r *RocketPack) GetExecutableName(platform runtime.Platform) (string, error) {
	if r.Build != nil {
		return r.Build.GetExecutableName(platform)
	}

	if r.Addon != nil {
		return r.Addon.GetExecutableName()
	}

	return "", fmt.Errorf("invalid rocket pack: neither build nor addon are defined")
}

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

func (i *Build) GetDownloadUrl(platform runtime.Platform) (string, error) {
	source := i.GetSourceForPlatform(platform)
	if source == nil {
		return "", fmt.Errorf("failed to find source for platform: %s", platform)
	}

	return source.URL, nil
}

func (i *Build) GetExecutableName(platform runtime.Platform) (string, error) {
	source := i.GetSourceForPlatform(platform)
	if source == nil {
		return "", fmt.Errorf("failed to find source for platform: %s", platform)
	}

	return source.Executable, nil
}

func (a *Addon) GetDownloadUrl() (string, error) {
	if a.Source == nil {
		return "", fmt.Errorf("failed to find source for addon: %s", a.Name)
	}

	return a.Source.URL, nil
}

func (a *Addon) GetExecutableName() (string, error) {
	if a.Source == nil {
		return "", fmt.Errorf("failed to find source for addon: %s", a.Name)
	}

	return a.Source.File, nil
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

func getFilenameFromURL(downloadURL string) string {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return ""
	}

	return path.Base(u.Path)
}
