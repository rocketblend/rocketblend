package types

import "context"

type (
	DownloadOpts struct {
		URI  *URI   `json:"uri" validate:"required"`
		Path string `json:"path" validate:"required"`
	}

	Downloader interface {
		Download(ctx context.Context, opts *DownloadOpts) error
	}
)
