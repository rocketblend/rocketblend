package types

import "context"

type (
	ExtractOpts struct {
		Path       string
		OutputPath string
	}

	Extractor interface {
		Extract(ctx context.Context, opts *ExtractOpts) error
	}
)
