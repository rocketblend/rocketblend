package semver

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		major int
		minor int
		patch int
		want  string
	}{
		{1, 2, 3, "1.2.3"},
		{4, 5, 6, "4.5.6"},
		{0, 0, 0, "0.0.0"},
	}

	for _, test := range tests {
		got := NewVersion(test.major, test.minor, test.patch)
		if got.String() != test.want {
			t.Errorf("NewVersion(%d, %d, %d) = %q, want %q", test.major, test.minor, test.patch, got.String(), test.want)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		s    string
		want string
		err  bool
	}{
		{"1.2.3", "1.2.3", false},
		{"4.5.6", "4.5.6", false},
		{"0.0.0", "0.0.0", false},
		{"1.2", "", true},
		{"1.2.3.4", "", true},
		{"invalid", "", true},
	}

	for _, test := range tests {
		got, err := Parse(test.s)
		if err != nil && !test.err {
			t.Errorf("Parse(%q) returned unexpected error: %v", test.s, err)
		}
		if err == nil && test.err {
			t.Errorf("Parse(%q) did not return an error as expected", test.s)
		}
		if got != nil && got.String() != test.want {
			t.Errorf("Parse(%q) = %q, want %q", test.s, got.String(), test.want)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		data []byte
		want string
		err  bool
	}{
		{[]byte(`"1.2.3"`), "1.2.3", false},
		{[]byte(`"4.5.6"`), "4.5.6", false},
		{[]byte(`"0.0.0"`), "0.0.0", false},
		{[]byte(`"1.2"`), "", true},
		{[]byte(`"1.2.3.4"`), "", true},
		{[]byte(`"invalid"`), "", true},
		{[]byte(`true`), "", true},
		{[]byte(`123`), "", true},
		{[]byte(`[]`), "", true},
		{[]byte(`{}`), "", true},
		{[]byte(`null`), "", true},
	}

	for _, test := range tests {
		// Unmarshal the JSON data.
		var v Version
		err := json.Unmarshal(test.data, &v)
		if err != nil && !test.err {
			t.Errorf("json.Unmarshal(%q) returned unexpected error: %v", test.data, err)
		}
		if err == nil && test.err {
			t.Errorf("json.Unmarshal(%q) did not return an error as expected", test.data)
		}

		if !test.err {
			if got, want := v.String(), test.want; got != want {
				t.Errorf("json.Unmarshal(%q).String() = %q, want %q", test.data, got, want)
			}
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		v    Version
		want []byte
	}{
		{NewVersion(1, 2, 3), []byte(`"1.2.3"`)},
		{NewVersion(4, 5, 6), []byte(`"4.5.6"`)},
		{NewVersion(0, 0, 0), []byte(`"0.0.0"`)},
	}

	for _, test := range tests {
		got, err := json.Marshal(test.v)
		if err != nil {
			t.Errorf("json.Marshal(%q) returned unexpected error: %v", test.v, err)
		}
		if !bytes.Equal(got, test.want) {
			t.Errorf("json.Marshal(%q) = %q, want %q", test.v, got, test.want)
		}
	}
}

func TestUnmarshalAndMarshalRoundTrip(t *testing.T) {
	tests := []struct {
		data []byte
		want string
		err  bool
	}{
		{[]byte(`"1.2.3"`), "1.2.3", false},
		{[]byte(`"4.5.6"`), "4.5.6", false},
		{[]byte(`"0.0.0"`), "0.0.0", false},
		{[]byte(`"1.2"`), "", true},
		{[]byte(`"1.2.3.4"`), "", true},
		{[]byte(`"invalid"`), "", true},
		{[]byte(`true`), "", true},
		{[]byte(`123`), "", true},
		{[]byte(`[]`), "", true},
		{[]byte(`{}`), "", true},
		{[]byte(`null`), "", true},
	}

	for _, test := range tests {
		// Unmarshal the JSON data.
		var v Version
		err := json.Unmarshal(test.data, &v)
		if err != nil && !test.err {
			t.Errorf("json.Unmarshal(%q) returned unexpected error: %v", test.data, err)
		}
		if err == nil && test.err {
			t.Errorf("json.Unmarshal(%q) did not return an error as expected", test.data)
		}

		// Marshal the SemVer back to JSON.
		b, err := json.Marshal(v)
		if err != nil {
			t.Errorf("json.Marshal(%q) returned unexpected error: %v", v, err)
		}

		// Compare the string representation of the original and unmarshaled SemVer.
		if !test.err {
			// Compare the original and marshaled JSON data.
			if !bytes.Equal(test.data, b) {
				t.Errorf("round-trip marshaling failed for %q: got %q", test.data, b)
			}

			if got, want := v.String(), test.want; got != want {
				t.Errorf("json.Unmarshal(%q).String() = %q, want %q", test.data, got, want)
			}
		}
	}
}
