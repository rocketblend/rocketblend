package installation

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/jot/downloader"
	"github.com/rocketblend/rocketblend/pkg/jot/extractor"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketpack"
)

type (
	Service interface {
		GetInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack, readOnly bool) (map[reference.Reference]*Installation, error)
		RemoveInstallations(ctx context.Context, rocketPacks map[reference.Reference]*rocketpack.RocketPack) error
	}

	Options struct {
		Logger      logger.Logger
		PackagePath string
		StoragePath string
		Platform    runtime.Platform
		Downloader  downloader.Downloader
		Extractor   extractor.Extractor
	}

	Option func(*Options)

	service struct {
		logger      logger.Logger
		packagePath string
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

func WithPackagePath(packagePath string) Option {
	return func(o *Options) {
		o.PackagePath = packagePath
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

	if options.Downloader == nil {
		return nil, fmt.Errorf("downloader is required")
	}

	if options.Extractor == nil {
		return nil, fmt.Errorf("extractor is required")
	}

	if options.StoragePath == "" {
		return nil, fmt.Errorf("storage path is required")
	}

	if options.PackagePath == "" {
		return nil, fmt.Errorf("package path is required")
	}

	if options.Platform == runtime.Undefined {
		return nil, fmt.Errorf("platform is required")
	}

	err := os.MkdirAll(options.StoragePath, 0755)
	if err != nil {
		return nil, err
	}

	return &service{
		logger:      options.Logger,
		storagePath: options.StoragePath,
		packagePath: options.PackagePath,
		downloader:  options.Downloader,
		extractor:   options.Extractor,
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
	return nil
}

func (s *service) getInstallation(ctx context.Context, reference reference.Reference, rocketPack *rocketpack.RocketPack, readOnly bool) (*Installation, error) {
	executableName, err := rocketPack.GetExecutableName(s.platform)
	if err != nil {
		return nil, err
	}

	installationPath := filepath.Join(s.storagePath, reference.String())
	executablePath := filepath.Join(installationPath, executableName)

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
		// Local installation
		return nil
	}

	err := s.downloader.DownloadWithContext(ctx, installationPath, downloadUrl)
	if err != nil {
		return err
	}

	downloadedFilePath := filepath.Join(installationPath, getFilenameFromURL(downloadUrl))
	if isArchive(downloadedFilePath) {
		err = s.extractor.ExtractWithContext(ctx, downloadedFilePath, filepath.Dir(downloadedFilePath))
		if err != nil {
			return err
		}

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
