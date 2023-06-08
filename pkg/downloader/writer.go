package downloader

import (
	"io"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
)

type progressWriter struct {
	id      string // Unique ID for this writer in logging
	w       io.Writer
	total   int64 // Total bytes transferred
	logged  int64 // This keeps track of the last logged amount
	maxSize int64 // The total size of the download
	logFreq int64 // How often to log
	logger  logger.Logger
	logCh   chan int64
	wg      *sync.WaitGroup
}

// Write method for progressWriter with logging
func (pw *progressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.w.Write(p)
	if err != nil {
		pw.logger.Error("Error during write operation", map[string]interface{}{"error": err.Error()})
		return n, err
	}

	pw.total += int64(n)
	pw.logCh <- int64(n)
	return n, nil
}

func (pw *progressWriter) startLogging() {
	pw.logger.Info("Download started", map[string]interface{}{
		"id":       pw.id,
		"maxBytes": pw.maxSize,
	})

	pw.wg.Add(1)
	go func() {
		defer pw.wg.Done()
		for n := range pw.logCh {
			pw.logged += n
			if pw.logged >= pw.logFreq {
				pw.logger.Info("Download progress", map[string]interface{}{
					"id":         pw.id,
					"bytes":      pw.logged,
					"totalBytes": pw.total,
					"maxBytes":   pw.maxSize,
				})
				pw.logged = 0
			}
		}
	}()
}

func (pw *progressWriter) stopLogging() {
	close(pw.logCh)
	msg := "Download finished"
	if pw.total != pw.maxSize {
		msg = "Download cancelled"
	}

	pw.logger.Info(msg, map[string]interface{}{
		"id":         pw.id,
		"totalBytes": pw.total,
		"maxBytes":   pw.maxSize,
	})
}
