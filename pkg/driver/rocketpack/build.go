package rocketpack

import (
	"encoding/json"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Build struct {
		Args    string                `json:"args,omitempty"`
		Version *semver.Version       `json:"version,omitempty"`
		Sources Sources               `json:"sources,omitempty"`
		Addons  []reference.Reference `json:"addons,omitempty"`
	}

	Sources map[runtime.Platform]*Source
)

func (s *Sources) MarshalJSON() ([]byte, error) {
	result := make(map[string]*Source)
	for k, v := range *s {
		result[k.String()] = v
	}

	return json.Marshal(result)
}

func (s *Sources) UnmarshalJSON(b []byte) error {
	var raw map[string]*Source
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	*s = make(Sources)
	for k, v := range raw {
		(*s)[runtime.PlatformFromString(k)] = v
	}

	return nil
}
