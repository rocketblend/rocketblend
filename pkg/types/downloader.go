package types

import "context"

type (
	DownloadOpts struct {
		URI  *URI
		Path string
	}

	Downloader interface {
		Download(ctx context.Context, opts *DownloadOpts) error
	}
)
