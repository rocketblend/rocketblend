package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/library"
)

func (c *Client) InstallPackage(ref string) error {
	return c.installPackageIgnorable(ref, false)
}

func (c *Client) installPackageIgnorable(ref string, ignore bool) error {
	// TODO: Move downloading packages/builds into library service.

	// Check if addon already exists
	adn, _ := c.findAddon(ref)
	if adn != nil {
		if !ignore {
			return fmt.Errorf("already installed")
		}
		return nil
	}

	// Fetch package from ref
	pack, err := c.library.FetchPackage(ref)
	if err != nil {
		return err
	}

	if pack == nil {
		return fmt.Errorf("invalid package")
	}

	// Output directory
	outPath := filepath.Join(c.conf.InstallationDir, ref)

	// Create output directories
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return err
	}

	// download file path
	name := filepath.Base(pack.Source.URL)
	filePath := filepath.Join(outPath, name)

	// Download file to file path
	err = c.downloader.Download(pack.Source.URL, filePath)
	if err != nil {
		return err
	}

	// Markshal pack
	data, err := json.Marshal(pack)
	if err != nil {
		return err
	}

	// Write out package.json
	if err := os.WriteFile(filepath.Join(outPath, library.PackgeFile), data, os.ModePerm); err != nil {
		return err
	}

	// Add addon to database
	err = c.addon.Create(c.newAddon(ref, outPath))
	if err != nil {
		return err
	}

	return nil
}
