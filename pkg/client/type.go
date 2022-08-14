package client

import (
	"github.com/blang/semver/v4"
	"go.lsp.dev/uri"
)

type Available struct {
	Hash    string
	Name    string
	Uri     uri.URI
	Version semver.Version
}
