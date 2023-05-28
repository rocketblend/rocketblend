package downloader

import (
	"io"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
)

type progressWriter struct {
	w       io.Writer
	total   int64
	maxSize int64
	logger  logger.Logger
	logCh   chan int64
	wg      *sync.WaitGroup
}

func (pw *progressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.w.Write(p)
	pw.total += int64(n)
	if pw.maxSize > 1<<20 { // log progress for files larger than 1 MB
		pw.logCh <- pw.total
	}
	return
}

func (pw *progressWriter) startLogging() {
	pw.logger.Info("download started", map[string]interface{}{
		"maxBytes": pw.maxSize,
	})

	if pw.maxSize > 1<<20 { // log progress for files larger than 1 MB
		pw.wg.Add(1)
		go func() {
			defer pw.wg.Done()
			for total := range pw.logCh {
				pw.logger.Info("download progress", map[string]interface{}{"bytes": total})
			}
		}()
	}
}

func (pw *progressWriter) stopLogging() {
	if pw.maxSize > 1<<20 { // log progress for files larger than 1 MB
		close(pw.logCh)
	}

	pw.logger.Info("download finished", map[string]interface{}{
		"totalBytes": pw.total,
	})
}
