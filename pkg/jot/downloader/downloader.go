package downloader

import (
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
}

func New() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(path string, downloadUrl string) error {
	tempPath := path + ".tmp"

	req, _ := http.NewRequest("GET", downloadUrl, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	f, _ := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
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
