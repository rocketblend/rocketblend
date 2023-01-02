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

	path := filepath.Join(c.conf.InstallationDir, ref)
	_, err := c.library.FindPackageByPath(path)
	if err == nil {
		if !ignore {
			return fmt.Errorf("already installed")
		}
		return nil
	}

	// Fetch package from ref
	pack, err := c.library.FetchPackage(path)
	if err != nil {
		return err
	}

	if pack == nil {
		return fmt.Errorf("invalid package")
	}

	// Create output directories
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// download file path
	name := filepath.Base(pack.Source.URL)
	filePath := filepath.Join(path, name)

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
	if err := os.WriteFile(filepath.Join(path, library.PackgeFile), data, os.ModePerm); err != nil {
		return err
	}

	return nil
}
