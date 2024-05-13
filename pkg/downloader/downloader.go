package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/types"
)

const (
	TempFileExtension = ".tmp"
	// DownloadInfoFile  = "download-info.json"
)

type (
	// DownloadInfo struct {
	// 	URI  string `json:"uri"`
	// 	Size int64  `json:"size"`
	// }

	Options struct {
		Logger         logger.Logger
		BufferSize     int
		UpdateInterval time.Duration
	}

	Option func(*Options)

	Downloader struct {
		logger         logger.Logger
		bufferSize     int
		updateInterval time.Duration
	}
)

// With Logger sets the logger to use. The default is no-op.
func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithBufferSize sets the buffer size for reading and writing. The default is 2MB.
func WithBufferSize(bufferSize int) Option {
	return func(o *Options) {
		o.BufferSize = bufferSize
	}
}

// WithUpdateInterval sets the frequency at which progress updates are sent. The default is 5 seconds.
func WithUpdateInterval(updateInterval time.Duration) Option {
	return func(o *Options) {
		o.UpdateInterval = updateInterval
	}
}

// New creates a new Downloader.
func New(opts ...Option) (*Downloader, error) {
	options := &Options{
		Logger:         logger.NoOp(),
		BufferSize:     2 << 20,         // Default buffer size is 2MB
		UpdateInterval: 5 * time.Second, // Default update interval is 5 seconds
	}

	for _, opt := range opts {
		opt(options)
	}

	options.Logger.Debug("initializing Downloader", map[string]interface{}{
		"bufferSize":      options.BufferSize,
		"updateFrequency": options.UpdateInterval,
	})

	return &Downloader{
		logger:         options.Logger,
		bufferSize:     options.BufferSize,
		updateInterval: options.UpdateInterval,
	}, nil
}

// DownloadOpts contains the options for downloading a file
func (d *Downloader) Download(ctx context.Context, opts *types.DownloadOpts) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	tempPath := opts.Path + TempFileExtension

	if err := os.MkdirAll(filepath.Dir(tempPath), 0755); err != nil {
		return err
	}

	fileSize := d.checkFileSize(tempPath)
	reader, contentLength, err := d.setupReader(ctx, opts.URI, fileSize)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := d.writeToFile(ctx, tempPath, fileSize, contentLength, reader, opts.ProgressChan); err != nil {
		return err
	}

	if err := os.Rename(tempPath, opts.Path); err != nil {
		return err
	}

	d.logger.Debug("file successfully downloaded", map[string]interface{}{
		"uri":  opts.URI.String(),
		"path": opts.Path,
	})

	return nil
}

// checkFileSize checks the size of the file on disk, returning 0 if it doesn't exist
func (d *Downloader) checkFileSize(tempPath string) int64 {
	var fileSize int64 = 0
	if fi, err := os.Stat(tempPath); err == nil {
		fileSize = fi.Size()
	}

	return fileSize
}

// writeToFile writes the file to disk, updating progress as it goes
func (d *Downloader) writeToFile(ctx context.Context, path string, initialSize int64, contentLength int64, reader io.ReadCloser, progress chan<- types.Progress) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	totalBytes := initialSize
	lastUpdateBytes := initialSize
	lastTime := time.Now()
	buffer := make([]byte, d.bufferSize)
	ticker := time.NewTicker(d.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			lastTime = d.reportProgress(path, lastUpdateBytes, totalBytes, contentLength, progress, lastTime)
			lastUpdateBytes = totalBytes
		default:
			totalBytes, err = d.processRead(reader, f, buffer, totalBytes)
			if err != nil {
				if err == io.EOF {
					d.reportProgress(path, lastUpdateBytes, totalBytes, contentLength, progress, lastTime)
					return nil
				}

				return err
			}
		}
	}
}

// reportProgress sends progress updates to the channel
func (d *Downloader) reportProgress(path string, lastUpdateBytes int64, totalBytes int64, totalSize int64, progress chan<- types.Progress, lastTime time.Time) time.Time {
	now := time.Now()
	timeElapsed := now.Sub(lastTime).Seconds()
	speed := float64(totalBytes-lastUpdateBytes) / timeElapsed

	if progress != nil {
		progress <- types.Progress{
			BytesRead: totalBytes,
			TotalSize: totalSize,
			Speed:     speed,
		}
	}

	d.logger.Info("download progress", map[string]interface{}{
		"path":      path,
		"bytesRead": totalBytes,
		"totalSize": totalSize,
		"speed":     speed,
	})

	return now
}

// processRead reads from the reader and writes to the file
func (d *Downloader) processRead(reader io.Reader, f *os.File, buffer []byte, totalBytes int64) (int64, error) {
	bytesRead, readErr := reader.Read(buffer)
	if bytesRead > 0 {
		if _, writeErr := f.Write(buffer[:bytesRead]); writeErr != nil {
			return totalBytes, writeErr
		}
		totalBytes += int64(bytesRead)
	}

	return totalBytes, readErr
}

// setupReader sets up an io.ReadCloser based on whether the file is local or remote
func (d *Downloader) setupReader(ctx context.Context, uri *types.URI, fileSize int64) (io.ReadCloser, int64, error) {
	if uri.IsRemote() {
		return d.setupRemoteReader(ctx, uri, fileSize)
	}

	if uri.IsLocal() {
		return d.setupLocalReader(uri)
	}

	return nil, 0, fmt.Errorf("unknown URI type: %s", uri.String())
}

// setupRemoteReader sets up an io.ReadCloser for a remote file
func (d *Downloader) setupRemoteReader(ctx context.Context, uri *types.URI, fileSize int64) (io.ReadCloser, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", fileSize))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		err = fmt.Errorf("received non 200/206 status code: %s", resp.Status)
		d.logger.Error("received non 200/206 status code", map[string]interface{}{"err": err.Error()})
		return nil, 0, err
	}

	contentLength := resp.ContentLength
	var resumed bool
	if resp.StatusCode == 206 && fileSize > 0 {
		resumed = true
	}

	d.logger.Debug("http request successful", map[string]interface{}{
		"status":        resp.Status,
		"contentLength": contentLength,
		"resumed":       resumed,
	})

	return resp.Body, contentLength, nil
}

// setupLocalReader sets up an io.ReadCloser for a local file
func (d *Downloader) setupLocalReader(uri *types.URI) (io.ReadCloser, int64, error) {
	file, err := os.Open(uri.Path)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening local file: %w", err)
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	return file, fi.Size(), nil
}
