package downloader

import (
	"context"
	"errors"
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
		bypass                bool
		wantErr               error
	}{
		{
			name:                  "trusted domain with block policy",
			uri:                   "https://download.blender.org/release/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyBlock,
		},
		{
			name:                  "untrusted domain with block policy",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyBlock,
			wantErr:               types.ErrUnsafeDomain,
		},
		{
			name:                  "untrusted domain with block policy ignores bypass",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyBlock,
			bypass:                true,
			wantErr:               types.ErrUnsafeDomain,
		},
		{
			name:                  "untrusted domain with warn policy",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyWarn,
			wantErr:               types.ErrUntrustedDomain,
		},
		{
			name:                  "untrusted domain with warn policy and bypass",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyWarn,
			bypass:                true,
		},
		{
			name:                  "untrusted domain with allow policy",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyAllow,
		},
		{
			name:                  "untrusted domain with unknown policy",
			uri:                   "https://untrusted-site.com/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: "unknown",
			bypass:                true,
			wantErr:               types.ErrUnsafeDomain,
		},
		{
			name:                  "local file always trusted",
			uri:                   "file:///local/path/file.tar.gz",
			trustedDomains:        []string{"download.blender.org"},
			untrustedDomainPolicy: types.UntrustedDomainPolicyBlock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri, err := types.NewURI(tt.uri)
			if err != nil {
				t.Fatalf("failed to parse URI: %v", err)
			}

			err = CheckDomain(uri, tt.trustedDomains, tt.untrustedDomainPolicy, tt.bypass)
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v but got: %v", tt.wantErr, err)
			}
		})
	}

	if err := CheckDomain(nil, nil, types.UntrustedDomainPolicyBlock, false); err != nil {
		t.Errorf("expected nil URI to be allowed but got: %v", err)
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
