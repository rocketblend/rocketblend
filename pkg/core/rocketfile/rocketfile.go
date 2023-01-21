package rocketfile

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

const FileName = "rocketfile.yaml"

type (
	RocketFile struct {
		Build   string   `json:"build"`
		ARGS    string   `json:"args,omitempty"`
		Version string   `json:"version,omitempty"`
		Addons  []string `json:"addons,omitempty"`
	}
)

func New(rocketFile *RocketFile) (*RocketFile, error) {
	rkt := RocketFile{}

	if rocketFile != nil {
		rkt = *rocketFile
	}

	return &rkt, nil
}

func Load(dir string) (*RocketFile, error) {
	f, err := os.ReadFile(filepath.Join(dir, FileName))
	if err != nil {
		return nil, fmt.Errorf("failed to read rocketfile: %s", err)
	}

	var rkt RocketFile
	if err := yaml.Unmarshal(f, &rkt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	return &rkt, nil
}

func Save(dir string, r *RocketFile) error {
	f, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Errorf("failed to marshal rocketfile: %s", err)
	}

	if err := os.WriteFile(filepath.Join(dir, FileName), f, 0644); err != nil {
		return fmt.Errorf("failed to write rocketfile: %s", err)
	}

	return nil
}
