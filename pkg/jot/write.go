package jot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

func (d *Driver) Write(reference reference.Reference, resource string, downloadUrl string) error {
	return d.WriteWithContext(context.Background(), reference, resource, downloadUrl)
}

func (d *Driver) WriteWithContext(ctx context.Context, reference reference.Reference, resource string, downloadUrl string) error {
	path, err := d.writeWithContext(ctx, reference, resource, downloadUrl)
	if err != nil {
		return err
	}

	if isArchive(path) {
		err = d.extractor.Extract(path, filepath.Dir(path))
		if err != nil {
			return err
		}
	}

	return nil
}

// Write locks the store and attempts to download the record to the store under
// the [reference] specified and with the [resource] name given
func (d *Driver) writeWithContext(ctx context.Context, reference reference.Reference, resource string, downloadUrl string) (string, error) {
	// ensure there is a place to save record
	if reference.String() == "" {
		return "", fmt.Errorf("missing reference - no place to save record")
	}

	// ensure there is a resource (name) to save record as
	if resource == "" {
		return "", fmt.Errorf("missing resource - unable to save record (no name)")
	}

	// ensure there is a download url to download record from
	if downloadUrl == "" {
		return "", fmt.Errorf("missing download url - unable to save record (no url)")
	}

	// create mutex on reference
	mutex := d.getOrCreateMutex(reference.String())
	mutex.Lock()
	defer mutex.Unlock()

	// create full paths to reference, final resource file, and temp file
	dir := filepath.Join(d.storageDir, reference.String())

	// create reference directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(dir, resource)
	err := d.downloader.DownloadWithContext(ctx, filePath, downloadUrl)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
