package types

import "context"

type (
	Progress struct {
		Current int64   `json:"current"` // current bytes read
		Total   int64   `json:"total"`   // total bytes of the file
		Speed   float64 `json:"speed"`   // bytes per second
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
