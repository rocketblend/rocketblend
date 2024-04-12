package types

import "context"

type (
	ExtractOpts struct {
		Path       string `json:"path" validate:"required"`
		OutputPath string `json:"outputPath" validate:"required"`
	}

	Extractor interface {
		Extract(ctx context.Context, opts *ExtractOpts) error
	}
)
