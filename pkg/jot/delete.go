package jot

import (
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

func (d *Driver) DeleteAll(reference reference.Reference) error {
	// create mutex on reference
	mutex := d.getOrCreateMutex(reference.String())
	mutex.Lock()
	defer mutex.Unlock()

	path := filepath.Join(d.storagePath, reference.String())

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
