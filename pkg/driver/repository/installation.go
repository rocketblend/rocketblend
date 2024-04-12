package repository

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/lockfile"
	"github.com/rocketblend/rocketblend/pkg/types"
)

const (
	LockFileName = "reference.lock"
)

type (
	getInstallationResult struct {
		ref   reference.Reference
		inst  *types.Installation
		error error
	}
)

func (r *repository) GetInstallations(ctx context.Context, opts *types.GetInstallationOpts) (*types.GetInstallationResult, error) {
	installations, err := r.getInstallations(ctx, opts.References, opts.Fetch)
	if err != nil {
		return nil, err
	}

	return &types.GetInstallationResult{
		Installations: installations,
	}, nil
}

func (r *repository) RemoveInstallations(ctx context.Context, opts *types.RemoveInstallationOpts) error {
	if err := r.removeInstallations(ctx, opts.References); err != nil {
		return err
	}

	return nil
}

// TODO: Return a map of reference to error instead of returning the first error encountered.
func (r *repository) getInstallations(ctx context.Context, references []reference.Reference, fetch bool) ([]*types.Installation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rocketPacks, err := r.getPackages(ctx, references, 0, false)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan getInstallationResult, len(rocketPacks))
	var wg sync.WaitGroup
	wg.Add(len(rocketPacks))

	for ref, pack := range rocketPacks {
		go func(ref reference.Reference, pack *types.RocketPack) {
			defer wg.Done()

			installation, err := r.getInstallation(ctx, ref, pack, fetch)
			if err != nil {
				cancel() // Cancel all operations as soon as an error is encountered
				results <- getInstallationResult{ref: ref, error: err}
				return
			}
			results <- getInstallationResult{ref: ref, inst: installation}
		}(ref, pack)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	installations := make([]*types.Installation, 0, len(rocketPacks))
	var retErr error

	for res := range results {
		if res.error != nil {
			if retErr == nil { // Store the first error encountered
				retErr = res.error
			}
			// Continue to drain the channel
		}

		if res.inst != nil {
			installations = append(installations, res.inst)
		}
	}

	return installations, retErr
}

func (r *repository) getInstallation(ctx context.Context, reference reference.Reference, rocketPack *types.RocketPack, fetch bool) (*types.Installation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.logger.Info("checking installation", map[string]interface{}{
		"bundled":   rocketPack.Bundled(),
		"reference": reference.String(),
		"fetch":     fetch,
	})

	var resourcePath string

	// Bundled rocketpacks are not downloaded as they are already available within the build.
	if !rocketPack.Bundled() {
		var source *rocketpack.Source

		if rocketPack.IsBuild() {
			source = rocketPack.Build.Sources[r.platform]
		}

		if rocketpack.IsAddon() {
			source = rocketPack.Addon.Source
		}

		if source == nil {
			return nil, fmt.Errorf("no source found for %s", reference.String())
		}

		installationPath := filepath.Join(r.installationPath, reference.String())
		resourcePath = filepath.Join(installationPath, source.Resource)

		_, err := os.Stat(resourcePath)
		if err != nil {
			if os.IsNotExist(err) && !readOnly {
				err := r.downloadInstallation(ctx, source.URI, installationPath)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	return &types.Installation{
		Reference: reference,
		Type:      rocketPack.Type,
		Path:      resourcePath,
		Name:      rocketPack.Name,
		Version:   rocketPack.Version,
	}, nil
}

// TODO: Return a map of reference to error instead of returning the first error encountered.
func (r *repository) removeInstallations(ctx context.Context, references []reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	rocketPacks, err := r.getPackages(ctx, references, 0, false)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(rocketPacks))
	var wg sync.WaitGroup
	wg.Add(len(rocketPacks))

	for ref := range rocketPacks {
		go func(ref reference.Reference) {
			defer wg.Done()
			if err := r.removeInstallation(ctx, ref); err != nil {
				cancel()
				errs <- err
				return
			}
		}(ref)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	var retErr error
	for err := range errs {
		if err != nil {
			if retErr == nil {
				retErr = err
			}
		}
	}

	return retErr
}

func (r *repository) downloadInstallation(ctx context.Context, downloadURI *types.URI, installationPath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if downloadURI == nil {
		return fmt.Errorf("no download URI provided")
	}

	// Create the installation path if it doesn't exist.
	if err := os.MkdirAll(installationPath, 0755); err != nil {
		r.logger.Error("failed to create installation path", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}

	// Lock the installation path to prevent concurrent downloads.
	cancel, err := r.lock(ctx, installationPath)
	if err != nil {
		return err
	}
	defer cancel()

	downloadedFilePath := filepath.Join(installationPath, path.Base(downloadURI.Path))
	r.logger.Info("downloading installation", map[string]interface{}{
		"downloadURI":        downloadURI.String(),
		"downloadedFilePath": downloadedFilePath,
	})

	// Download the file.
	if err := r.downloader.Download(ctx, &types.DownloadOpts{
		URI:  downloadURI,
		Path: downloadedFilePath,
	}); err != nil {
		return err
	}

	// Extract the file if it's an archive.
	if isArchive(downloadedFilePath) {
		if err := r.extractor.Extract(ctx, &types.ExtractOpts{
			Path:       downloadedFilePath,
			OutputPath: filepath.Dir(downloadedFilePath),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) removeInstallation(ctx context.Context, reference reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	installationPath := filepath.Join(r.installationPath, reference.String())

	cancel, err := r.lock(ctx, installationPath)
	if err != nil {
		return err
	}
	defer cancel()

	if err := os.RemoveAll(installationPath); err != nil {
		return err
	}

	return nil
}

func (r *repository) lock(ctx context.Context, dir string) (cancelFunc func(), err error) {
	r.logger.Debug("creating new file lock", map[string]interface{}{"path": dir, "lockFile": LockFileName})
	return lockfile.New(ctx, lockfile.WithPath(filepath.Join(dir, LockFileName)), lockfile.WithLogger(r.logger))
}
