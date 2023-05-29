package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
)

type (
	Downloader interface {
		Download(path string, downloadUrl string) error
		DownloadWithContext(ctx context.Context, path string, downloadUrl string) error
	}

	Options struct {
		Logger  logger.Logger
		LogFreq int64
	}

	Option func(*Options)

	downloader struct {
		logger  logger.Logger
		logFreq int64
	}
)

// With Logger sets the logger to use. The default is no-op.
func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithLogFrequency sets the frequency of log messages. The default is 1MB.
func WithLogFrequency(logFreq int64) Option {
	return func(o *Options) {
		o.LogFreq = logFreq
	}
}

// New creates a new Downloader.
func New(opts ...Option) Downloader {
	options := &Options{
		Logger:  logger.NoOp(),
		LogFreq: 1 << 20, // Default log frequency is 1MB
	}

	for _, opt := range opts {
		opt(options)
	}

	return &downloader{
		logger:  options.Logger,
		logFreq: options.LogFreq,
	}
}

// Download downloads a file from downloadUrl to path.
func (d *downloader) Download(path string, downloadUrl string) error {
	return d.DownloadWithContext(context.Background(), path, downloadUrl)
}

// DownloadWithContext downloads a file from downloadUrl to path. It uses the provided context to cancel the download.
func (d *downloader) DownloadWithContext(ctx context.Context, path string, downloadUrl string) error {
	downloadID := hashString(downloadUrl)
	startTime := time.Now()
	tempPath := path + ".tmp"

	logContext := map[string]interface{}{
		"downloadID": downloadID,
		"url":        downloadUrl,
		"path":       path,
		"tempPath":   tempPath,
	}

	d.logger.Debug("Starting download", logContext)

	req, err := http.NewRequestWithContext(ctx, "GET", downloadUrl, nil)
	if err != nil {
		logContext["error"] = err.Error()
		d.logger.Error("Error creating HTTP request", logContext)
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logContext["error"] = err.Error()
		d.logger.Error("Error making HTTP request", logContext)
		return err
	}
	defer resp.Body.Close()

	logContext["status"] = resp.Status
	logContext["contentLength"] = resp.ContentLength

	d.logger.Debug("HTTP request successful", logContext)

	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logContext["error"] = err.Error()
		d.logger.Error("Error opening temporary file", logContext)
		return err
	}
	defer f.Close()

	wg := &sync.WaitGroup{}
	pw := &progressWriter{
		id:      downloadID,
		w:       f,
		maxSize: resp.ContentLength,
		logFreq: d.logFreq,
		logger:  d.logger,
		logCh:   make(chan int64),
		wg:      wg,
	}
	pw.startLogging()

	// Wrap resp.Body in a contextReader
	cr := &contextReader{
		id:     downloadID,
		r:      resp.Body,
		ctx:    ctx,
		logger: d.logger,
	}

	err = d.downloadToFile(pw, cr)
	pw.stopLogging()
	wg.Wait() // Wait for the logging goroutine to finish

	if err != nil {
		logContext["error"] = err.Error()
		d.logger.Error("Error downloading file", logContext)
		return err
	}

	d.logger.Debug("Download completed to temporary file", logContext)

	// Close the file without defer so it can happen before Rename()
	f.Close()

	if err = os.Rename(tempPath, path); err != nil {
		logContext["error"] = err.Error()
		logContext["from"] = tempPath
		logContext["to"] = path
		d.logger.Error("Error renaming temporary file", logContext)
		return err
	}

	logContext["elapsedTime"] = time.Since(startTime).String()
	d.logger.Info("File successfully downloaded", logContext)

	return nil
}

// downloadToFile downloads the contents of an io.Reader to an io.Writer
func (d *downloader) downloadToFile(w io.Writer, r io.Reader) error {
	bufferSize := 1 << 20 // 1MB
	buffer := make([]byte, bufferSize)
	_, err := io.CopyBuffer(w, r, buffer)
	return err
}

// hashString returns the SHA256 hash of a string
func hashString(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
