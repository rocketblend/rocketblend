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

func Load(path string) (*RocketFile, error) {
	f, err := os.ReadFile(filepath.Join(path, FileName))
	if err != nil {
		return nil, fmt.Errorf("failed to read rocketfile: %s", err)
	}

	var rkt RocketFile
	if err := yaml.Unmarshal(f, &rkt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	return &rkt, nil
}
