package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/google/uuid"
)

const (
	TempFileExtension = ".tmp"
	DownloadInfoFile  = "download-info.json"
)

type (
	DownloadInfo struct {
		URI  string `json:"uri"`
		Size int64  `json:"size"`
	}

	Downloader interface {
		Download(ctx context.Context, path string, uri *URI) error
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

	options.Logger.Debug("initializing downloader", map[string]interface{}{"logFreq": options.LogFreq})

	return &downloader{
		logger:  options.Logger,
		logFreq: options.LogFreq,
	}
}

// Download downloads a file from downloadUrl to path. It uses the provided context to cancel the download.
func (d *downloader) Download(ctx context.Context, path string, uri *URI) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	downloadID := uuid.New().String()
	tempPath := path + TempFileExtension
	infoPath := filepath.Join(filepath.Dir(path), DownloadInfoFile)

	if err := d.prepareDirectoryAndInfoFile(tempPath, infoPath, uri); err != nil {
		return err
	}
	defer os.Remove(infoPath) // Ensures removal of the info file if the download fails

	reader, contentLength, err := d.setupReader(ctx, uri, 0) // fileSize is determined inside setupReader if needed
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := d.writeToFile(ctx, tempPath, contentLength, reader); err != nil {
		return err
	}

	if err := d.finalizeDownload(tempPath, path, downloadID, uri); err != nil {
		return err
	}

	return nil
}

// prepareDirectoryAndInfoFile creates the directory for the download and the info file
func (d *downloader) prepareDirectoryAndInfoFile(tempPath, infoPath string, uri *URI) error {
	if err := os.MkdirAll(filepath.Dir(tempPath), 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	fileSize, err := d.checkFileSize(tempPath)
	if err != nil {
		return err
	}

	infoFile, err := os.Create(infoPath)
	if err != nil {
		return fmt.Errorf("error creating download info JSON file: %w", err)
	}
	defer infoFile.Close()

	info := DownloadInfo{
		URI:  uri.String(),
		Size: fileSize,
	}
	infoBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("error marshalling download info JSON: %w", err)
	}

	if _, err := infoFile.Write(infoBytes); err != nil {
		return fmt.Errorf("error writing download info JSON: %w", err)
	}

	return nil
}

// writeToFile writes the downloaded file to disk
func (d *downloader) writeToFile(ctx context.Context, path string, contentLength int64, reader io.ReadCloser) error {
	return d.openAndWriteToFile(ctx, path, contentLength, reader, uuid.New().String())
}

// finalizeDownload renames the temporary file to its final name once the download is complete
func (d *downloader) finalizeDownload(tempPath, path, downloadID string, uri *URI) error {
	if err := d.renameFile(tempPath, path); err != nil {
		return err
	}

	d.logger.Debug("file successfully downloaded", map[string]interface{}{
		"downloadID": downloadID,
		"uri":        uri.String(),
		"path":       path,
	})

	return nil
}

// checkFileSize checks if a file exists and returns its size if it does
func (d *downloader) checkFileSize(tempPath string) (int64, error) {
	var fileSize int64 = 0
	if fi, err := os.Stat(tempPath); err == nil {
		fileSize = fi.Size()
	}

	return fileSize, nil
}

// setupReader sets up an io.ReadCloser based on whether the file is local or remote
func (d *downloader) setupReader(ctx context.Context, uri *URI, fileSize int64) (io.ReadCloser, int64, error) {
	if uri.IsRemote() {
		return d.setupRemoteReader(ctx, uri, fileSize)
	}

	if uri.IsLocal() {
		return d.setupLocalReader(uri)
	}

	return nil, 0, fmt.Errorf("unknown URI type: %s", uri.String())
}

// setupRemoteReader sets up an io.ReadCloser for a remote file
func (d *downloader) setupRemoteReader(ctx context.Context, uri *URI, fileSize int64) (io.ReadCloser, int64, error) {
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

	d.logger.Debug("HTTP request successful", map[string]interface{}{
		"status":        resp.Status,
		"contentLength": contentLength,
		"resumed":       resumed,
	})

	return resp.Body, contentLength, nil
}

// setupLocalReader sets up an io.ReadCloser for a local file
func (d *downloader) setupLocalReader(uri *URI) (io.ReadCloser, int64, error) {
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

// openAndWriteToFile opens the file for writing and starts the download process
func (d *downloader) openAndWriteToFile(ctx context.Context, tempPath string, contentLength int64, reader io.ReadCloser, downloadID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	wg := &sync.WaitGroup{}
	pw := &progressWriter{
		id:      downloadID,
		w:       f,
		maxSize: contentLength,
		logFreq: d.logFreq,
		logger:  d.logger,
		logCh:   make(chan int64),
		wg:      wg,
	}
	pw.startLogging()

	// Wrap resp.Body in a contextReader
	cr := &contextReader{
		id:     downloadID,
		r:      reader,
		ctx:    ctx,
		logger: d.logger,
	}

	bufferSize := 2 << 20 // 2MB
	buffer := make([]byte, bufferSize)
	_, err = io.CopyBuffer(pw, cr, buffer)
	pw.stopLogging()
	wg.Wait() // Wait for the logging goroutine to finish

	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			d.logger.Info("download cancelled", map[string]interface{}{"downloadID": downloadID})
		}

		return err
	}

	return nil
}

// renameFile renames the temporary file to its final name once the download is complete
func (d *downloader) renameFile(tempPath string, path string) error {
	if err := os.Rename(tempPath, path); err != nil {
		return err
	}

	return nil
}
