package installation

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/driver/source"
	"github.com/rocketblend/rocketblend/pkg/extractor"
	"github.com/rocketblend/rocketblend/pkg/lockfile"
)

const LockFileName = "reference.lock"

type (
	Service interface {
		GetInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack, readOnly bool) (map[reference.Reference]*Installation, error)
		RemoveInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack) error

		// TODO: Should update GetInstallations to work with sources instead of rocketpacks.
		InsertInstallations(ctx context.Context, sources map[reference.Reference]*source.Source) error
	}

	Options struct {
		Logger      logger.Logger
		StoragePath string
		Platform    runtime.Platform
		Downloader  downloader.Downloader
		Extractor   extractor.Extractor
	}

	Option func(*Options)

	getResult struct {
		ref   reference.Reference
		inst  *Installation
		error error
	}

	service struct {
		logger      logger.Logger
		storagePath string
		platform    runtime.Platform
		downloader  downloader.Downloader
		extractor   extractor.Extractor
	}
)

func WithStoragePath(storagePath string) Option {
	return func(o *Options) {
		o.StoragePath = storagePath
	}
}

func WithPlatform(platform runtime.Platform) Option {
	return func(o *Options) {
		o.Platform = platform
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithDownloader(downloader downloader.Downloader) Option {
	return func(o *Options) {
		o.Downloader = downloader
	}
}

func WithExtractor(extractor extractor.Extractor) Option {
	return func(o *Options) {
		o.Extractor = extractor
	}
}

func NewService(opts ...Option) (Service, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.StoragePath == "" {
		return nil, fmt.Errorf("storage path is required")
	}

	if options.Downloader == nil {
		downloader := downloader.New(
			downloader.WithLogger(options.Logger),
			downloader.WithLogFrequency(10<<20), // 10MB
		)

		options.Downloader = downloader
	}

	if options.Extractor == nil {
		extractor := extractor.New(
			extractor.WithLogger(options.Logger),
			extractor.WithCleanup(),
		)

		options.Extractor = extractor
	}

	options.Logger.Debug("Initializing installation service", map[string]interface{}{
		"storagePath": options.StoragePath,
		"platform":    options.Platform,
	})

	err := os.MkdirAll(options.StoragePath, 0755)
	if err != nil {
		options.Logger.Error(fmt.Sprintf("Unable to create directory %s", options.StoragePath), map[string]interface{}{
			"error": err,
		})
		return nil, err
	}

	return &service{
		logger:      options.Logger,
		storagePath: options.StoragePath,
		downloader:  options.Downloader,
		extractor:   options.Extractor,
		platform:    options.Platform,
	}, nil
}

func (s *service) GetInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack, readOnly bool) (map[reference.Reference]*Installation, error) {
	results := make(chan getResult, len(rocketPacks))

	var wg sync.WaitGroup
	wg.Add(len(rocketPacks))

	for ref, pack := range rocketPacks {
		go func(ref reference.Reference, pack *rocketpack.RocketPack) {
			defer wg.Done()
			installation, err := s.getInstallation(ctx, ref, pack, readOnly)
			if err != nil {
				s.logger.Error("Failed to get installation", map[string]interface{}{
					"error":     err,
					"reference": ref.String(),
				})
				results <- getResult{ref: ref, error: err}
				return
			}

			results <- getResult{ref: ref, inst: installation}
		}(ref, pack)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	installations := make(map[reference.Reference]*Installation, len(rocketPacks))

	for res := range results {
		if res.error != nil {
			return nil, res.error
		}

		if res.inst != nil {
			installations[res.ref] = res.inst
		}
	}

	return installations, nil
}

func (s *service) RemoveInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack) error {
	errs := make(chan error, len(rocketPacks))
	var wg sync.WaitGroup
	wg.Add(len(rocketPacks))

	for ref := range rocketPacks {
		go func(r reference.Reference) {
			defer wg.Done()
			err := s.removeInstallation(ctx, r)
			if err != nil {
				s.logger.Error("Failed to remove installation", map[string]interface{}{
					"error":     err,
					"reference": r.String(),
				})
				errs <- fmt.Errorf("failed to remove installation for %s: %w", r.String(), err)
				return
			}
		}(ref)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %v", <-errs) // return first error for simplicity
	}

	return nil
}

// InsertInstallations inserts the given rocketpacks into the service's storage. This is used for local packages.
func (s *service) InsertInstallations(ctx context.Context, sources map[reference.Reference]*source.Source) error {
	errs := make(chan error, len(sources))
	var wg sync.WaitGroup
	wg.Add(len(sources))

	for ref, src := range sources {
		go func(r reference.Reference, src *source.Source) {
			defer wg.Done()
			err := s.insertInstallation(ctx, r, src)
			if err != nil {
				s.logger.Error("Failed to insert installation", map[string]interface{}{
					"error":     err,
					"reference": r.String(),
				})
				errs <- fmt.Errorf("failed to insert installation for %s: %w", r.String(), err)
				return
			}
		}(ref, src)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %v", <-errs) // return first error for simplicity
	}

	return nil
}

func (s *service) insertInstallation(ctx context.Context, reference reference.Reference, source *source.Source) error {
	s.logger.Info("Inserting installation", map[string]interface{}{
		"reference": reference.String(),
		"source":    source.String(),
	})

	// TODO: Implement remote source insertion. Move download logic here.
	if !source.IsLocal() {
		return fmt.Errorf("remote sources are not supported")
	}

	// Check if the source file exists.
	_, err := os.Stat(source.URI.Path)
	if err != nil {
		return err
	}

	installationPath := filepath.Join(s.storagePath, reference.String())

	// Lock the installation path to prevent concurrent usage.
	locker := s.newLocker(installationPath)
	if err := locker.Lock(ctx); err != nil {
		s.logger.Error("Failed to acquire lock", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}
	defer locker.Unlock()

	// Create the installation directory if it doesn't exist.
	err = os.MkdirAll(installationPath, 0755)
	if err != nil {
		return err
	}

	// Move the file into the installation directory.
	err = s.moveFile(source.URI.Path, filepath.Join(installationPath, source.FileName))
	if err != nil {
		return err
	}

	return nil
}

func (s *service) getInstallation(ctx context.Context, reference reference.Reference, rocketPack *rocketpack.RocketPack, readOnly bool) (*Installation, error) {
	s.logger.Info("Checking installation", map[string]interface{}{
		"preInstalled": rocketPack.IsPreInstalled(),
		"readOnly":     readOnly,
		"reference":    reference.String(),
	})

	var executablePath string

	// Pre-installed rocketpacks are not downloaded as they are already available within the build.
	if !rocketPack.IsPreInstalled() {
		executableName, err := rocketPack.GetExecutableName(s.platform)
		if err != nil {
			return nil, err
		}

		installationPath := filepath.Join(s.storagePath, reference.String())
		executablePath = filepath.Join(installationPath, executableName)

		_, err = os.Stat(executablePath)
		if err != nil {
			if os.IsNotExist(err) && !readOnly {
				downloadUrl, err := rocketPack.GetDownloadUrl(s.platform)
				if err != nil {
					return nil, err
				}

				err = s.downloadInstallation(ctx, downloadUrl, installationPath)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	installation := &Installation{}
	if rocketPack.IsBuild() {
		installation.Build = &Build{
			FilePath: executablePath,
			ARGS:     rocketPack.Build.Args,
		}
	}

	if rocketPack.IsAddon() {
		installation.Addon = &Addon{
			FilePath: executablePath,
			Name:     rocketPack.Addon.Name,
			Version:  rocketPack.Addon.Version,
		}
	}

	return installation, nil
}

func (s *service) downloadInstallation(ctx context.Context, downloadUrl string, installationPath string) error {
	if downloadUrl == "" {
		return nil // No download URL, nothing to do.
	}

	fileName, err := getFilenameFromURL(downloadUrl)
	if err != nil {
		return err
	}

	// Create the installation path if it doesn't exist.
	err = os.MkdirAll(installationPath, 0755)
	if err != nil {
		s.logger.Error("Failed to create installation path", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}

	// Lock the installation path to prevent concurrent downloads.
	locker := s.newLocker(installationPath)
	if err := locker.Lock(ctx); err != nil {
		s.logger.Error("Failed to acquire lock", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}
	defer locker.Unlock()

	downloadedFilePath := filepath.Join(installationPath, fileName)
	s.logger.Info("Downloading installation", map[string]interface{}{
		"downloadUrl":        downloadUrl,
		"downloadedFilePath": downloadedFilePath,
		"installationPath":   installationPath,
	})

	// Download the file.
	err = s.downloader.DownloadWithContext(ctx, downloadedFilePath, downloadUrl)
	if err != nil {
		return err
	}

	// Extract the file if it's an archive.
	if isArchive(downloadedFilePath) {
		err = s.extractor.ExtractWithContext(ctx, downloadedFilePath, filepath.Dir(downloadedFilePath))
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *service) removeInstallation(ctx context.Context, reference reference.Reference) error {
	installationPath := filepath.Join(s.storagePath, reference.String())

	locker := s.newLocker(installationPath)
	if err := locker.Lock(ctx); err != nil {
		s.logger.Error("Failed to acquire lock", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}
	defer locker.Unlock()

	err := os.RemoveAll(installationPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) newLocker(dir string) *lockfile.Locker {
	s.logger.Debug("Creating new file lock", map[string]interface{}{"path": dir, "lockFile": LockFileName})
	return lockfile.NewLocker(filepath.Join(dir, LockFileName))
}

func (s *service) moveFile(src, dst string) error {
	// Check if source file exists
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if destination file exists
	_, err = os.Stat(dst)
	if err == nil {
		// Destination file exists, check if it's the same as the source file
		dstStat, err := os.Stat(dst)
		if err != nil {
			return err
		}

		// Check if source and destination files are the same
		if os.SameFile(srcStat, dstStat) {
			// Source and destination files are the same, nothing to do
			return nil
		}
	}

	// Use os.Rename to move the file
	err = os.Rename(src, dst)
	if err != nil {
		return err
	}

	return nil
}

func getFilenameFromURL(downloadURL string) (string, error) {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	return path.Base(u.Path), nil
}
