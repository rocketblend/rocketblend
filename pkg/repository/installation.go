package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/lockfile"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

const (
	LockFileName     = "reference.lock"
	ProgressFileName = "download.json"
)

type (
	getInstallationResult struct {
		reference    reference.Reference
		installation *types.Installation
	}
)

func (r *Repository) GetInstallations(ctx context.Context, opts *types.GetInstallationsOpts) (*types.GetInstallationsResult, error) {
	if err := r.validator.Validate(opts); err != nil {
		return nil, err
	}

	installations, err := r.getInstallations(ctx, opts.Dependencies, opts.Fetch)
	if err != nil {
		return nil, err
	}

	return &types.GetInstallationsResult{
		Installations: installations,
	}, nil
}

func (r *Repository) RemoveInstallations(ctx context.Context, opts *types.RemoveInstallationsOpts) error {
	if err := r.validator.Validate(opts); err != nil {
		return err
	}

	if err := r.removeInstallations(ctx, opts.References); err != nil {
		return err
	}

	return nil
}

func (r *Repository) getInstallations(ctx context.Context, dependencies []*types.Dependency, fetch bool) (map[reference.Reference]*types.Installation, error) {
	references := make([]reference.Reference, 0, len(dependencies))
	for _, dep := range dependencies {
		references = append(references, dep.Reference)
	}

	packs, err := r.getPackages(ctx, references, false)
	if err != nil {
		return nil, err
	}

	for _, dep := range dependencies {
		if pack, ok := packs[dep.Reference]; ok {
			if pack.Type != dep.Type {
				return nil, fmt.Errorf("dependency type mismatch: %s", dep.Reference.String())
			}
		} else {
			return nil, fmt.Errorf("dependency not found: %s", dep.Reference.String())
		}
	}

	tasks := make([]taskrunner.Task[*getInstallationResult], 0, len(packs))
	for ref, pack := range packs {
		tasks = append(tasks, func(ctx context.Context) (*getInstallationResult, error) {
			installation, err := r.getInstallation(ctx, ref, pack, fetch)
			if err != nil {
				return nil, err
			}

			return &getInstallationResult{reference: ref, installation: installation}, nil
		})
	}
	results, err := taskrunner.Run(ctx, &taskrunner.RunOpts[*getInstallationResult]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	})
	if err != nil {
		return nil, err
	}

	installations := make(map[reference.Reference]*types.Installation, len(results))
	for _, res := range results {
		installations[res.reference] = res.installation
	}

	return installations, nil
}

func (r *Repository) getInstallation(ctx context.Context, reference reference.Reference, pack *types.Package, fetch bool) (*types.Installation, error) {
	r.logger.Info("checking installation", map[string]interface{}{
		"bundled":   pack.Bundled(),
		"reference": reference.String(),
		"fetch":     fetch,
	})

	var resourcePath string

	// Bundled packages are not downloaded as they are already available within the build.
	if !pack.Bundled() {
		// TODO: Clean up this platform stuff.
		source := pack.Source(types.Platform(r.platform.String()))

		installationPath := filepath.Join(r.installationPath, reference.String())
		resourcePath = filepath.Join(installationPath, source.Resource)

		_, err := os.Stat(resourcePath)
		if err != nil {
			if os.IsNotExist(err) && fetch {
				packageFilePath := filepath.Join(r.packagePath, reference.String(), types.PackageFileName)
				err := r.downloadInstallation(ctx, source.URI, packageFilePath, installationPath)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	return &types.Installation{
		Type:    pack.Type,
		Path:    resourcePath,
		Name:    pack.Name,
		Version: pack.Version,
	}, nil
}

func (r *Repository) removeInstallations(ctx context.Context, references []reference.Reference) error {
	packs, err := r.getPackages(ctx, references, false)
	if err != nil {
		return err
	}

	tasks := make([]taskrunner.Task[struct{}], 0, len(packs))
	for ref := range packs {
		tasks = append(tasks, func(ctx context.Context) (struct{}, error) {
			return struct{}{}, r.removeInstallation(ctx, ref)
		})
	}
	_, err = taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) downloadInstallation(ctx context.Context, downloadURI *types.URI, packageFilePath string, installationPath string) error {
	if downloadURI == nil {
		return fmt.Errorf("no download URI provided")
	}

	if err := os.MkdirAll(installationPath, 0755); err != nil {
		r.logger.Error("failed to create installation path", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}

	cancel, err := r.lock(ctx, installationPath)
	if err != nil {
		return err
	}
	defer cancel()

	downloadedFilePath := filepath.Join(installationPath, path.Base(downloadURI.Path))
	r.logger.Info("downloading installation", map[string]interface{}{
		"uri":      downloadURI.String(),
		"filePath": downloadedFilePath,
	})

	progressFilePath := filepath.Join(installationPath, ProgressFileName)
	defer os.Remove(progressFilePath)

	progressChan := make(chan types.Progress)
	go func() {
		for p := range progressChan {
			if err := writeProgressFile(progressFilePath, p); err != nil {
				r.logger.Error("failed to write progress file", map[string]interface{}{
					"error": err,
					"path":  progressFilePath,
				})
			}

			if err := touchFile(packageFilePath); err != nil {
				r.logger.Error("failed to touch package", map[string]interface{}{
					"error": err,
					"path":  packageFilePath,
				})
			}
		}
	}()

	if err := r.downloader.Download(ctx, &types.DownloadOpts{
		URI:          downloadURI,
		Path:         downloadedFilePath,
		ProgressChan: progressChan,
	}); err != nil {
		return err
	}

	if helpers.IsSupportedArchive(downloadedFilePath) {
		if err := r.extractor.Extract(ctx, &types.ExtractOpts{
			Path:       downloadedFilePath,
			OutputPath: filepath.Dir(downloadedFilePath),
		}); err != nil {
			return err
		}
	}

	if err := touchFile(packageFilePath); err != nil {
		r.logger.Error("failed to touch package", map[string]interface{}{
			"error": err,
			"path":  packageFilePath,
		})
	}

	return nil
}

func (r *Repository) removeInstallation(ctx context.Context, reference reference.Reference) error {
	installationPath := filepath.Join(r.installationPath, reference.String())

	cancel, err := r.lock(ctx, installationPath)
	if err != nil {
		return err
	}
	defer cancel()

	if err := os.RemoveAll(installationPath); err != nil {
		return err
	}

	packageFilePath := filepath.Join(r.packagePath, reference.String(), types.PackageFileName)
	if err := touchFile(packageFilePath); err != nil {
		r.logger.Error("failed to touch package", map[string]interface{}{
			"error":     err,
			"reference": reference.String(),
			"path":      packageFilePath,
		})
	}

	return nil
}

func (r *Repository) lock(ctx context.Context, dir string) (cancelFunc func(), err error) {
	r.logger.Debug("creating new file lock", map[string]interface{}{"path": dir, "lockFile": LockFileName})
	return lockfile.New(ctx, lockfile.WithPath(filepath.Join(dir, LockFileName)), lockfile.WithLogger(r.logger))
}

func writeProgressFile(path string, data types.Progress) error {
	infoFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer infoFile.Close()

	infoBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := infoFile.Write(infoBytes); err != nil {
		return err
	}

	return nil
}

// touchFile used to trigger a file change event on the file system.
func touchFile(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	mode := fileInfo.Mode()
	if err := os.Chmod(path, mode^0200); err != nil {
		return err
	}

	return os.Chmod(path, mode)
}
