package types

import "context"

type (
	Progress struct {
		BytesRead int64   `json:"bytesRead"` // current bytes read
		TotalSize int64   `json:"totalSize"` // total bytes of the file
		Speed     float64 `json:"speed"`     // bytes per second
	}

	DownloadOpts struct {
		URI          *URI            `json:"uri" validate:"required"`
		Path         string          `json:"path" validate:"required"`
		ProgressChan chan<- Progress `json:"-"`
	}

	Downloader interface {
		Download(ctx context.Context, opts *DownloadOpts) error
	}
)
