package installation

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/downloader"
	"github.com/rocketblend/rocketblend/pkg/extractor"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
)

type (
	Service interface {
		GetInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack, readOnly bool) (map[reference.Reference]*Installation, error)
		RemoveInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack) error
	}

	Options struct {
		Logger      logger.Logger
		StoragePath string
		Platform    runtime.Platform
		Downloader  downloader.Downloader
		Extractor   extractor.Extractor
	}

	Option func(*Options)

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
	installations := make(map[reference.Reference]*Installation, len(rocketPacks))

	for ref, pack := range rocketPacks {
		installation, err := s.getInstallation(ctx, ref, pack, readOnly)
		if err != nil {
			return nil, err
		}

		installations[ref] = installation
	}

	return installations, nil
}

func (s *service) RemoveInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack) error {
	for ref := range rocketPacks {
		err := s.removeInstallation(ctx, ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) getInstallation(ctx context.Context, reference reference.Reference, rocketPack *rocketpack.RocketPack, readOnly bool) (*Installation, error) {
	s.logger.Debug("Checking installation", map[string]interface{}{"preInstalled": rocketPack.IsPreInstalled(), "readOnly": readOnly, "reference": reference.String()})

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
		return nil
	}

	downloadedFilePath := filepath.Join(installationPath, getFilenameFromURL(downloadUrl))
	err := s.downloader.DownloadWithContext(ctx, downloadedFilePath, downloadUrl)
	if err != nil {
		return err
	}

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

	err := os.RemoveAll(installationPath)
	if err != nil {
		return err
	}

	return nil
}

func getFilenameFromURL(downloadURL string) string {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return ""
	}

	return path.Base(u.Path)
}
