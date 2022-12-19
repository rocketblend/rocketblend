package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/library"
)

func (c *Client) InstallBuild(ref string) error {
	// Check if install already exists
	inst, _ := c.findInstall(ref)
	if inst != nil {
		return fmt.Errorf("already installed")
	}

	// Fetch build from ref
	build, err := c.library.FetchBuild(ref)
	if err != nil {
		return err
	}

	if build == nil {
		return fmt.Errorf("invalid build")
	}

	// Output directory
	outPath := filepath.Join(c.conf.InstallationDir, ref)

	// Create output directories
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return err
	}

	// build info for current platform
	source := build.GetSourceForPlatform(c.conf.Platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", c.conf.Platform)
	}

	// Download URL
	downloadURL := source.URL

	// Download file path
	name := filepath.Base(downloadURL)
	filePath := filepath.Join(outPath, name)

	// Download file to file path
	err = c.downloader.Download(downloadURL, filePath)
	if err != nil {
		return err
	}

	// Extract the archived file
	if err := c.archiver.Extract(filePath); err != nil {
		return err
	}

	// Markshal build
	data, err := json.Marshal(build)
	if err != nil {
		return err
	}

	// Write out build.json
	if err := os.WriteFile(filepath.Join(outPath, library.BuildFile), data, os.ModePerm); err != nil {
		return err
	}

	// Add install to database
	err = c.install.Create(c.newInstall(ref, outPath, build.Packages))
	if err != nil {
		return err
	}

	// TODO: call asynchronously
	// Install packages
	for _, p := range build.Packages {
		err = c.installPackageIgnorable(p, true)
		if err != nil {
			return err
		}
	}

	return nil
}
