package installation

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/driver/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/driver/runtime"
	"github.com/rocketblend/rocketblend/pkg/extractor"
	"github.com/rocketblend/rocketblend/pkg/lockfile"
)

const LockFileName = "reference.lock"

type (
	Service interface {
		Get(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack, readOnly bool) (map[reference.Reference]*Installation, error)
		Remove(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack) error
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

	options.Logger.Debug("initializing installation service", map[string]interface{}{
		"storagePath": options.StoragePath,
		"platform":    options.Platform,
	})

	err := os.MkdirAll(options.StoragePath, 0755)
	if err != nil {
		options.Logger.Error(fmt.Sprintf("unable to create directory %s", options.StoragePath), map[string]interface{}{
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

// TODO: Return a map of reference to error instead of returning the first error encountered.
func (s *service) Get(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack, readOnly bool) (map[reference.Reference]*Installation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan getResult, len(rocketPacks))
	var wg sync.WaitGroup
	wg.Add(len(rocketPacks))

	for ref, pack := range rocketPacks {
		go func(ref reference.Reference, pack *rocketpack.RocketPack) {
			defer wg.Done()

			installation, err := s.getInstallation(ctx, ref, pack, readOnly)
			if err != nil {
				cancel() // Cancel all operations as soon as an error is encountered
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
	var retErr error

	for res := range results {
		if res.error != nil {
			if retErr == nil { // Store the first error encountered
				retErr = res.error
			}
			// Continue to drain the channel
		}

		if res.inst != nil {
			installations[res.ref] = res.inst
		}
	}

	return installations, retErr
}

// TODO: Return a map of reference to error instead of returning the first error encountered.
func (s *service) Remove(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(rocketPacks))
	var wg sync.WaitGroup
	wg.Add(len(rocketPacks))

	for ref := range rocketPacks {
		go func(r reference.Reference) {
			defer wg.Done()
			if err := s.removeInstallation(ctx, r); err != nil {
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

func (s *service) getInstallation(ctx context.Context, reference reference.Reference, rocketPack *rocketpack.RocketPack, readOnly bool) (*Installation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.logger.Info("checking installation", map[string]interface{}{
		"preInstalled": rocketPack.IsPreInstalled(),
		"readOnly":     readOnly,
		"reference":    reference.String(),
	})

	var resourcePath string

	// Pre-installed rocketpacks are not downloaded as they are already available within the build.
	if !rocketPack.IsPreInstalled() {
		var source *rocketpack.Source

		if rocketPack.IsBuild() {
			source = rocketPack.Build.Sources[s.platform]
		}

		if rocketPack.IsAddon() {
			source = rocketPack.Addon.Source
		}

		if source == nil {
			return nil, fmt.Errorf("no source found for %s", reference.String())
		}

		installationPath := filepath.Join(s.storagePath, reference.String())
		resourcePath = filepath.Join(installationPath, source.Resource)

		_, err := os.Stat(resourcePath)
		if err != nil {
			if os.IsNotExist(err) && !readOnly {
				err := s.downloadInstallation(ctx, source.URI, installationPath)
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
			FilePath: resourcePath,
			ARGS:     rocketPack.Build.Args,
		}
	}

	if rocketPack.IsAddon() {
		installation.Addon = &Addon{
			FilePath: resourcePath,
			Name:     rocketPack.Addon.Name,
			Version:  rocketPack.Addon.Version,
		}
	}

	return installation, nil
}

func (s *service) downloadInstallation(ctx context.Context, downloadURI *downloader.URI, installationPath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if downloadURI == nil {
		return fmt.Errorf("no download URI provided")
	}

	// Create the installation path if it doesn't exist.
	if err := os.MkdirAll(installationPath, 0755); err != nil {
		s.logger.Error("failed to create installation path", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}

	// Lock the installation path to prevent concurrent downloads.
	locker := s.newLocker(installationPath)
	if err := locker.Lock(ctx); err != nil {
		s.logger.Error("failed to acquire lock", map[string]interface{}{"error": err, "installationPath": installationPath})
		return err
	}
	defer locker.Unlock()

	downloadedFilePath := filepath.Join(installationPath, path.Base(downloadURI.Path))
	s.logger.Info("downloading installation", map[string]interface{}{
		"downloadURI":        downloadURI.String(),
		"downloadedFilePath": downloadedFilePath,
	})

	// Download the file.
	if err := s.downloader.DownloadWithContext(ctx, downloadedFilePath, downloadURI); err != nil {
		return err
	}

	// Extract the file if it's an archive.
	if isArchive(downloadedFilePath) {
		if err := s.extractor.ExtractWithContext(ctx, downloadedFilePath, filepath.Dir(downloadedFilePath)); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) removeInstallation(ctx context.Context, reference reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	installationPath := filepath.Join(s.storagePath, reference.String())

	locker := s.newLocker(installationPath)
	if err := locker.Lock(ctx); err != nil {
		s.logger.Error("failed to acquire lock", map[string]interface{}{"error": err, "installationPath": installationPath})
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
	s.logger.Debug("creating new file lock", map[string]interface{}{"path": dir, "lockFile": LockFileName})
	return lockfile.NewLocker(filepath.Join(dir, LockFileName))
}
