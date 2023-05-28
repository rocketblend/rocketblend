package downloader

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type (
	Downloader interface {
		Download(path string, downloadUrl string) error
	}

	downloader struct {
	}
)

func New() Downloader {
	return &downloader{}
}

func (d *downloader) Download(path string, downloadUrl string) error {
	tempPath := path + ".tmp"

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	bar := progressBar(
		resp.ContentLength,
		"Downloading",
	)

	// Close the progress bar when the download is complete
	defer bar.Finish()

	bufferSize := 1 << 20 // 1MB
	buffer := make([]byte, bufferSize)
	if _, err := io.CopyBuffer(io.MultiWriter(f, bar), resp.Body, buffer); err != nil {
		f.Close()
		return err
	}

	// Close the file without defer so it can happen before Rename()
	f.Close()

	if err = os.Rename(tempPath, path); err != nil {
		return err
	}

	return nil
}

func progressBar(maxBytes int64, description ...string) *progressbar.ProgressBar {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	return progressbar.NewOptions64(
		maxBytes,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	)
}
