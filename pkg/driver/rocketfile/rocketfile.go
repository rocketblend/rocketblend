package rocketfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/driver/helpers"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"sigs.k8s.io/yaml"
)

const FileName = "rocketfile.yaml"

type (
	rocketFileJSON struct {
		Build   reference.Reference   `json:"build"`
		ARGS    string                `json:"args,omitempty"`
		Version string                `json:"version,omitempty"`
		Addons  []reference.Reference `json:"addons,omitempty"`
	}

	RocketFile struct {
		build   reference.Reference
		args    string
		version string
		addons  []reference.Reference
	}
)

func New(build reference.Reference, addons ...reference.Reference) *RocketFile {
	return &RocketFile{
		build:  build,
		addons: addons,
	}
}

func (r *RocketFile) GetBuild() reference.Reference {
	return r.build
}

func (r *RocketFile) SetBuild(build reference.Reference) {
	r.build = build
}

func (r *RocketFile) GetArgs() string {
	return r.args
}

func (r *RocketFile) SetArgs(args string) {
	r.args = args
}

func (r *RocketFile) GetVersion() string {
	return r.version
}

func (r *RocketFile) SetVersion(version string) {
	r.version = version
}

func (r *RocketFile) GetAddons() []reference.Reference {
	return r.addons
}

func (r *RocketFile) GetDependencies() []reference.Reference {
	var dependencies []reference.Reference
	if r.build != "" {
		dependencies = append(dependencies, r.build)
	}

	return append(dependencies, r.addons...)
}

func (r *RocketFile) AddAddons(addons ...reference.Reference) {
	r.addons = append(r.addons, addons...)
}

func (r *RocketFile) RemoveAddons(removals ...reference.Reference) {
	for _, removal := range removals {
		for i, addon := range r.addons {
			if addon == removal {
				r.addons = append(r.addons[:i], r.addons[i+1:]...)
				break
			}
		}
	}
}

func (r *RocketFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(rocketFileJSON{
		Build:   r.build,
		ARGS:    r.args,
		Version: r.version,
		Addons:  r.addons,
	})
}

func (r *RocketFile) UnmarshalJSON(data []byte) error {
	var rfj rocketFileJSON
	if err := json.Unmarshal(data, &rfj); err != nil {
		return err
	}

	r.build = rfj.Build
	r.args = rfj.ARGS
	r.version = rfj.Version
	r.addons = rfj.Addons
	return nil
}

func Load(filePath string) (*RocketFile, error) {
	if err := validateFilePath(filePath); err != nil {
		return nil, fmt.Errorf("failed to validate file path: %s", err)
	}

	if err := helpers.FileExists(filePath); err != nil {
		return nil, fmt.Errorf("failed to find rocketfile: %s", err)
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	var rocketFile RocketFile
	if err := yaml.Unmarshal(f, &rocketFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	if err := Validate(&rocketFile); err != nil {
		return nil, fmt.Errorf("failed to validate rocketfile: %w", err)
	}

	return &rocketFile, nil
}

func Save(filePath string, rocketfile *RocketFile) error {
	err := Validate(rocketfile)
	if err != nil {
		return fmt.Errorf("failed to validate rocketfile: %s", err)
	}

	err = validateFilePath(filePath)
	if err != nil {
		return fmt.Errorf("failed to validate file path: %s", err)
	}

	f, err := yaml.Marshal(rocketfile)
	if err != nil {
		return fmt.Errorf("failed to marshal rocketfile: %s", err)
	}

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	if err := os.WriteFile(filePath, f, 0644); err != nil {
		return fmt.Errorf("failed to write rocketfile: %s", err)
	}

	return nil
}

func Validate(r *RocketFile) error {
	return nil
}

func validateFilePath(filePath string) error {
	return helpers.ValidateFilePath(filePath, FileName)
}
