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
	select {
	case <-cr.ctx.Done():
		cr.logger.Debug("context cancelled during read operation", map[string]interface{}{"id": cr.id, "error": cr.ctx.Err().Error()})
		return 0, cr.ctx.Err()
	default:
		return cr.r.Read(p)
	}
}
