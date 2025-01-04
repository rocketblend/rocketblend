package reference_test

import (
	"testing"

	"github.com/rocketblend/rocketblend/pkg/reference"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
		expected  string
	}{
		{
			name:     "Valid reference",
			input:    "domain.com/base/repo/v1/builds/module/1.0",
			expected: "domain.com/base/repo/v1/builds/module/1.0",
		},
		{
			name:      "Invalid reference (missing parts)",
			input:     "domain.com/base/repo",
			expectErr: true,
		},
		{
			name:     "Valid local reference",
			input:    "local/builds/module",
			expected: "local/builds/module",
		},
		{
			name:      "Invalid local reference (missing parts)",
			input:     "local/",
			expectErr: true,
		},
		{
			name:      "Reference with extra slashes",
			input:     "domain.com/base//repo/v1//builds/module",
			expectErr: true,
		},
		{
			name:      "Empty input",
			input:     "",
			expectErr: true,
		},
		{
			name:      "Reference with only slashes",
			input:     "///",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := reference.Parse(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if !tt.expectErr && string(result) != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestAliased(t *testing.T) {
	aliases := map[string]string{
		"domain.com/base/repo/v1/builds":           "builds",
		"domain.com/base/repo/v1/addons":           "addons",
		"domain.com/base/repo/v1/builds/tools":     "tools",
		"domain.com/base/repo/v1/builds/tools/dev": "dev",
	}

	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:     "No alias match",
			input:    "domain.com/base/repo/v1/builds/example/1.0",
			expected: "domain.com/base/repo/v1/builds/example/1.0",
		},
		{
			name:     "Simple alias match",
			input:    "builds/module/1.0",
			expected: "domain.com/base/repo/v1/builds/module/1.0",
		},
		{
			name:     "Nested alias match",
			input:    "tools/utilities/v2",
			expected: "domain.com/base/repo/v1/builds/tools/utilities/v2",
		},
		{
			name:     "Deep nested alias match",
			input:    "dev/project/alpha",
			expected: "domain.com/base/repo/v1/builds/tools/dev/project/alpha",
		},
		{
			name:     "Another alias match",
			input:    "addons/theme/2.3.4",
			expected: "domain.com/base/repo/v1/addons/theme/2.3.4",
		},
		{
			name:      "No matching alias",
			input:     "unknown/path",
			expectErr: true,
		},
		{
			name:      "Empty input",
			input:     "",
			expectErr: true,
		},
		{
			name:     "Alias without additional path",
			input:    "builds",
			expected: "domain.com/base/repo/v1/builds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := reference.Aliased(tt.input, aliases)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if !tt.expectErr && string(result) != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}
