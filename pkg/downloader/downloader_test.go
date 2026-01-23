package downloader

import (
	"context"
	"testing"

	"github.com/rocketblend/rocketblend/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func TestDomainSafety(t *testing.T) {
	tests := []struct {
		name                  string
		uri                   string
		trustedDomains        []string
		untrustedDomainPolicy types.UntrustedDomainPolicy
		expectError           bool
	}{
		{
			name:                  "trusted domain with block policy",
			uri:                   "https://download.blender.org/release/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyBlock,
			expectError:           false,
		},
		{
			name:                  "untrusted domain with block policy",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyBlock,
			expectError:           true,
		},
		{
			name:                  "untrusted domain with warn policy",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyWarn,
			expectError:           true,
		},
		{
			name:                  "untrusted domain with allow policy",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyAllow,
			expectError:           false,
		},
		{
			name:                  "local file always trusted",
			uri:                   "file:///local/path/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyBlock,
			expectError:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New(
				WithLogger(logger.NoOp()),
				WithTrustedDomains(tt.trustedDomains),
				WithUntrustedDomainPolicy(tt.untrustedDomainPolicy),
			)
			if err != nil {
				t.Fatalf("failed to create downloader: %v", err)
			}

			uri, err := types.NewURI(tt.uri)
			if err != nil {
				t.Fatalf("failed to parse URI: %v", err)
			}

			err = d.checkDomainSafety(uri, false)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestIsDomainTrusted(t *testing.T) {
	d, err := New(
		WithTrustedDomains([]string{
			"download.blender.org",
			"github.com",
			"gitlab.com",
		}),
	)
	if err != nil {
		t.Fatalf("failed to create downloader: %v", err)
	}

	tests := []struct {
		domain  string
		trusted bool
	}{
		{"download.blender.org", true},
		{"github.com", true},
		{"gitlab.com", true},
		{"untrusted-site.com", false},
		{"malicious.org", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			result := d.isDomainTrusted(tt.domain)
			if result != tt.trusted {
				t.Errorf("isDomainTrusted(%s) = %v, want %v", tt.domain, result, tt.trusted)
			}
		})
	}
}

func TestDownloadWithUntrustedDomain(t *testing.T) {
	ctx := context.Background()

	d, err := New(
		WithLogger(logger.NoOp()),
		WithTrustedDomains([]string{"download.blender.org"}),
		WithUntrustedDomainPolicy(types.UntrustedDomainPolicyBlock),
	)
	if err != nil {
		t.Fatalf("failed to create downloader: %v", err)
	}

	uri, err := types.NewURI("https://untrusted-site.com/file.tar.gz")
	if err != nil {
		t.Fatalf("failed to parse URI: %v", err)
	}

	opts := &types.DownloadOpts{
		URI:  uri,
		Path: "/tmp/test-file.tar.gz",
	}

	err = d.Download(ctx, opts)
	if err == nil {
		t.Error("expected error when downloading from untrusted domain with block policy")
	}
	if err != nil {
		t.Logf("got expected error: %v", err)
	}
}
