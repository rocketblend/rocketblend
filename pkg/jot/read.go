package jot

import (
	"fmt"
	"os"
	"path/filepath"
)

// Read a record from the database
func (d *Driver) Read(reference, resource string, b []byte) error {
	// ensure there is a place to save record
	if reference == "" {
		return fmt.Errorf("missing reference - no place to save record")
	}

	// ensure there is a resource (name) to save record as
	if resource == "" {
		return fmt.Errorf("missing resource - unable to save record (no name)")
	}

	// create full path to record and check to see if file exists
	record := filepath.Join(d.dir, reference, resource)
	if _, err := stat(record); err != nil {
		return err
	}

	// read record from store
	b, err := os.ReadFile(record)
	if err != nil {
		return err
	}

	return nil
}
