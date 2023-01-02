package jot

import (
	"fmt"
	"os"
	"path/filepath"
)

// Write locks the store and attempts to download the record to the store under
// the [reference] specified and with the [resource] name given
func (d *Driver) Write(reference string, resource string, downloadUrl string) error {
	// ensure there is a place to save record
	if reference == "" {
		return fmt.Errorf("missing reference - no place to save record")
	}

	// ensure there is a resource (name) to save record as
	if reference == "" {
		return fmt.Errorf("missing resource - unable to save record (no name)")
	}

	// create mutex on reference
	mutex := d.getOrCreateMutex(reference)
	mutex.Lock()
	defer mutex.Unlock()

	// create full paths to reference, final resource file, and temp file
	dir := filepath.Join(d.dir, reference, resource)

	// create reference directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// filePath := filepath.Join(dir, resource)
	err := d.downloader.Download(dir, downloadUrl)
	if err != nil {
		return err
	}

	return nil
}
