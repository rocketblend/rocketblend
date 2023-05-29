package jot

import (
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

func (d *Driver) DeleteAll(reference reference.Reference) error {
	mutex := d.getOrCreateMutex(reference.String())
	mutex.Lock()
	defer mutex.Unlock()

	path := filepath.Join(d.storagePath, reference.String())

	d.logger.Debug("Starting to delete all files", map[string]interface{}{"path": path})

	err := os.RemoveAll(path)
	if err != nil {
		d.logger.Error("Failed to delete all files", map[string]interface{}{"error": err.Error(), "path": path})
		return err
	}

	d.logger.Info("Successfully deleted all files", map[string]interface{}{"path": path})

	return nil
}
