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
	if reference == "" {
		return nil, fmt.Errorf("missing reference - no place to save record")
	}

	// ensure there is a resource (name) to save record as
	if resource == "" {
		return nil, fmt.Errorf("missing resource - unable to save record (no name)")
	}

	// create full path to record and check to see if file exists
	record := filepath.Join(d.dir, reference.String(), resource)
	if _, err := stat(record); err != nil {
		return nil, err
	}

	// read record from store
	b, err := os.ReadFile(record)
	if err != nil {
		return nil, err
	}

	return b, nil
}
