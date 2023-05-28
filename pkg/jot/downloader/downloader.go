package downloader

import (
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
)

type (
	Downloader interface {
		Download(path string, downloadUrl string) error
	}

	Options struct {
		Logger logger.Logger
	}

	Option func(*Options)

	downloader struct {
		logger logger.Logger
	}
)

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func New(opts ...Option) Downloader {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return &downloader{
		logger: options.Logger,
	}
}

func (d *downloader) Download(path string, downloadUrl string) error {
	d.logger.Debug("Starting download", map[string]interface{}{"url": downloadUrl, "path": path})

	tempPath := path + ".tmp"
	d.logger.Debug("Temporary path for download", map[string]interface{}{"tempPath": tempPath})

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		d.logger.Error("Error creating HTTP request", map[string]interface{}{"error": err.Error()})
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		d.logger.Error("Error making HTTP request", map[string]interface{}{"error": err.Error()})
		return err
	}
	defer resp.Body.Close()

	d.logger.Debug("HTTP request successful", map[string]interface{}{"status": resp.Status, "contentLength": resp.ContentLength})

	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		d.logger.Error("Error opening temporary file", map[string]interface{}{"error": err.Error()})
		return err
	}
	defer f.Close()

	wg := &sync.WaitGroup{}
	pw := &progressWriter{
		w:       f,
		maxSize: resp.ContentLength,
		logger:  d.logger,
		logCh:   make(chan int64),
		wg:      wg,
	}
	pw.startLogging()

	err = d.downloadToFile(pw, resp.Body)
	pw.stopLogging()
	wg.Wait() // Wait for the logging goroutine to finish

	if err != nil {
		d.logger.Error("Error downloading file", map[string]interface{}{"error": err.Error()})
		return err
	}

	d.logger.Debug("Download completed to temporary file", map[string]interface{}{"tempPath": tempPath})

	// Close the file without defer so it can happen before Rename()
	f.Close()

	if err = os.Rename(tempPath, path); err != nil {
		d.logger.Error("Error renaming temporary file", map[string]interface{}{"error": err.Error(), "from": tempPath, "to": path})
		return err
	}

	d.logger.Info("File successfully downloaded", map[string]interface{}{"path": path})

	return nil
}

func (d *downloader) downloadToFile(w io.Writer, r io.Reader) error {
	bufferSize := 1 << 20 // 1MB
	buffer := make([]byte, bufferSize)
	_, err := io.CopyBuffer(w, r, buffer)
	return err
}
