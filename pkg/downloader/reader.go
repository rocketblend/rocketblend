package downloader

import (
	"context"
	"io"

	"github.com/flowshot-io/x/pkg/logger"
)

type (
	contextReader struct {
		r      io.Reader
		id     string
		ctx    context.Context
		logger logger.Logger
	}
)

// Read method for contextReader
func (cr *contextReader) Read(p []byte) (n int, err error) {
	if err := cr.ctx.Err(); err != nil {
		return 0, err
	}

	return cr.r.Read(p)
}
