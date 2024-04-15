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
	return [...]string{"undefined", "windows", "linux", "macos/intel", "macos/apple"}[p]
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

	*p = PlatformFromString(s)
	return nil
}

func PlatformFromString(str string) Platform {
	return map[string]Platform{
		"undefined":   Undefined,
		"windows":     Windows,
		"linux":       Linux,
		"macos/intel": DarwinAmd,
		"macos/apple": DarwinArm,
	}[str]
}
