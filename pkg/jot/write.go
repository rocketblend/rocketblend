package jot

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

func (d *Driver) Write(reference reference.Reference, resource string, downloadUrl string) error {
	_, err := d.write(reference, resource, downloadUrl)
	return err
}

func (d *Driver) WriteAndExtract(reference reference.Reference, resource string, downloadUrl string) error {
	path, err := d.write(reference, resource, downloadUrl)
	if err != nil {
		return err
	}

	err = d.extractor.Extract(path, filepath.Dir(path))
	if err != nil {
		return err
	}

	return nil
}

// Write locks the store and attempts to download the record to the store under
// the [reference] specified and with the [resource] name given
func (d *Driver) write(reference reference.Reference, resource string, downloadUrl string) (string, error) {
	// ensure there is a place to save record
	if reference.String() == "" {
		return "", fmt.Errorf("missing reference - no place to save record")
	}

	// ensure there is a resource (name) to save record as
	if reference.String() == "" {
		return "", fmt.Errorf("missing resource - unable to save record (no name)")
	}

	// create mutex on reference
	mutex := d.getOrCreateMutex(reference.String())
	mutex.Lock()
	defer mutex.Unlock()

	// create full paths to reference, final resource file, and temp file
	dir := filepath.Join(d.dir, reference.String())

	// create reference directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(dir, resource)
	err := d.downloader.Download(filePath, downloadUrl)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
