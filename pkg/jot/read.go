package jot

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

// Read a record from the store
func (d *Driver) Read(reference reference.Reference, resource string) ([]byte, error) {
	// ensure there is a place to save record
	if reference.String() == "" {
		d.logger.Warn("Missing reference, no place to save record")
		return nil, fmt.Errorf("missing reference - no place to save record")
	}

	// ensure there is a resource (name) to save record as
	if resource == "" {
		d.logger.Warn("Missing resource, unable to save record (no name)")
		return nil, fmt.Errorf("missing resource - unable to save record (no name)")
	}

	// create full path to record and check to see if file exists
	record := filepath.Join(d.storagePath, reference.String(), resource)
	if _, err := stat(record); err != nil {
		d.logger.Error("Failed to stat record", map[string]interface{}{"error": err.Error(), "record": record})
		return nil, err
	}

	// read record from store
	b, err := os.ReadFile(record)
	if err != nil {
		d.logger.Error("Failed to read record", map[string]interface{}{"error": err.Error(), "record": record})
		return nil, err
	}

	d.logger.Debug("Successfully read record", map[string]interface{}{"record": record})
	return b, nil
}
