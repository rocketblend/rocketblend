package runtime

import (
	"encoding/json"
)

type Platform int

const (
	Undefined Platform = iota
	Windows
	Linux
	DarwinAmd
	DarwinArm
)

func (p Platform) String() string {
	return [...]string{"unknown", "windows", "linux", "macos/intel", "macos/apple"}[p]
}

func (p *Platform) FromString(str string) Platform {
	return map[string]Platform{
		"unknown":     Undefined,
		"windows":     Windows,
		"linux":       Linux,
		"macos/intel": DarwinAmd,
		"macos/apple": DarwinArm,
	}[str]
}

func (p Platform) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Platform) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*p = p.FromString(s)
	return nil
}
