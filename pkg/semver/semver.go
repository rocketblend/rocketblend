package semver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

// NewVersion returns a new Version with the given major, minor, and patch numbers.
func NewVersion(major, minor, patch int) Version {
	return Version{Major: major, Minor: minor, Patch: patch}
}

// String returns a string representation of the Version in the format "MAJOR.MINOR.PATCH".
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Version) UnmarshalJSON(data []byte) error {
	// Extract the string value from the JSON data.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Parse the string as a SemVer.
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid SemVer string: %q", s)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid SemVer major version: %q", parts[0])
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid SemVer minor version: %q", parts[1])
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid SemVer patch version: %q", parts[2])
	}
	*v = Version{Major: major, Minor: minor, Patch: patch}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}
