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
	d.logger.Debug("Starting write", map[string]interface{}{"reference": reference, "resource": resource, "downloadUrl": downloadUrl})
	path, err := d.writeWithContext(ctx, reference, resource, downloadUrl)
	if err != nil {
		d.logger.Error("Error writing", map[string]interface{}{"error": err.Error(), "reference": reference, "resource": resource, "downloadUrl": downloadUrl})
		return err
	}

	if isArchive(path) {
		extractionPath := filepath.Dir(path)
		d.logger.Debug("Starting archive extraction", map[string]interface{}{"path": path, "extractionPath": extractionPath})
		err = d.extractor.ExtractWithContext(ctx, path, extractionPath)
		if err != nil {
			d.logger.Error("Error during archive extraction", map[string]interface{}{"error": err.Error(), "path": path, "extractionPath": extractionPath})
			return err
		}
		d.logger.Debug("Finished archive extraction", map[string]interface{}{"path": path, "extractionPath": extractionPath})
	}

	d.logger.Debug("Finished write", map[string]interface{}{"path": path, "reference": reference, "resource": resource, "downloadUrl": downloadUrl})

	return nil
}

// Write locks the store and attempts to download the record to the store under
// the [reference] specified and with the [resource] name given
func (d *Driver) writeWithContext(ctx context.Context, reference reference.Reference, resource string, downloadUrl string) (string, error) {
	// ensure there is a place to save record
	if reference.String() == "" {
		d.logger.Warn("Missing reference - no place to save record")
		return "", fmt.Errorf("missing reference - no place to save record")
	}

	// ensure there is a resource (name) to save record as
	if resource == "" {
		d.logger.Warn("Missing resource - unable to save record (no name)")
		return "", fmt.Errorf("missing resource - unable to save record (no name)")
	}

	// ensure there is a download url to download record from
	if downloadUrl == "" {
		d.logger.Warn("Missing download URL")
		return "", fmt.Errorf("missing download url - unable to save record (no url)")
	}

	// create mutex on reference
	mutex := d.getOrCreateMutex(reference.String())
	mutex.Lock()
	defer mutex.Unlock()

	// create full paths to reference, final resource file, and temp file
	dir := filepath.Join(d.storagePath, reference.String())

	// create reference directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		d.logger.Error("Error creating directory", map[string]interface{}{"error": err.Error(), "directory": dir})
		return "", err
	}

	filePath := filepath.Join(dir, resource)
	err := d.downloader.DownloadWithContext(ctx, filePath, downloadUrl)
	if err != nil {
		d.logger.Error("Error during download", map[string]interface{}{"error": err.Error(), "filePath": filePath, "downloadUrl": downloadUrl})
		return "", err
	}

	return filePath, nil
}
